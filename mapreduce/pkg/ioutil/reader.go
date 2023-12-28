package ioutil

import (
	"encoding/binary"
	"io"
	"math/bits"

	"google.golang.org/protobuf/proto"
)

type integer interface {
	~uint32 | ~uint64 | ~int32 | ~int64 | ~int
}

func varIntLen64(v uint64) int {
	if v == 0 {
		return 1
	}
	lz := 64 - bits.LeadingZeros64(v)
	return (lz + 6) / 7
}

func varIntLen[T integer, U integer](v T) U {
	return U(varIntLen64(uint64(v)))
}

type ByteReadReader interface {
	io.ByteReader
	io.Reader
	io.Closer
}

type MessagePointer[T any] interface {
	proto.Message
	*T
}

type MessageReader[M any, T MessagePointer[M]] interface {
	Recv() (T, error)
}

// Reader reads length-prefixed proto.Message objects from an io.Reader.
//
// The anticipated format is a length encoded as a varint [1] followed
// immediately by a number of bytes precisely equaling the encoded length.
// These message blocks MAY be separated by any number of NULL bytes.
//
// If BlockSize is nonzero, it is assumed that NULL bytes outside of
// messages only occur in contiguous chunks at the end of a block. As such,
// a Reader MAY seek ahead to the beginning of the next block the first time
// it encounters a NULL byte outside a message.
//
// As a consequence of allowing this padding, it is not possible to encode
// a zero-valued message, since it would have zero length and its length
// prefix would therefore be a NULL byte outside a message. To compensate
// for this, we have a special case where if the length is 1, we consider
// it to encode a zero-valued message. This is possible because a protobuf
// message is encoded as pairs of <tag, value> and can therefore never be
// exactly one byte in length.
//
// 1: https://protobuf.dev/programming-guides/encoding/#varints
type Reader[M any, T MessagePointer[M]] struct {
	r ByteReadReader

	i         int64
	block     int64
	blockSize int64
}

func (r *Reader[M, T]) BlockSize() int64 {
	return r.blockSize
}

func (r *Reader[M, T]) Close() error {
	return r.r.Close()
}

// Recv gets the next encoded message from the Reader.
//
// First we read the length, which is encoded as an unsigned varint.
// If the length is zero, we skip ahead until we read a non-zero length
// or hit EOF. If the BlockSize is nonzero and the underlying io.Reader
// is an io.Seeker, the skip-ahead will be achieved by seeking to the
// start of the next block; otherwise we read varints repeatedly until
// a non-zero value or EOF is found.
//
// If the length is exactly 1, that is a special case, and we return a
// zero-value message. This is to allow for the use of NULL bytes as padding.
//
// After the length has been determined, we read precisely that many bytes
// off the underlying io.Reader and unmarshal them into the target type.
// If fewer bytes are read, io.ErrUnexpectedEOF is returned.
//
// If a decode error occurs in the unmarshaling, it is propagated.
// When the end of the file is reached, io.EOF is returned as the error.
func (r *Reader[M, T]) Recv() (T, error) {
	var size uint64
	var err error
	for size == 0 {
		size, err = binary.ReadUvarint(r.r)
		if err != nil {
			return nil, err
		}
		{
			tmp := r.i + varIntLen[uint64, int64](size)
			r.i = tmp % r.blockSize
			r.block += tmp / r.blockSize
		}

		if s, ok := r.r.(io.Seeker); r.blockSize != 0 && ok {
			n, err := s.Seek(r.blockSize-r.i, io.SeekCurrent)
			if err != nil {
				return nil, err
			}
			if n != (r.block+1)*r.blockSize {
				return nil, io.EOF
			}
			r.block += 1
			r.i = 0
		}
	}

	if size == 1 {
		return new(M), nil
	}

	b := make([]byte, size)
	n, err := r.r.Read(b)
	if err != nil {
		return nil, err
	}

	if n < int(size) {
		return nil, io.ErrUnexpectedEOF
	}
	{
		tmp := r.i + int64(size)
		r.i = tmp % r.blockSize
		r.block += tmp / r.blockSize
	}

	var out T = new(M)
	err = proto.Unmarshal(b, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// CombinedReader combines the output of inputs implementing MessageReader.
type CombinedReader[M any, T MessagePointer[M]] struct {
	rs   []MessageReader[M, T]
	less func(a T, b T) bool

	buf  []T
	dead []MessageReader[M, T]
}

// NewConcatReader returns a new MergedReader that concatenates inputs.
//
// Every MessageReader passed as input is drained to first error, which
// is typically io.EOF, sequentially.
func NewConcatReader[M any, T MessagePointer[M]](
	in ...MessageReader[M, T],
) *CombinedReader[M, T] {
	return &CombinedReader[M, T]{
		rs: in,
	}
}

// NewMergeSortingReader returns a new MergedReader that sorts output.
//
// The sorting requirement is recursive: any MessageReader passed MUST
// also return sorted output. The MergedReader will buffer one response
// from each MessageReader and return the smallest one, as evaluated by
// the comparator passed in. This process is repeated until every
// MessageReader has responded with a non-nil error, typically io.EOF.
func NewMergeSortingReader[M any, T MessagePointer[M]](
	less func(a T, b T) bool,
	in ...MessageReader[M, T],
) *CombinedReader[M, T] {
	return &CombinedReader[M, T]{
		rs:   in,
		less: less,
	}
}

func defaultLess[T any](_, _ T) bool {
	return false
}

// Recv returns the next response from an input MessageReader.
//
// One response is buffered from each MessageReader. If a non-nil comparator
// has been defined, the least element is selected using that comparator.
// otherwise, the first non-nil element returned from any MessageReader
// is used.
//
// If every MessageReader has been exhausted (read until it returns an error),
// returns nil, io.EOF.
func (r *CombinedReader[M, T]) Recv() (T, error) {
	if r.buf == nil {
		r.buf = make([]T, len(r.rs))
	}

	for i, v := range r.rs {
		if r.buf[i] != nil || v == nil {
			continue
		}

		m, err := v.Recv()
		if err != nil {
			r.dead = append(r.dead, r.rs[i])
			r.rs[i] = nil
			continue
		}

		r.buf[i] = m
	}

	less := r.less
	if less == nil {
		less = defaultLess[T]
	}

	minIdx := 0
	for i := 1; i < len(r.buf); i++ {
		if r.buf[i] == nil {
			continue
		}
		if r.buf[minIdx] == nil {
			minIdx = i
			continue
		}

		if less(r.buf[i], r.buf[minIdx]) {
			minIdx = i
		}
	}

	if r.buf[minIdx] == nil {
		return nil, io.EOF
	}

	out := r.buf[minIdx]
	r.buf[minIdx] = nil

	return out, nil
}

// Close will close any MessageReader that has been exhausted, if possible.
//
// MessageReader does not imply io.Closer, but many file-backed MessageReader
// implementations MAY implement io.Closer. This will close any exhausted
// MessageReader that implements io.Closer, silently ignoring any that do not.
//
// Only the LAST error returned by any io.Closer is returned; any other errors
// are suppressed, and nil indicates no errors occurred.
func (r *CombinedReader[M, T]) Close() error {
	var err error
	for _, v := range r.dead {
		if c, ok := v.(io.Closer); ok {
			err = c.Close()
		}
	}
	r.dead = nil

	return err
}

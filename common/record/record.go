package record

import (
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
	"math/bits"
	"sort"

	"github.com/google/uuid"
	"github.com/hxtk/yggdrasil/third_party/crc"
)

// Chunk format:
// The record ID and chunk ID are the first two values, and the chunk ID uses
// big-endian, fixed-length encoding so that sorting the chunks will group record
// members and put them in order.
// +------------------------------------------------------------------------------+
// | recordID (16B) | chunkID (4B) | kind (1B) | checksum (4B) | size (uvarint64) |
// +------------------------------------------------------------------------------+
const fixedHeaderSize = 16 + 4 + 1 + 4

var (
	ErrBadChecksum     = errors.New("record: did not pass checksum")
	ErrBadSize         = errors.New("record: could not parse payload size")
	ErrChunkTooSmall   = errors.New("chunk: requested chunk size too small")
	ErrNoFirstChunk    = errors.New("record: need a first chunk combine")
	ErrNotEnoughChunks = errors.New("record: need at least two chunks to combine")
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

func varIntLen[T integer](v T) int {
	return varIntLen64(uint64(v))
}

func headerSize[T integer](payloadSize T) int {
	return varIntLen(payloadSize) + fixedHeaderSize
}

func chunkSize[T integer](payloadSize T) int {
	return headerSize(payloadSize) + int(payloadSize)
}

type ChunkKind byte

const (
	FullRecord ChunkKind = iota
	FirstChunk
	MiddleChunk
	LastChunk
)

type Record struct {
	recordID uuid.UUID
	chunkID  uint32
	kind     ChunkKind
	payload  []byte
}

var _ encoding.BinaryMarshaler = new(Record)
var _ encoding.BinaryUnmarshaler = new(Record)

func NewRecord(p []byte) *Record {
	id, uuidErr := uuid.NewV7()
	return &Record{
		recordID: uuid.Must(id, uuidErr),
		chunkID:  0,
		kind:     FullRecord,
		payload:  p,
	}
}

func MergeRecords(rs []*Record) (*Record, error) {
	if len(rs) < 2 {
		return nil, ErrNotEnoughChunks
	}
	var id *uuid.UUID
	var payload []byte
	chunkID := uint32(0)
	for _, v := range rs {
		if chunkID == 0 && v.kind != FirstChunk {
			return nil, ErrNoFirstChunk
		}
		if chunkID != v.chunkID {
			//TODO: Handle Error
			//return nil,
		}
		if id == nil {
			*id = v.recordID
		} else if v.recordID != *id {
			continue // TODO: skip or error? Need use cases.
		}
		payload = append(payload, v.payload...)
		if v.kind == LastChunk {
			break
		}
		chunkID++
	}

	return &Record{
		recordID: *id,
		chunkID:  0,
		kind:     FullRecord,
		payload:  payload,
	}, nil
}

func newRecordWithID(id uuid.UUID, p []byte) *Record {
	return &Record{
		recordID: id,
		chunkID:  0,
		kind:     FullRecord,
		payload:  p,
	}
}

func (r *Record) MarshalBinary() ([]byte, error) {
	id, err := r.recordID.MarshalBinary()
	if err != nil {
		return nil, err
	}

	b := make([]byte, 0, chunkSize(len(r.payload)))
	b = append(b, id...)
	b = binary.BigEndian.AppendUint32(b, r.chunkID)
	b = append(b, byte(r.kind))
	b = binary.BigEndian.AppendUint32(b, r.Checksum())
	b = binary.AppendUvarint(b, uint64(len(r.payload)))
	//fmt.Println(r.payload)
	b = append(b, r.payload...)

	//fmt.Println(b[16:])

	return b, nil
}

func (r *Record) UnmarshalBinary(p []byte) error {
	err := r.recordID.UnmarshalBinary(p[:16])
	if err != nil {
		return err
	}

	r.chunkID = binary.BigEndian.Uint32(p[16:20])
	r.kind = ChunkKind(p[20])
	checksum := binary.BigEndian.Uint32(p[21:25])

	size, n := binary.Uvarint(p[25:])
	if n <= 0 {
		return ErrBadSize
	}
	fmt.Println(n, size, p)
	r.payload = p[25+n : 25+n+int(size)]

	if checksum != r.Checksum() {
		return ErrBadChecksum
	}
	return nil
}

// Checksum calculates a CRC checksum over the RecordID, ChunkID, Kind, and Payload.
func (r *Record) Checksum() uint32 {
	idBytes, _ := r.recordID.MarshalBinary()
	idBytes = binary.BigEndian.AppendUint32(idBytes, r.chunkID)

	checksum := crc.New(idBytes)
	checksum.Update([]byte{byte(r.kind)})
	checksum.Update(r.payload)

	return checksum.Value()
}

// Size calculates the encoded size of the chunk.
func (r *Record) Size() int {
	return chunkSize(len(r.payload))
}

func (r *Record) Chunk(maxSize int) ([]*Record, error) {
	// In general, this would be handled below, but testing for it specifically
	// is not just a micro-optimization. It is essential for correctness in the
	// case of an empty payload. If we find that the largest payload size we can
	// put in a chunk is 0, that is interpreted below to mean that the requested
	// chunk size is too small for us to make progress; however, it is also
	// possible that we cannot produce a chunk with a payload size greater than
	// zero because the whole record's payload is empty. By checking here to see
	// if the whole record fits in one chunk, we simplify the check that would
	// otherwise require a special case for len(r.payload)==0 && maxSize >= 26.
	if r.Size() <= maxSize {
		return []*Record{r}, nil
	}

	// sort.Search returns the first index at which the search function
	// evaluates as true. In this case, the partition we want is the point
	// where the encoded size of a record exceeds the maximum allowed size.
	// However, we want to be under that limit, so we evaluate at x+1 to find
	// the largest payload whose encoded size does not exceed the maxSize.
	size := sort.Search(len(r.payload)-1, func(x int) bool {
		size := chunkSize(x + 1)
		return size > maxSize
	})

	if size == 0 {
		return nil, ErrChunkTooSmall
	}

	payload := r.payload
	var out []*Record
	for len(payload) > 0 {
		chunkID := uint32(len(out))

		chunkKind := MiddleChunk
		if chunkID == 0 {
			chunkKind = FirstChunk
		} else if len(payload) <= size {
			chunkKind = LastChunk
		}

		cut := min(size, len(payload))
		var p []byte
		p, payload = payload[:cut], payload[cut:]

		out = append(out, &Record{
			recordID: r.recordID,
			chunkID:  chunkID,
			kind:     chunkKind,
			payload:  p,
		})
	}

	return out, nil
}

package record

import (
	"bytes"
	"io"
	"os"
	"sync"
	"sync/atomic"
)

// Log represents a file on a device supporting atomic, block-aligned appends.
type Log interface {
	io.Closer

	// RecordAppend writes a buffer to a file with alignment.
	//
	// If writing succeeds without error, it is guaranteed that the contents
	// of `p` shall appear in the underlying File aligned to a block boundary.
	// Critically, this interface does not have any functional concepts related
	// to blocks: the alignment occurs transparently.
	//
	// If `p` will fit in the current block then it is appended unaltered.
	// Otherwise, the current block is atomically filled with nul bytes and the
	// record is atomically written to the next block. These MAY be two
	// separate atomic writes or one.
	//
	// If a buffer is bigger than RecordLimit, it is rejected with
	// ErrRecordTooBig.
	//
	// On success, returns the index just past the end of the written record.
	RecordAppend(p []byte) (int64, error)

	// RecordLimit returns the maximum size of a single record. Records larger
	// than the block size are not supported, as it must be possible to do
	// operations on the blocks of a file independently, which is not possible
	// if a record spans multiple blocks.
	//
	// Implementors MUST return a RecordLimit no larger than the file's
	// block size. Smaller record limits are less versatile because they limit
	// the total amount of data that can be written in a single record, but
	// a high ratio of RecordLimit to block size reduces overall utilization.
	// In the worst case, if RecordLimit = BlockSize/2 + 1, it is possible that
	// nearly half of each block would be wasted by padding.
	RecordLimit() int
}

type LogWriter struct {
	f  Log
	mu sync.Mutex
	ws []*singleWriter
}

func (l *LogWriter) Flush() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	var ws []*singleWriter
	for _, v := range l.ws {
		if !v.done.Load() {
			ws = append(ws, v)
			continue
		}

		_, err := l.WriteRecord(v.buf.Bytes())
		if err != nil {
			return err
		}
	}

	l.ws = ws
	return nil
}

func (l *LogWriter) Next() io.WriteCloser {
	l.mu.Lock()
	defer l.mu.Unlock()

	out := &singleWriter{
		p: l,
	}

	l.ws = append(l.ws, out)
	return out
}

func (l *LogWriter) WriteRecord(p []byte) (int64, error) {
	chunks, err := NewRecord(p).Chunk(l.f.RecordLimit())
	if err != nil {
		return 0, err
	}

	var end int64
	for _, v := range chunks {
		recordBytes, err := v.MarshalBinary()
		if err != nil {
			return 0, err
		}
		end, err = l.f.RecordAppend(recordBytes)
		if err != nil {
			return 0, err
		}
	}

	return end, nil
}

type singleWriter struct {
	p   *LogWriter
	buf bytes.Buffer

	done atomic.Bool
}

var _ io.WriteCloser = new(singleWriter)

func (s *singleWriter) Write(p []byte) (n int, err error) {
	if s.done.Load() {
		return 0, os.ErrClosed
	}
	return s.buf.Write(p)
}

func (s *singleWriter) Close() error {
	s.done.Store(true)
	return nil
}

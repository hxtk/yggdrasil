package record

import (
	"errors"
	"io"
	"io/fs"

	"github.com/hxtk/yggdrasil/third_party/filelock"
)

var (
	ErrRecordTooBig = errors.New("cannot atomically append a record that large")
)

type StatLockWriter interface {
	filelock.File
	io.WriterAt
	io.Writer

	// Stat returns the FileInfo structure describing file.
	Stat() (fs.FileInfo, error)

	// Sync flushes the file contents to disk.
	Sync() error
}

// AtomicFile is a file that supports atomic writes.
//
// In particular, there is a primitive for an atomic RecordAppend,
// which guarantees that the bytes written shall occur entirely within a
// single block of the file, using a user-defined block size.
//
// Atomicity is achieved with OS-level advisory locks, which support
// cooperative mutual exclusion. Care must be taken to ensure that all
// processes contending for write access are honoring OS-specific advisory
// locks.
//
// Every RecordAppend is synchronized to the disk before the lock is released.
// This results in relatively low throughput in exchange for safety guarantees.
type AtomicFile struct {
	StatLockWriter

	blockSize int
	maxRecord int
}

// NewAtomicFile creates a File with the specified blockSize.
//
// The maximum record size defaults to 1/4 of the blockSize to limit
// wasted space at the end of incomplete blocks. This value, along with the
// Record Append concept itself is borrowed from GFS [1].
//
// 1: https://research.google/pubs/the-google-file-system/
func NewAtomicFile(w StatLockWriter, blockSize int) *AtomicFile {
	return &AtomicFile{
		StatLockWriter: w,
		blockSize:      blockSize,
		maxRecord:      blockSize >> 2,
	}
}

func NewAtomicFileWithRecordSize(w StatLockWriter, blockSize, recordSize int) *AtomicFile {
	return &AtomicFile{
		StatLockWriter: w,
		blockSize:      blockSize,
		maxRecord:      recordSize,
	}
}

func (f *AtomicFile) Write(p []byte) (n int, err error) {
	err = filelock.Lock(f.StatLockWriter)
	if err != nil {
		return
	}

	defer func() {
		unlockErr := filelock.Unlock(f.StatLockWriter)
		if err == nil {
			err = unlockErr
		}
	}()

	return f.StatLockWriter.Write(p)
}

func (f *AtomicFile) RecordAppend(p []byte) (n int64, err error) {
	if f.maxRecord == 0 {
		f.maxRecord = f.blockSize >> 2
	}
	if len(p) > f.maxRecord {
		return 0, ErrRecordTooBig
	}

	err = filelock.Lock(f.StatLockWriter)
	if err != nil {
		return
	}

	defer func() {
		unlockErr := filelock.Unlock(f.StatLockWriter)
		if err == nil {
			err = unlockErr
		}
	}()

	info, err := f.StatLockWriter.Stat()
	if err != nil {
		return 0, err
	}

	size := info.Size()
	remaining := f.blockSize - (int(size) % f.blockSize)
	payload := p
	if remaining < len(p) {
		payload = make([]byte, remaining+len(p))
		copy(payload[remaining:], p)
	}

	written, err := f.StatLockWriter.WriteAt(payload, size)
	if err != nil {
		return 0, err
	}

	err = f.StatLockWriter.Sync()
	if err != nil {
		return 0, err
	}
	return int64(written) + size, nil
}

func (f *AtomicFile) Size() int64 {
	stat, _ := f.StatLockWriter.Stat()
	return stat.Size()
}

func (f *AtomicFile) RecordLimit() int {
	return f.maxRecord
}

func (f *AtomicFile) Close() error {
	if closer, ok := f.StatLockWriter.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

type StatWriter interface {
	io.WriterAt
	io.Writer

	// Stat returns the FileInfo structure describing file.
	Stat() (fs.FileInfo, error)

	// Sync flushes the file contents to disk.
	Sync() error
}

// ExclusiveFile is a file that supports RecordAppend for a single writer.
//
// Under contention, writes may be interleaved unless maxRecord is less than
// PIPE_BUF (at least 512B per POSIX standard; 4096B on Linux and macOS).
// However, in cases where there is only a single writer process, it is much
// faster because it does not synchronize to the disk for every write.
type ExclusiveFile struct {
	StatWriter

	blockSize int
	maxRecord int
}

// NewFile creates a File with the specified blockSize.
//
// The maximum record size defaults to 1/4 of the blockSize to limit
// wasted space at the end of incomplete blocks. This value, along with the
// Record Append concept itself is borrowed from GFS [1].
//
// 1: https://research.google/pubs/the-google-file-system/
func NewFile(w StatWriter, blockSize int) *ExclusiveFile {
	return &ExclusiveFile{
		StatWriter: w,
		blockSize:  blockSize,
		maxRecord:  blockSize >> 2,
	}
}

func NewFileWithRecordSize(w StatWriter, blockSize, recordSize int) *ExclusiveFile {
	return &ExclusiveFile{
		StatWriter: w,
		blockSize:  blockSize,
		maxRecord:  recordSize,
	}
}

func (f *ExclusiveFile) innerWrite(b []byte) (int64, int, error) {
	if f.maxRecord == 0 {
		f.maxRecord = f.blockSize >> 2
	}
	if len(b) > f.maxRecord {
		return 0, 0, ErrRecordTooBig
	}

	var totalWritten int
	var size int64
	for {
		info, err := f.StatWriter.Stat()
		if err != nil {
			return 0, 0, err
		}
		size = info.Size()
		remaining := f.blockSize - (int(size) % f.blockSize)
		if remaining < len(b) {
			payload := make([]byte, remaining)
			written, err := f.StatWriter.WriteAt(payload, size)
			totalWritten += written
			if err != nil {
				return size, totalWritten, err
			}
		} else {
			break
		}
	}

	written, err := f.StatWriter.WriteAt(b, size)
	totalWritten += written

	return size + int64(totalWritten), totalWritten, err
}

// RecordAppend writes a buffer contiguously into a block of the file.
//
// will fit within the block size of the file.
func (f *ExclusiveFile) Write(p []byte) (n int, err error) {
	_, n, err = f.innerWrite(p)
	return
}

// RecordAppend writes a buffer contiguously into a block of the file.
//
// will fit within the block size of the file.
func (f *ExclusiveFile) RecordAppend(p []byte) (n int64, err error) {
	size, _, err := f.innerWrite(p)
	return size, err
}

func (f *ExclusiveFile) Size() int64 {
	stat, _ := f.StatWriter.Stat()
	return stat.Size()
}

func (f *ExclusiveFile) RecordLimit() int {
	return f.maxRecord
}

func (f *ExclusiveFile) Close() error {
	if closer, ok := f.StatWriter.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

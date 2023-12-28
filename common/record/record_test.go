package record

import (
	"errors"
	"testing"
	"time"
	_ "unsafe"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

//go:linkname uuidTimeNow github.com/google/uuid.timeNow
var uuidTimeNow func() time.Time

type fuzzReader struct {
	b []byte
	i int
}

func (f *fuzzReader) Read(b []byte) (int, error) {
	x := copy(b, f.b[f.i:])
	f.i += x
	for i := x; i < len(b); i++ {
		b[i] = 0
	}
	return len(b), nil
}

func FuzzRecordCodec(f *testing.F) {
	f.Fuzz(func(t *testing.T, payload, rand []byte, epoch int64) {
		oldTimeNow := uuidTimeNow
		defer func() {
			uuidTimeNow = oldTimeNow
		}()

		uuidTimeNow = func() time.Time {
			return time.UnixMilli(epoch)
		}
		id, _ := uuid.NewV7FromReader(&fuzzReader{b: rand})
		want := newRecordWithID(id, payload)
		encoded, err := want.MarshalBinary()
		if err != nil {
			t.Fatalf("Error marshaling record: %v", err)
		}

		got := new(Record)
		err = got.UnmarshalBinary(encoded)
		if err != nil {
			t.Fatalf("Error unmarshaling record: %v", err)
		}

		if !cmp.Equal(want, got, cmp.AllowUnexported(*want)) {
			t.Fatalf(
				"Decoded record did not match: %v",
				cmp.Diff(want, got, cmp.AllowUnexported(*want)),
			)
		}
	})
}

func FuzzRecordSize(f *testing.F) {
	f.Fuzz(func(t *testing.T, payload, rand []byte, epoch int64) {
		oldTimeNow := uuidTimeNow
		defer func() {
			uuidTimeNow = oldTimeNow
		}()

		uuidTimeNow = func() time.Time {
			return time.UnixMilli(epoch)
		}
		id, _ := uuid.NewV7FromReader(&fuzzReader{b: rand})
		record := newRecordWithID(id, payload)
		encoded, err := record.MarshalBinary()
		if err != nil {
			t.Fatalf("Error marshaling record: %v", err)
		}

		got := record.Size()
		want := len(encoded)
		if !cmp.Equal(want, got) {
			t.Fatalf(
				"record.Size() = %v; want %v",
				got,
				want,
			)
		}
	})
}

func FuzzRecordChunk(f *testing.F) {
	f.Fuzz(func(t *testing.T, payload, rand []byte, epoch int64, size uint32) {
		if int(size) < fixedHeaderSize+2 {
			t.Skip("Chunks too small.")
		}
		oldTimeNow := uuidTimeNow
		defer func() {
			uuidTimeNow = oldTimeNow
		}()

		uuidTimeNow = func() time.Time {
			return time.UnixMilli(epoch)
		}
		id, _ := uuid.NewV7FromReader(&fuzzReader{b: rand})

		want := newRecordWithID(id, payload)
		chunks, err := want.Chunk(int(size))
		if errors.Is(err, ErrChunkTooSmall) {
			t.Fatalf("Got chunks too small, but they aren't: %v", err)
		} else if err != nil {
			t.Fatalf("Error chunking: %v.", err)
		}
		for i, v := range chunks {
			d, err := v.MarshalBinary()
			if err != nil {
				t.Fatalf("chunks[%d].MarshalBinary() = _, %v", i, err)
			}
			if len(d) > int(size) {
				t.Fatalf("d, _ := chunks[%d]; len(d) = %d; want <= %d", i, len(d), size)
			}
		}
	})
}

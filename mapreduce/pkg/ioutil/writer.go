package ioutil

import (
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"io"
)

type Writer[M any, T MessagePointer[M]] struct {
	w io.Writer
}

func NewWriter[M any, T MessagePointer[M]](w io.Writer) *Writer[M, T] {
	return &Writer[M, T]{
		w: w,
	}
}

func (w *Writer[M, T]) Close() error {
	if c, ok := w.w.(io.Closer); ok {
		return c.Close()
	}

	return nil
}

func (w *Writer[M, T]) Send(p T) error {
	data, err := proto.Marshal(p)
	if err != nil {
		return err
	}

	// Special case for empty-value.
	if len(data) == 0 {
		_, err := w.w.Write([]byte{1})
		return err
	}

	b := make([]byte, 0, len(data)+varIntLen[int, int](len(data)))
	b = binary.AppendUvarint(b, uint64(len(data)))
	b = append(b, data...)

	_, err = w.w.Write(b)
	return err
}

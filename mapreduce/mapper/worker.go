package mapper

import (
	"context"
	"io"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/hxtk/yggdrasil/mapreduce/v1alpha1"
)

type RecordReader interface {
	Recv() (*pb.Record, error)
}

type RecordWriter interface {
	Send(*pb.Record) error
	CloseAndRecv() (*emptypb.Empty, error)
}

type MapWorker struct {
	m pb.MapperClient
	c pb.ReducerClient
	r RecordReader
	w RecordWriter
}

func NewWorker(mapper pb.MapperClient, reader RecordReader) *MapWorker {
	return &MapWorker{
		m: mapper,
		r: reader,
	}
}

func (m *MapWorker) Run(ctx context.Context) error {
	for {
		rec, err := m.r.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		recv, err := m.m.MapRecord(ctx, rec)
		if err != nil {
			return err
		}

	L:
		for {
			select {
			case <-ctx.Done():
				break L
			default:
				//noop
			}

			r, err := recv.Recv()
			if err != nil {
				break
			}

			err = m.w.Send(r)
			if err != nil {
				break
			}
		}
	}

	_, err := m.w.CloseAndRecv()
	return err
}

package reducer

import (
	"bytes"
	"context"
	"github.com/hxtk/yggdrasil/mapreduce/pkg/ioutil"
	pb "github.com/hxtk/yggdrasil/mapreduce/v1alpha1"
	"io"
)

type RecordFile interface {
	io.Reader
	io.ByteReader
}

func recordLess(a *pb.Record, b *pb.Record) bool {
	if a == nil {
		return false
	}
	if b == nil {
		return true
	}
	return bytes.Compare(a.GetKey(), b.GetKey()) < 0
}

type ReduceWorker struct {
	in ioutil.MessageReader[pb.Record, *pb.Record]
	c  pb.ReducerClient
}

func NewReduceWorker(c pb.ReducerClient, in ...ioutil.MessageReader[pb.Record, *pb.Record]) *ReduceWorker {
	return &ReduceWorker{
		in: ioutil.NewMergeSortingReader[pb.Record, *pb.Record](recordLess, in...),
		c:  c,
	}
}

func (r *ReduceWorker) Run(ctx context.Context) (err error) {
	var key []byte
	var stream pb.Reducer_ReduceRecordsClient
	for {
		rec, err := r.in.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if !bytes.Equal(rec.GetKey(), key) {
			stream, err = r.c.ReduceRecords(ctx)
			if err != nil {
				return err
			}
		}
		err = stream.Send(rec)
		if err != nil {
			return err
		}
	}

	return stream.CloseSend()
}

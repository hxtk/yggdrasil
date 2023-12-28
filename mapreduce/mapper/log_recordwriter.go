package mapper

import (
	"bytes"
	"io"

	"github.com/google/btree"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hxtk/yggdrasil/mapreduce/pkg/ioutil"
	pb "github.com/hxtk/yggdrasil/mapreduce/v1alpha1"
	"github.com/hxtk/yggdrasil/third_party/crc"
)

// PartitionFunc accepts a key and returns a non-negative int.
//
// The value may be arbitrary. When a value is written to a partition,
// the partition is chosen by taking this int modulo the number of partitions.
// As a result, any regularities in this return value will result in uneven
// distribution of work. For example, a constant PartitionFunc would result in
// one reducer getting all the work and the rest getting none.
//
// The easiest way to ensure equal distribution is to have the PartitionFunc
// extract a key from b and then hash that key. For example, if b is a URL then
// the following PartitionFunc would ensure all URLs with the same host went to
// the same shard, without otherwise impacting the distribution:
//
//	func URLPartitionFunc(b []byte) int {
//	    uri, err := url.Parse(string(b))
//	    if err != nil {
//	        return 0  // Not expected to ever occur.
//	    }
//	    u32 := crc.New([]byte(uri.Hostname())
//	    return int(u32 >> 1)  // Trim highest bit to ensure it fits
//	}
type PartitionFunc func(b []byte) int

func defaultPartitionFunc(b []byte) int {
	hash := crc.New(b)
	positive := hash.Value() >> 1
	return int(positive)
}

func itemLess(a, b item) bool {
	return bytes.Compare(a.key, b.key) < 0
}

type item struct {
	key    []byte
	values [][]byte
}

func (it item) merge(r *pb.Record) item {
	out := item{
		key:    it.key,
		values: append(it.values, r.GetValue()),
	}

	if out.key == nil {
		out.key = r.Key
	}

	return out
}

type LogGroup interface {
	Next() (io.WriteCloser, error)
}

// LogRecordWriter receives records sent by the mapper and writes them.
//
// Writes are partitioned into the provided list of LogPrefix paths,
// which SHOULD be non-overlapping prefixes corresponding to the number
// of Reduce tasks in the MapReduce job.
//
// There may be arbitrarily many log files in a prefix directory.
// Log files are limited in size: a log file is written one time. Precisely,
// one batch of log files (one per partition) is written whenever the size of
// the buffered records is about to exceed the size limit.
//
// By writing to log files precisely once, we are able to sort the records as
// we receive them, buffering them in a btree, before writing them and ensure
// that all files are internally sorted, permitting efficient merge of the
// output files into a single sorted file.
type LogRecordWriter struct {
	ps []LogGroup
	p  PartitionFunc
	b  []byte
	t  *btree.BTreeG[item]

	size  int
	limit int
}

var _ RecordWriter = new(LogRecordWriter)

// NewRecordWriter constructs a RecordWriter that writes records to files.
//
// The list of LogGroups passed in represents the collection of output
// partitions, corresponding to Reduce tasks.
//
// The PartitionFunc maps keys to non-negative integers, and is used to
// determine to which LogGroup a record will be written. Though this facility
// MAY be used to map specific records to specific log groups, implementors
// SHOULD prefer only to provide consistent results for items that are to be
// grouped together. This allows for more horizontal scalability of the
// Reducer tasks without rewriting the PartitionFunc.
//
// The bufSize is the maximum size of the buffer contents. When the buffer has
// grown to this size, a new file will be created for each LogGroup, containing
// all data for that LogGroup currently in the buffer. The calculated size is
// only an approximation, accurately reflecting neither the in-memory nor
// on-disk size. The size is approximated as the sum of the byte lengths of all
// distinct keys and all values.
func NewRecordWriter(prefixes []LogGroup, partition PartitionFunc, bufSize int) *LogRecordWriter {
	return &LogRecordWriter{
		ps: prefixes,
		p:  partition,
		t:  btree.NewG[item](2, itemLess),

		limit: bufSize,
	}
}

// CloseAndRecv flushes the current contents of the buffer to disk.
//
// This is purely a wrapper for Flush to implement RecordWriter.
func (l *LogRecordWriter) CloseAndRecv() (*emptypb.Empty, error) {
	return &emptypb.Empty{}, l.Flush()
}

func (l *LogRecordWriter) Send(r *pb.Record) error {
	it, ok := l.t.Get(item{key: r.GetKey()})
	if ok {
		l.size += len(r.GetKey()) + len(r.GetValue())
	} else {
		l.size += len(r.GetValue())
	}

	it.merge(r)
	l.t.ReplaceOrInsert(it)

	if l.size > l.limit {
		return l.Flush()
	}
	return nil

}

// Flush writes the current record buffer to disk.
//
// Iterate over the buffer tree in-order and calculate to which partition
// each buffered record belongs. Open a file for each partition to which any
// record belongs and write those matching records to the file in iteration
// order.
//
// Each call to Flush will create zero or one record files in the
func (l *LogRecordWriter) Flush() (ret error) {
	var outerErr error
	ls := make([]*ioutil.Writer[pb.Record, *pb.Record], len(l.ps))
	defer func() {
		for _, v := range ls {
			err := v.Close()
			if ret != nil {
				ret = err
			}
		}
	}()

	pf := l.p
	if pf == nil {
		pf = defaultPartitionFunc
	}

	l.t.Ascend(func(i item) bool {
		for _, v := range i.values {
			r := &pb.Record{
				Key:   i.key,
				Value: v,
			}

			partition := pf(r.GetKey()) % len(ls)
			if ls[partition] == nil {
				w, err := l.ps[partition].Next()
				if err != nil {
					outerErr = err

					// We cannot continue because we failed to open the
					// output file for this partition.
					return false
				}
				ls[partition] = ioutil.NewWriter[pb.Record, *pb.Record](w)
			}

			err := ls[partition].Send(r)
			if err != nil {
				outerErr = err
				return true
			}
		}

		return true
	})

	l.size = 0
	l.t.Clear(false)

	if outerErr != nil {
		return outerErr
	}

	return nil
}

var _ RecordWriter = new(LogRecordWriter)

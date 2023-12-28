package mapper

import (
	"crypto/rand"
	"io"
	"os"
	"path"

	"github.com/google/uuid"

	"github.com/hxtk/yggdrasil/common/record"
)

var randReader io.Reader = rand.Reader

type Opener interface {
	Open(path string) (record.StatWriter, error)
}

type osOpener struct{}

func (o *osOpener) Open(path string) (record.StatWriter, error) {
	return os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

type LogPrefix struct {
	root      string
	blockSize int
	fs        Opener
}

func (l *LogPrefix) Next() (io.WriteCloser, error) {
	id, err := uuid.NewV7FromReader(randReader)
	if err != nil {
		return nil, err
	}
	name := path.Join(l.root, id.String())
	f, err := l.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return record.NewFile(f, l.blockSize), nil
}

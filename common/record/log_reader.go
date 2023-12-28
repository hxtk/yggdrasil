package record

import "io"

type LogReader struct {
	f io.ReadSeeker
}

package nextpass

import "io"

type readCounter struct {
	count  int
	reader io.Reader
}

func newReadCounter(reader io.Reader) *readCounter {
	return &readCounter{0, reader}
}

func (ctr *readCounter) Read(b []byte) (n int, err error) {
	n, err = ctr.reader.Read(b)
	ctr.count += n
	return n, err
}

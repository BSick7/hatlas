package terraform

import (
	"bytes"
	"io"
)

type Payload struct {
	MD5  []byte
	Data []byte
}

func (p *Payload) GetReader() io.ReadSeeker {
	return bytes.NewReader(p.Data)
}

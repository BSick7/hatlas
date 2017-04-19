package terraform

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

type Payload struct {
	MD5  []byte
	Data []byte
}

func NewPayloadFromString(str string) *Payload {
	data := bytes.NewBufferString(str).Bytes()
	hash := md5.Sum(data)
	return &Payload{
		Data: data,
		MD5:  hash[:],
	}
}

func NewPayloadFromResponse(res *http.Response) (*Payload, error) {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, res.Body); err != nil {
		return nil, err
	}

	data := buf.Bytes()
	if len(data) == 0 {
		return nil, nil
	}

	md5, err := decodeMD5(res, data)
	if err != nil {
		return nil, err
	}

	return &Payload{
		MD5:  md5,
		Data: data,
	}, nil
}

func decodeMD5(res *http.Response, data []byte) ([]byte, error) {
	// Check for the MD5
	if raw := res.Header.Get("Content-MD5"); raw != "" {
		md5, err := base64.StdEncoding.DecodeString(raw)
		if err != nil {
			return nil, fmt.Errorf("failed to decode Content-MD5 '%s': %v", raw, err)
		}
		return md5, nil
	}

	// Generate the MD5
	hash := md5.Sum(data)
	return hash[:], nil
}

func (p *Payload) GetReader() io.ReadSeeker {
	return bytes.NewReader(p.Data)
}

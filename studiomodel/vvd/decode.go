package vvd

import (
	"bytes"
	"io"
)

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{reader: r}
}

type Decoder struct {
	reader io.Reader
	r      *bytes.Reader
}

func (d *Decoder) Decode() (*VVD, error) {
	data, err := io.ReadAll(d.reader)
	if err != nil {
		return nil, err
	}
	d.r = bytes.NewReader(data)
	return d.decodeVVD()
}

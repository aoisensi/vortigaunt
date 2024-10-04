package vtx

import (
	"bytes"
	"encoding/binary"
	"io"
)

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{reader: r}
}

var le = binary.LittleEndian

type Decoder struct {
	reader io.Reader
	r      *bytes.Reader
}

func (d *Decoder) Decode() (*VTX, error) {
	data, err := io.ReadAll(d.reader)
	if err != nil {
		return nil, err
	}
	d.r = bytes.NewReader(data)
	vtx := new(VTX)
	vtx.Header = new(Header)
	if err := binary.Read(d.r, le, vtx.Header); err != nil {
		return nil, err
	}
	if err := d.decodeBodyParts(vtx); err != nil {
		return nil, err
	}
	return vtx, nil
}

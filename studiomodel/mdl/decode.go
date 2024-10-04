package mdl

import (
	"bytes"
	"encoding/binary"
	"io"
	"strings"
)

type Decoder struct {
	reader io.Reader
	data   []byte
	r      *bytes.Reader
	mdl    *MDL
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{reader: r}
}

var le = binary.LittleEndian

func (d *Decoder) Decode() (*MDL, error) {
	var err error
	d.data, err = io.ReadAll(d.reader)
	if err != nil {
		return nil, err
	}
	d.r = bytes.NewReader(d.data)
	return d.decodeMDL()
}

func (d *Decoder) ppush(p int32, f func() error) error {
	o, _ := d.r.Seek(0, io.SeekCurrent)
	d.r.Seek(int64(p), io.SeekCurrent)
	if err := f(); err != nil {
		return err
	}
	d.r.Seek(o, io.SeekStart)
	return nil
}

func (d *Decoder) read(data any) error {
	return binary.Read(d.r, le, data)
}

func (d *Decoder) readName(p int32) string {
	o, _ := d.r.Seek(0, io.SeekCurrent)
	d.r.Seek(int64(p), io.SeekCurrent)
	str := make([]rune, 0, 256)
	for {
		c, _, err := d.r.ReadRune()
		if err != nil {
			panic(err)
		}
		if c == 0 {
			d.r.Seek(o, io.SeekStart)
			return string(str)
		}
		str = append(str, c)
	}
}

func (d *Decoder) readNames(p int32, count int32) []string {
	o, _ := d.r.Seek(0, io.SeekCurrent)
	d.r.Seek(int64(p), io.SeekCurrent)
	strs := make([]string, count)
	for i := range strs {
		str := make([]rune, 0, 256)
		for {
			c, _, err := d.r.ReadRune()
			if err != nil {
				panic(err)
			}
			if c == 0 {
				strs[i] = string(str)
				break
			}
			str = append(str, c)
		}
	}
	d.r.Seek(o, io.SeekStart)
	return strs
}

func nameString(name [64]byte) string {
	return strings.TrimRight(string(name[:]), "\x00")
}

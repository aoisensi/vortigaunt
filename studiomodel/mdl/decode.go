package mdl

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
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
	defer d.r.Seek(o, io.SeekStart)
	if err := f(); err != nil {
		return err
	}
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

//lint:ignore U1000 will be used later
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

func (d *Decoder) quat48() (*[4]float32, error) {
	var v [3]int16
	if err := binary.Read(d.r, le, &v); err != nil {
		return nil, err
	}
	x := float32(v[0]) / 32768.0
	y := float32(v[1]) / 32768.0
	z := float32(v[2] & ^1 << 15) / 16384.0
	w := float32(math.Sqrt(float64(1.0 - x*x - y*y - z*z)))
	if v[2]&1<<15 != 0 {
		w = -w
	}
	return &[4]float32{x, y, z, w}, nil
}

func (d *Decoder) quat64() (*[4]float32, error) {
	var v uint64
	if err := binary.Read(d.r, le, &v); err != nil {
		return nil, err
	}
	const div = 1.0 / 1048576.5
	const bits = uint64(0b111111111111111111111) // 21 bits
	x := float32(int32(v>>00&bits)-1048576) * div
	y := float32(int32(v>>21&bits)-1048576) * div
	z := float32(int32(v>>42&bits)-1048576) * div
	w := float32(math.Sqrt(float64(1.0 - x*x - y*y - z*z)))
	if v&1<<63 != 0 {
		w = -w
	}
	return &[4]float32{x, y, z, w}, nil
}

func (d *Decoder) vec32() (*[3]float32, error) {
	var v uint32
	if err := binary.Read(d.r, le, &v); err != nil {
		return nil, err
	}
	bits := uint32(0b1111111111) // 10 bits
	x := int32(v>>00&bits) - 512
	y := int32(v>>10&bits) - 512
	z := int32(v>>20&bits) - 512
	e := v >> 30 & 0b11
	f := []float32{4.0, 16.0, 32.0, 64.0}[e]
	return &[3]float32{float32(x) / f, float32(y) / f, float32(z) / f}, nil
}

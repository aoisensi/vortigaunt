package vtx

import (
	"encoding/binary"
	"io"
)

type BodyPart struct {
	Header *BodyPartHeader
	Models []*Model
}
type BodyPartHeader struct {
	NumModels   int32
	ModelOffset int32
}

func (d *Decoder) decodeBodyParts(vtx *VTX) error {
	p, _ := d.r.Seek(int64(vtx.Header.BodyPartOffset), io.SeekStart)
	vtx.BodyParts = make([]*BodyPart, vtx.Header.NumBodyParts)
	for i := range vtx.BodyParts {
		header := new(BodyPartHeader)
		if err := binary.Read(d.r, le, header); err != nil {
			return err
		}
		bp := new(BodyPart)
		bp.Header = header
		vtx.BodyParts[i] = bp
	}
	for i, bp := range vtx.BodyParts {
		d.r.Seek(p+int64(i*8), io.SeekStart)
		if err := d.decodeModels(bp); err != nil {
			return err
		}
	}
	return nil
}

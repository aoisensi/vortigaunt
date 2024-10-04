package vtx

import (
	"encoding/binary"
	"io"
)

type Model struct {
	Header *ModelHeader
	LODs   []*LOD
}

type ModelHeader struct {
	NumLODs   int32
	LODOffset int32
}

func (d *Decoder) decodeModels(bp *BodyPart) error {
	p, _ := d.r.Seek(int64(bp.Header.ModelOffset), io.SeekCurrent)
	bp.Models = make([]*Model, bp.Header.NumModels)
	for i := range bp.Models {
		model := new(Model)
		header := new(ModelHeader)
		if err := binary.Read(d.r, le, header); err != nil {
			return err
		}
		model.Header = header
		bp.Models[i] = model
	}
	for i, model := range bp.Models {
		d.r.Seek(p+int64(i*8), io.SeekStart)
		if err := d.decodeLODs(model); err != nil {
			return err
		}
	}
	return nil
}

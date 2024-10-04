package vtx

import (
	"encoding/binary"
	"io"
)

type Strip struct {
	Header *StripHeader
}

type StripHeader struct {
	NumIndices            int32
	IndexOffset           int32
	NumVerts              int32
	VertOffset            int32
	NumBones              int16
	Flags                 byte
	NumBoneStateChanges   int32
	BoneStateChangeOffset int32
}

func (d *Decoder) decodeStrips(sg *StripGroup) error {
	d.r.Seek(int64(sg.Header.StripOffset), io.SeekCurrent)
	sg.Strips = make([]*Strip, sg.Header.NumStrips)
	for i := range sg.Strips {
		header := new(StripHeader)
		if err := binary.Read(d.r, le, header); err != nil {
			return err
		}
		sg.Strips[i] = &Strip{Header: header}
	}
	return nil
}

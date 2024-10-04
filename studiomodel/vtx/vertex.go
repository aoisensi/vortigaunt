package vtx

import (
	"encoding/binary"
	"io"
)

type Vertex struct {
	BoneWeightIndex      [3]uint8
	NumBones             uint8
	OriginalMeshVertexID uint16
	BoneID               [3]int8
}

func (d *Decoder) decodeVertexes(sg *StripGroup) error {
	d.r.Seek(int64(sg.Header.VertOffset), io.SeekCurrent)
	sg.Vertexes = make([]*Vertex, sg.Header.NumVerts)
	for i := range sg.Vertexes {
		vertex := new(Vertex)
		if err := binary.Read(d.r, le, vertex); err != nil {
			return err
		}
		sg.Vertexes[i] = vertex
	}
	return nil
}

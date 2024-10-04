package vtx

import (
	"encoding/binary"
	"io"
)

type StripGroup struct {
	Header *StripGroupHeader
	// Header2
	Vertexes []*Vertex
	Indices  []uint16
	Strips   []*Strip
}

type StripGroupHeader struct {
	NumVerts    int32
	VertOffset  int32
	NumIndices  int32
	IndexOffset int32
	NumStrips   int32
	StripOffset int32
	Flags       byte
}

func (d *Decoder) decodeStripGroups(mesh *Mesh) error {
	p, _ := d.r.Seek(int64(mesh.Header.StripGroupOffset), io.SeekCurrent)
	mesh.StripGroups = make([]*StripGroup, mesh.Header.NumStripGroups)
	for i := range mesh.StripGroups {
		sg := new(StripGroup)
		header := new(StripGroupHeader)
		if err := binary.Read(d.r, le, header); err != nil {
			return err
		}
		sg.Header = header
		mesh.StripGroups[i] = sg
	}
	for i, sg := range mesh.StripGroups {
		d.r.Seek(p+int64(i*25), io.SeekStart)
		if err := d.decodeVertexes(sg); err != nil {
			return err
		}
		sg.Indices = make([]uint16, sg.Header.NumIndices)
		d.r.Seek(p+int64(i*25), io.SeekStart)
		d.r.Seek(int64(sg.Header.IndexOffset), io.SeekCurrent)
		if err := binary.Read(d.r, le, sg.Indices); err != nil {
			return err
		}
		d.r.Seek(p+int64(i*25), io.SeekStart)
		if err := d.decodeStrips(sg); err != nil {
			return err
		}
	}
	return nil
}

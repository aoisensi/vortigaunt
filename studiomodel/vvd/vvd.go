package vvd

import (
	"encoding/binary"
	"io"
)

type VVD struct {
	Header   *Header
	Fixups   []*Fixup
	Vertexes []*Vertex // usually use LODsData
	Tangents [][4]float32
	LODsData [][]*Vertex
}

func (d *Decoder) decodeVVD() (*VVD, error) {
	var err error
	vvd := new(VVD)
	vvd.Header, err = d.decodeHeader()
	if err != nil {
		return nil, err
	}

	if _, err := d.r.Seek(int64(vvd.Header.FixupTableStart), io.SeekStart); err != nil {
		return nil, err
	}
	vvd.Fixups = make([]*Fixup, vvd.Header.NumFixups)
	for i := range vvd.Fixups {
		vvd.Fixups[i], err = d.decodeFixup()
		if err != nil {
			return nil, err
		}
	}

	if _, err := d.r.Seek(int64(vvd.Header.VertexDataStart), io.SeekStart); err != nil {
		return nil, err
	}
	vvd.Vertexes = make([]*Vertex, vvd.Header.NumLODVertexes[0])
	for i := range vvd.Vertexes {
		vvd.Vertexes[i], err = d.decodeVertex()
		if err != nil {
			return nil, err
		}
	}

	if _, err := d.r.Seek(int64(vvd.Header.TangentDataStart), io.SeekStart); err != nil {
		return nil, err
	}
	vvd.Tangents = make([][4]float32, vvd.Header.NumLODVertexes[0])
	err = binary.Read(d.r, binary.LittleEndian, &vvd.Tangents)
	if err != nil {
		return nil, err
	}

	if len(vvd.Fixups) > 0 {
		vvd.LODsData = make([][]*Vertex, vvd.Header.NumLODs)
		for lodID := 0; lodID < int(vvd.Header.NumLODs); lodID++ {
			offset := 0
			data := make([]*Vertex, vvd.Header.NumLODVertexes[lodID])
			for _, fixup := range vvd.Fixups {
				if fixup.LOD >= int32(lodID) {
					idx := int(fixup.SourceVertexID)
					cnt := int(fixup.NumVertexes)
					copy(data[offset:offset+cnt], vvd.Vertexes[idx:idx+cnt])
					offset += cnt
				}
			}
			vvd.LODsData[lodID] = data
		}
	} else {
		vvd.LODsData = [][]*Vertex{vvd.Vertexes}
	}
	return vvd, nil
}

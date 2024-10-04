package vtx

import (
	"encoding/binary"
	"io"
)

type LOD struct {
	Header *LODHeader
	Meshes []*Mesh
}

type LODHeader struct {
	NumMeshes   int32
	MeshOffset  int32
	SwitchPoint float32
}

func (d *Decoder) decodeLODs(model *Model) error {
	p, _ := d.r.Seek(int64(model.Header.LODOffset), io.SeekCurrent)
	model.LODs = make([]*LOD, model.Header.NumLODs)
	for i := range model.LODs {
		lod := new(LOD)
		header := new(LODHeader)
		if err := binary.Read(d.r, le, header); err != nil {
			return err
		}
		lod.Header = header
		model.LODs[i] = lod
	}
	for i, lod := range model.LODs {
		d.r.Seek(p+int64(i*12), io.SeekStart)
		if err := d.decodeMeshs(lod); err != nil {
			return err
		}
	}
	return nil
}

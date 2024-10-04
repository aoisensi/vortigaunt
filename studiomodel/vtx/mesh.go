package vtx

import (
	"encoding/binary"
	"io"
)

type Mesh struct {
	Header      *MeshHeader
	StripGroups []*StripGroup
}

type MeshHeader struct {
	NumStripGroups   int32
	StripGroupOffset int32
	Flags            byte
}

func (d *Decoder) decodeMeshs(lod *LOD) error {
	p, _ := d.r.Seek(int64(lod.Header.MeshOffset), io.SeekCurrent)
	lod.Meshes = make([]*Mesh, lod.Header.NumMeshes)
	for i := range lod.Meshes {
		mesh := new(Mesh)
		header := new(MeshHeader)
		if err := binary.Read(d.r, le, header); err != nil {
			return err
		}
		mesh.Header = header
		lod.Meshes[i] = mesh
	}
	for i, mesh := range lod.Meshes {
		d.r.Seek(p+int64(i*9), io.SeekStart)
		if err := d.decodeStripGroups(mesh); err != nil {
			return err
		}
	}
	return nil
}

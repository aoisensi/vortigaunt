package mdl

type Mesh struct {
	Header *MeshHeader
}

type MeshHeader struct {
	Material      int32
	ModelIndex    int32
	NumVertices   int32
	VertexOffset  int32
	NumFlexes     int32
	FlexIndex     int32
	MaterialType  int32
	MaterialParam int32
	MeshID        int32
	Center        [3]float32
	_             [9]int32 // mstudio_meshvertexdata_t
	_             [8]int32
}

func (d *Decoder) decodeMesh(model *Model) error {
	model.Meshes = make([]*Mesh, model.Header.NumMeshes)
	for i := range model.Meshes {
		header := new(MeshHeader)
		if err := d.read(header); err != nil {
			return err
		}
		mesh := new(Mesh)
		mesh.Header = header
		model.Meshes[i] = mesh
	}
	return nil
}

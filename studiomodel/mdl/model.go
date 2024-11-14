package mdl

// mstudiomodel_t
type Model struct {
	Header *ModelHeader
	Name   string
	Meshes []*Mesh
}

type ModelHeader struct {
	Name            [64]byte
	Type            int32
	BoundingRadius  float32
	NumMeshes       int32
	MeshIndex       int32
	NumVertices     int32
	VertexIndex     int32
	TangentsIndex   int32
	NumAttachments  int32
	AttachmentIndex int32
	NumEyeballs     int32
	EyeballIndex    int32
	_               [2]int32 // mstudio_modelvertexdata_t
	_               [8]int32
} // 148 0x9c bytes

func (d *Decoder) decodeModel(bp *BodyPart) error {
	bp.Models = make([]*Model, bp.Header.NumModels)
	for i := range bp.Models {
		header := new(ModelHeader)
		if err := d.read(header); err != nil {
			return err
		}
		model := new(Model)
		model.Header = header
		model.Name = nameString(model.Header.Name)

		err := d.ppush(
			header.MeshIndex-148,
			func() error { return d.decodeMesh(model) },
		)
		if err != nil {
			return err
		}

		bp.Models[i] = model
	}
	return nil
}

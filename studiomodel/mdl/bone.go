package mdl

type Bone struct {
	Header *BoneHeader
	Name   string
}

type BoneHeader struct {
	NameIndex        int32
	ParentID         int32
	BoneController   [6]int32
	Pos              [3]float32
	Quat             [4]float32
	Rot              [3]float32
	PosScale         [3]float32
	RotScale         [3]float32
	PosToBone        [3][4]float32
	QAlignment       [4]float32
	Flags            int32
	ProcType         int32
	ProcIndex        int32
	PhysicsBone      int32
	SurfacePropIndex int32
	Contents         int32
	_                [8]int32
}

func (d *Decoder) decodeBones(mdl *MDL) error {
	mdl.Bones = make([]*Bone, mdl.Header.BoneCount)
	for i := range mdl.Bones {
		header := new(BoneHeader)

		if err := d.read(header); err != nil {
			return err
		}
		bone := new(Bone)
		bone.Header = header
		bone.Name = d.readName(header.NameIndex - 216)

		mdl.Bones[i] = bone
	}
	return nil
}

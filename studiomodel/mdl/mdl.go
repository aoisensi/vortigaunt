package mdl

import "encoding/binary"

type MDL struct {
	Header *Header
	Name   string

	Bones       []*Bone
	Textures    []*Texture
	TextureDirs []string
	Skins       [][]*Texture
	AnimDescs   []*AnimDesc
	BodyParts   []*BodyPart
}

type Header struct {
	ID         int32
	Version    int32
	CheckSum   int32
	Name       [64]byte
	DataLength int32

	EyePosition   [3]float32
	IllumPosition [3]float32
	HullMin       [3]float32
	HullMax       [3]float32
	ViewBBMin     [3]float32
	ViewBBMax     [3]float32

	Flags MDLFlags

	BoneCount  int32
	BoneOffset int32

	BoneControllerCount     int32
	BoneControllerOffset    int32
	HitBoxCount             int32
	HitBoxOffset            int32
	LocalAnimCount          int32
	LocalAnimOffset         int32
	LocalSeqCount           int32
	LocalSeqOffset          int32
	ActivityListVersion     int32
	EeventsIndexed          int32
	TextureCount            int32
	TextureOffset           int32
	TextureDirCount         int32
	TextureDirOffset        int32
	SkinReferenceCount      int32
	SkinFamilyCount         int32
	SkinReferenceIndex      int32
	BodyPartCount           int32
	BodyPartOffset          int32
	AttachmentCount         int32
	AttachmentOffset        int32
	LocalNodeCount          int32
	LocalNodeIndex          int32
	LocalNodeNameIndex      int32
	FlexDescCount           int32
	FlexDescIndex           int32
	FlexControllerCount     int32
	FlexControllerIndex     int32
	FlexRulesCount          int32
	FlexRulesIndex          int32
	IKChainCount            int32
	IKChainIndex            int32
	MouthsCount             int32
	MouthsIndex             int32
	LocalPoseParamCount     int32
	LocalPoseParamIndex     int32
	SurfacePropIndex        int32
	KeyValueIndex           int32
	KeyValueCount           int32
	IKLockCount             int32
	IKLockIndex             int32
	Mass                    float32
	_                       int32
	IncludeModelCount       int32
	IncludeModelIndex       int32
	VirtualModel            int32
	AnimBlocksNameIndex     int32
	AnimBlocksCount         int32
	AnimBlocksIndex         int32
	_                       int32
	BoneTableNameIndex      int32
	_                       int32
	_                       int32
	DirectionalDotProduct   byte
	RootLOD                 byte
	NumAllowedRootLODs      byte
	_                       byte
	_                       int32
	FlexControllerUICount   int32
	FlexControllerUIIndex   int32
	VertAnimFixedPointScale float32
	_                       int32
	StudioHDR2Index         int32
	_                       int32
}

type MDLFlags uint32

func (f MDLFlags) IsStaticProp() bool {
	return f&1<<4 != 0 // https://github.com/ValveSoftware/source-sdk-2013/blob/0d8dceea4310fde5706b3ce1c70609d72a38efdf/sp/src/public/studio.h#L1965
}

func (d *Decoder) decodeMDL() (*MDL, error) {

	d.mdl = new(MDL)
	header := new(Header)
	err := d.ppush(0, func() error { return binary.Read(d.r, le, header) })
	if err != nil {
		return nil, err
	}
	d.mdl.Header = header
	d.mdl.Name = nameString(header.Name)

	// Bones
	err = d.ppush(
		d.mdl.Header.BoneOffset,
		func() error { return d.decodeBones(d.mdl) },
	)
	if err != nil {
		return nil, err
	}

	// Textures
	err = d.ppush(
		d.mdl.Header.TextureOffset,
		func() error { return d.decodeTexture(d.mdl) },
	)
	if err != nil {
		return nil, err
	}

	// TextureDirs
	d.mdl.TextureDirs = make([]string, d.mdl.Header.TextureDirCount)
	if err := d.ppush(0, func() error {
		ps := make([]int32, d.mdl.Header.TextureDirCount)
		if err := d.ppush(d.mdl.Header.TextureDirOffset, func() error {
			return d.read(ps)
		}); err != nil {
			return err
		}
		for i, p := range ps {
			d.mdl.TextureDirs[i] = d.readName(p)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	// Skins
	d.mdl.Skins = make([][]*Texture, header.SkinFamilyCount)
	for i := range d.mdl.Skins {
		d.mdl.Skins[i] = make([]*Texture, header.SkinReferenceCount)
	}
	err = d.ppush(header.SkinReferenceIndex, func() error {
		for i, ref := range d.mdl.Skins {
			for j := range ref {
				var id int16
				if err := d.read(&id); err != nil {
					return err
				}
				d.mdl.Skins[i][j] = d.mdl.Textures[id]
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// AnimDescs
	err = d.ppush(
		d.mdl.Header.LocalAnimOffset,
		func() error { return d.decodeAnimDesc(d.mdl) },
	)
	if err != nil {
		return nil, err
	}

	// BodyParts
	err = d.ppush(
		d.mdl.Header.BodyPartOffset,
		func() error { return d.decodeBodyParts(d.mdl) },
	)
	if err != nil {
		return nil, err
	}
	return d.mdl, nil
}

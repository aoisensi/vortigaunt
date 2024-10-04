package vtx

type VTX struct {
	Header    *Header
	BodyParts []*BodyPart
}

type Header struct {
	Version                       int32
	VertCacheSize                 int32
	MaxBonesPerStrip              uint16
	MaxBonesPerTri                uint16
	MaxBonesPerVert               int32
	CheckSum                      int32
	NumLODs                       int32
	MaterialReplacementListOffset int32
	NumBodyParts                  int32
	BodyPartOffset                int32
}

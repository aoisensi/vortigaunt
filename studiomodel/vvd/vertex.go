package vvd

import (
	"encoding/binary"

	"github.com/aoisensi/vortigaunt/studiomodel/internal"
)

type Vertex struct {
	BoneWeight BoneWeight
	Position   [3]float32
	Normal     [3]float32
	TexCoord   [2]float32
}

type BoneWeight struct {
	Weight   [internal.MaxNumBonesPerVert]float32
	Bone     [internal.MaxNumBonesPerVert]int8
	NumBones uint8
}

func (d *Decoder) decodeVertex() (*Vertex, error) {
	vertex := new(Vertex)
	return vertex, binary.Read(d.r, binary.LittleEndian, vertex)
}

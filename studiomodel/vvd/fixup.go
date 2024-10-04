package vvd

import "encoding/binary"

type Fixup struct {
	LOD            int32
	SourceVertexID int32
	NumVertexes    int32
}

func (d *Decoder) decodeFixup() (*Fixup, error) {
	fixup := new(Fixup)
	return fixup, binary.Read(d.r, binary.LittleEndian, fixup)
}

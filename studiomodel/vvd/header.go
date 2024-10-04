package vvd

import (
	"encoding/binary"
	"errors"

	"github.com/aoisensi/vortigaunt/studiomodel/internal"
)

const magic = 0x56534449

type Header struct {
	ID               int32
	Version          int32
	Checksum         int32
	NumLODs          int32
	NumLODVertexes   [internal.MaxNumLODs]int32
	NumFixups        int32
	FixupTableStart  int32
	VertexDataStart  int32
	TangentDataStart int32
}

func (d *Decoder) decodeHeader() (*Header, error) {
	header := new(Header)
	err := binary.Read(d.r, binary.LittleEndian, header)
	if err != nil {
		return nil, err
	}
	if header.ID != magic {
		return nil, errors.New("vvd: invalid id")
	}
	return header, nil
}

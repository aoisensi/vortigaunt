package studiomodel

import (
	"fmt"
	"os"
	"strings"

	"github.com/aoisensi/vortigaunt/studiomodel/mdl"
	"github.com/aoisensi/vortigaunt/studiomodel/vtx"
	"github.com/aoisensi/vortigaunt/studiomodel/vvd"
)

type Model struct {
	MDL *mdl.MDL
	VTX *vtx.VTX
	VVD *vvd.VVD
}

func LoadFromDisc(mdlName string) (*Model, error) {
	if !strings.HasSuffix(mdlName, ".mdl") {
		return nil, fmt.Errorf("studiomodel: mld filename must end with .mdl")
	}
	name := mdlName[:len(mdlName)-4]
	mdlF, err := os.Open(mdlName)
	if err != nil {
		return nil, err
	}
	defer mdlF.Close()
	model := new(Model)
	model.MDL, err = mdl.NewDecoder(mdlF).Decode()
	if err != nil {
		return nil, err
	}
	if len(model.MDL.BodyParts) == 0 {
		return model, nil
	}

	vtxF, err := os.Open(name + ".dx90.vtx")
	if err != nil {
		return nil, err
	}
	defer vtxF.Close()
	model.VTX, err = vtx.NewDecoder(vtxF).Decode()
	if err != nil {
		return nil, err
	}

	vvdF, err := os.Open(name + ".vvd")
	if err != nil {
		return nil, err
	}
	defer vvdF.Close()
	model.VVD, err = vvd.NewDecoder(vvdF).Decode()
	if err != nil {
		return nil, err
	}
	return model, nil
}

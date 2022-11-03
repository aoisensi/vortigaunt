package main

import (
	"io"
	"log"
	"os"

	"github.com/galaco/studiomodel"
	"github.com/galaco/studiomodel/mdl"
	"github.com/galaco/studiomodel/phy"
	"github.com/galaco/studiomodel/vtx"
	"github.com/galaco/studiomodel/vvd"
)

func load(mdlName string) *studiomodel.StudioModel {
	name := mdlName[:len(mdlName)-4]
	model := studiomodel.NewStudioModel(mdlName)
	model.Mdl = loadSingle(mdlName, mdl.ReadFromStream)
	model.Vtx = loadSingle(name+".dx90.vtx", vtx.ReadFromStream)
	model.Vvd = loadSingle(name+".vvd", vvd.ReadFromStream)
	model.Phy = loadSingle(name+".phy", phy.ReadFromStream)

	return model
}

func loadSingle[T any](name string, reader func(io.Reader) (*T, error)) *T {
	f, err := os.Open(name)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer f.Close()
	v, err := reader(f)
	if err != nil {
		log.Println(err)
		return nil
	}
	return v
}

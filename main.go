package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/aoisensi/go-fbx/pkg/fbx75"
)

func main() {
	for _, fname := range flag.Args() {
		if !strings.HasSuffix(fname, ".mdl") {
			log.Printf("%v is not .mdl file.", fname)
			continue
		}
		mdl := load(fname)
		splitedName := strings.Split(string(mdl.Mdl.Header.Name[:]), "/")
		name := splitedName[len(splitedName)-1]
		name = name[0 : len(name)-4]

		fbx := fbx75.NewFBX()

		fbx.GlobalSettings.UpAxis = fbx75.AxisYP
		fbx.GlobalSettings.FrontAxis = fbx75.AxisXP

		geometry := &fbx75.ObjectGeometry{
			ID:   rand.Int63(),
			Name: name,
		}
		{
			vertices := make([]float64, len(mdl.Vvd.Vertices)*3)
			for i, v := range mdl.Vvd.Vertices {
				for jm, jf := range []int{2, 0, 1} {
					vertices[i*3+jf] = float64(v.Position[jm])
				}
			}
			sg := mdl.Vtx.BodyParts[0].Models[0].LODS[0].Meshes[0].StripGroups[0]
			indeces := make([]int32, len(sg.Indices))
			polygonNum := len(sg.Indices) / 3
			for i := 0; i < polygonNum; i++ {
				indeces[i*3+0] = int32(sg.Vertexes[sg.Indices[i*3+0]].OriginalMeshVertexID)
				indeces[i*3+1] = int32(sg.Vertexes[sg.Indices[i*3+2]].OriginalMeshVertexID)
				indeces[i*3+2] = ^int32(sg.Vertexes[sg.Indices[i*3+1]].OriginalMeshVertexID)
			}
			geometry.Vertices = vertices
			geometry.PolygonVertexIndex = indeces
		}
		model := &fbx75.ObjectModel{
			ID:   rand.Int63(),
			Name: name,
		}
		fbx.Objects.Objects = append(
			fbx.Objects.Objects,
			geometry,
			model,
		)

		fbx.Connections.Cs = []fbx75.C{
			{model.ID, 0},
			{geometry.ID, model.ID},
		}

		f, err := os.Create(fname[:len(fname)-4] + ".fbx")
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		if _, err := fbx.WriteTo(f); err != nil {
			log.Println(err)
		}
	}
}

func init() {
	flag.Parse()
}

package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/aoisensi/go-fbx/pkg/fbx75"
	"github.com/go-gl/mathgl/mgl32"
)

func main() {
	for _, fname := range flag.Args() {
		if !strings.HasSuffix(fname, ".mdl") {
			log.Printf("%v is not .mdl file.", fname)
			continue
		}
		mdl := load(fname)
		textureName := mdl.Mdl.TextureNames[0]
		longName := strings.TrimRight(string(mdl.Mdl.Header.Name[:]), "\x00")
		splitedName := strings.Split(longName, "/")
		name := splitedName[len(splitedName)-1]
		name = name[0 : len(name)-4]

		fbx := fbx75.NewFBX()

		material := &fbx75.Material{
			ID:   rand.Int63(),
			Name: textureName,
		}

		geometry := &fbx75.Geometry{
			ID:   rand.Int63(),
			Name: name,
		}
		{
			verticesNum := len(mdl.Vvd.Vertices)
			vertices := make([]float64, 0, verticesNum*3)
			verticesMap := make(map[mgl32.Vec3]int, verticesNum)
			indexIndex := make([]int, 0, verticesNum)
			for _, v := range mdl.Vvd.Vertices {
				if index, ok := verticesMap[v.Position]; ok {
					indexIndex = append(indexIndex, index)
				} else {
					index := len(verticesMap)
					indexIndex = append(indexIndex, index)
					verticesMap[v.Position] = index
					vp := src2fbxXYZ(v.Position)
					vertices = append(vertices, vp[:]...)
				}
			}
			geometry.Vertices = vertices

			sg := mdl.Vtx.BodyParts[0].Models[0].LODS[0].Meshes[0].StripGroups[0]
			indecesNum := len(sg.Indices)
			indeces := make([]int32, indecesNum)
			uvs := make([]float64, 0, indecesNum*2)
			polygonNum := indecesNum / 3
			for i := 0; i < polygonNum; i++ {
				for j := 0; j < 3; j++ {
					v := sg.Vertexes[sg.Indices[i*3+j]]
					vid := v.OriginalMeshVertexID

					index := int32(indexIndex[vid])
					if j == 2 {
						index = ^index
					}
					indeces[i*3+j] = index

					vv := mdl.Vvd.Vertices[vid]
					uv := src2fbxUV(vv.UVs)
					uvs = append(uvs, uv[:]...)
				}
			}
			geometry.PolygonVertexIndex = indeces
			geometry.LayerElementUV = &fbx75.LayerElementUV{
				UV: uvs,
			}
			geometry.LayerElementMaterial = &fbx75.LayerElementMaterial{}
			geometry.Layer = &fbx75.Layer{
				LayerElements: []*fbx75.LayerElement{
					{
						Type:       "LayerElementUV",
						TypedIndex: 0,
					},
					{
						Type:       "LayerElementMaterial",
						TypedIndex: 0,
					},
				},
			}
		}
		model := &fbx75.Model{
			ID:   rand.Int63(),
			Name: name,
		}
		fbx.Objects.Objects = append(
			fbx.Objects.Objects,
			geometry,
			model,
			material,
		)

		fbx.Connections.Cs = []fbx75.C{
			{model.ID, 0},
			{geometry.ID, model.ID},
			{material.ID, model.ID},
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

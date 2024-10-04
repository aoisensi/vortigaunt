package vortigaunt

import (
	"path/filepath"
	"strings"

	"github.com/aoisensi/vortigaunt/studiomodel"
	"github.com/aoisensi/vortigaunt/vmath"
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
)

func convert(n string) error {
	m, err := studiomodel.LoadFromDisc(n)
	if err != nil {
		return err
	}
	document := gltf.NewDocument()

	document.Asset.Generator = "vortigaunt"
	scene := document.Scenes[0]
	scene.Name = strings.TrimRight(filepath.Base(m.MDL.Name), ".mdl")

	// Bones
	if !m.MDL.Header.Flags.IsStaticProp() {
		boneNodes := make([]*gltf.Node, len(m.MDL.Bones))
		for i, mdlBone := range m.MDL.Bones {
			node := &gltf.Node{
				Name:        mdlBone.Name,
				Translation: vmath.FtoD3(vmath.VecMulScalar(vmath.VecToGL(mdlBone.Header.Pos), float32(flagScale))),
				Rotation:    vmath.FtoD4(vmath.QuatToGL(mdlBone.Header.Quat)),
				Children:    make([]int, 0, 4),
			}
			if mdlBone.Header.ParentID >= 0 {
				boneNodes[mdlBone.Header.ParentID].Children = append(boneNodes[mdlBone.Header.ParentID].Children, len(document.Nodes))
			}
			boneNodes[i] = node
			document.Nodes = append(document.Nodes, node)
		}
	}

	if len(m.MDL.BodyParts) > 0 {
		positions := make([][3]float32, 0, len(m.VVD.Vertexes))
		normals := make([][3]float32, 0, len(m.VVD.Vertexes))
		texcoords := make([][2]float32, 0, len(m.VVD.Vertexes))

		if len(m.VVD.Fixups) > 0 {
			for _, fixup := range m.VVD.Fixups {
				vid := fixup.SourceVertexID
				num := fixup.NumVertexes
				for i := vid; i < vid+num; i++ {
					p := vmath.VecToGL(m.VVD.Vertexes[i].Position)
					p = vmath.VecMulScalar(p, float32(flagScale))
					positions = append(positions, p)

					n := vmath.VecToGL(m.VVD.Vertexes[i].Normal)
					normals = append(normals, n)

					t := m.VVD.Vertexes[i].TexCoord
					texcoords = append(texcoords, [2]float32{t[0], t[1]})
				}
			}
		} else {
			for _, v := range m.VVD.Vertexes {
				p := vmath.VecToGL(v.Position)
				p = vmath.VecMulScalar(p, float32(flagScale))
				positions = append(positions, p)

				n := vmath.VecToGL(v.Normal)
				normals = append(normals, n)

				t := v.TexCoord
				texcoords = append(texcoords, [2]float32{t[0], t[1]})
			}
		}
		positionID := modeler.WritePosition(document, positions)
		normalID := modeler.WriteNormal(document, normals)
		texcoordID := modeler.WriteTextureCoord(document, texcoords)

		for bpID, mdlBP := range m.MDL.BodyParts {
			vtxBP := m.VTX.BodyParts[bpID]
			node := &gltf.Node{
				Name: mdlBP.Name,
			}
			mesh := &gltf.Mesh{
				Name: mdlBP.Name,
			}
			for modelID, mdlModel := range mdlBP.Models {
				vtxModel := vtxBP.Models[modelID]

				for meshID, mdlMesh := range mdlModel.Meshes {
					primitive := &gltf.Primitive{
						Attributes: map[string]int{
							"POSITION":   positionID,
							"NORMAL":     normalID,
							"TEXCOORD_0": texcoordID,
						},
					}
					vtxMesh := vtxModel.LODs[0].Meshes[meshID]
					indices := make([]uint16, 0, mdlMesh.Header.NumVertices)
					for _, vtxSG := range vtxMesh.StripGroups {
						for _, vtxStrip := range vtxSG.Strips {
							for i := 0; i < int(vtxStrip.Header.NumIndices); i += 3 {
								for _, j := range []int{0, 2, 1} {
									idx1 := i + j + int(vtxStrip.Header.IndexOffset)
									idx2 := vtxSG.Indices[idx1]
									vertex := vtxSG.Vertexes[idx2]
									idx3 := vertex.OriginalMeshVertexID
									idx4 := int(mdlMesh.Header.VertexOffset) + int(idx3)
									index := idx4 + int(mdlModel.Header.VertexIndex)/48
									indices = append(indices, uint16(index))
								}
							}
						}
					}
					primitive.Indices = gltf.Index(modeler.WriteIndices(document, indices))
					mesh.Primitives = append(mesh.Primitives, primitive)
				}
			}
			document.Meshes = append(document.Meshes, mesh)
			node.Mesh = gltf.Index(len(document.Meshes) - 1)
			document.Nodes = append(document.Nodes, node)
			scene.Nodes = append(scene.Nodes, len(document.Nodes)-1)
		}
	}

	gltfName := strings.TrimSuffix(n, ".mdl") + ".gltf"
	gltf.Save(document, gltfName)
	return nil
}

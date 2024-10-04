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

	if len(m.MDL.BodyParts) > 0 {
		positions := make([][3]float32, 0, len(m.VVD.Vertexes))

		if len(m.VVD.Fixups) > 0 {
			for _, fixup := range m.VVD.Fixups {
				vid := fixup.SourceVertexID
				num := fixup.NumVertexes
				for i := vid; i < vid+num; i++ {
					p := vmath.Vec3ToGL(m.VVD.Vertexes[i].Position)
					positions = append(positions, p)
				}
			}
		} else {
			for _, v := range m.VVD.Vertexes {
				p := vmath.Vec3ToGL(v.Position)
				positions = append(positions, p)
			}
		}
		posID := modeler.WritePosition(document, positions)

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
							"POSITION": posID,
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

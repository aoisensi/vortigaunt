package vortigaunt

import (
	"fmt"
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

	flagScale := float32(flagScale)

	var skinIndex *int

	// Bones
	if !m.MDL.Header.Flags.IsStaticProp() {
		boneNodes := make([]*gltf.Node, len(m.MDL.Bones))
		skin := &gltf.Skin{}
		inverseBindMatrices := make([][4][4]float32, 0, len(m.MDL.Bones))
		for i, mdlBone := range m.MDL.Bones {
			translation := vmath.VecMulScalar(vmath.VecToGL(mdlBone.Header.Pos), flagScale)
			rotation := vmath.QuatToGL(mdlBone.Header.Quat)
			node := &gltf.Node{
				Name:        mdlBone.Name,
				Translation: vmath.FtoD3(translation),
				Rotation:    vmath.FtoD4(rotation),
				Children:    make([]int, 0, 4),
			}
			parentIBM := vmath.IdentityMat()
			if mdlBone.Header.ParentID >= 0 {
				parentIBM = inverseBindMatrices[mdlBone.Header.ParentID]
			}
			rotationMatrix := vmath.MakeRotateMat(rotation)
			translationMatrix := vmath.MakeTranslateMat(translation)
			localMatrix := vmath.MultiplyMat(translationMatrix, rotationMatrix)

			combinedMatrix := vmath.MultiplyMat(vmath.InverseMat(parentIBM), localMatrix)

			ibm := vmath.RoundMat(vmath.InverseMat(combinedMatrix))
			fmt.Println(ibm[0][3], ibm[1][3], ibm[2][3], ibm[3][3])
			inverseBindMatrices = append(inverseBindMatrices, ibm)

			if mdlBone.Header.ParentID >= 0 {
				boneNodes[mdlBone.Header.ParentID].Children = append(boneNodes[mdlBone.Header.ParentID].Children, len(document.Nodes))
			}
			boneNodes[i] = node
			document.Nodes = append(document.Nodes, node)
			if mdlBone.Header.ParentID < 0 {
				scene.Nodes = append(scene.Nodes, len(document.Nodes)-1)
			}

			skin.Joints = append(skin.Joints, len(document.Nodes)-1)
		}
		ibmID := modeler.WriteInverseBindMatrices(document, inverseBindMatrices)
		{ // TODO: Remove this block when glTF loader supports TargetNone
			bv := *document.Accessors[ibmID].BufferView
			document.BufferViews[bv].Target = gltf.TargetNone
		}
		skin.InverseBindMatrices = gltf.Index(ibmID)
		document.Skins = append(document.Skins, skin)
		skinIndex = gltf.Index(0)
	}

	if len(m.MDL.BodyParts) > 0 {
		positions := make([][3]float32, 0, len(m.VVD.Vertexes))
		normals := make([][3]float32, 0, len(m.VVD.Vertexes))
		texcoords := make([][2]float32, 0, len(m.VVD.Vertexes))
		var joints [][4]uint16
		var weights [][4]float32
		if skinIndex != nil {
			joints = make([][4]uint16, 0, len(m.VVD.Vertexes)*4)
			weights = make([][4]float32, 0, len(m.VVD.Vertexes)*4)
		}

		if len(m.VVD.Fixups) > 0 {
			for _, fixup := range m.VVD.Fixups {
				vid := fixup.SourceVertexID
				num := fixup.NumVertexes
				for i := vid; i < vid+num; i++ {
					v := m.VVD.Vertexes[i]
					p := vmath.VecToGL(v.Position)
					p = vmath.VecMulScalar(p, float32(flagScale))
					positions = append(positions, p)

					n := vmath.VecToGL(v.Normal)
					normals = append(normals, n)

					t := v.TexCoord
					texcoords = append(texcoords, [2]float32{t[0], t[1]})

					if skinIndex != nil {
						joint := [4]uint16{}
						weight := [4]float32{}
						for j := 0; j < int(v.BoneWeight.NumBones); j++ {
							joint[j] = uint16(v.BoneWeight.Bone[j])
							weight[j] = v.BoneWeight.Weight[j]
						}
						joints = append(joints, joint)
						weights = append(weights, weight)
					}
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

				if skinIndex != nil {
					joint := [4]uint16{}
					weight := [4]float32{}
					for j := 0; j < int(v.BoneWeight.NumBones); j++ {
						joint[j] = uint16(v.BoneWeight.Bone[j])
						weight[j] = v.BoneWeight.Weight[j]
					}
					joints = append(joints, joint)
					weights = append(weights, weight)
				}
			}
		}
		positionID := modeler.WritePosition(document, positions)
		normalID := modeler.WriteNormal(document, normals)
		texcoordID := modeler.WriteTextureCoord(document, texcoords)
		var jointID, weightID int
		if skinIndex != nil {
			jointID = modeler.WriteJoints(document, joints)
			weightID = modeler.WriteWeights(document, weights)
		}

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
					attributes := map[string]int{
						"POSITION":   positionID,
						"NORMAL":     normalID,
						"TEXCOORD_0": texcoordID,
					}
					if skinIndex != nil {
						attributes["JOINTS_0"] = jointID
						attributes["WEIGHTS_0"] = weightID
					}
					primitive := &gltf.Primitive{
						Attributes: attributes,
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
			node.Skin = skinIndex
			document.Nodes = append(document.Nodes, node)
			scene.Nodes = append(scene.Nodes, len(document.Nodes)-1)
		}
	}

	gltfName := strings.TrimSuffix(n, ".mdl") + ".gltf"
	gltf.Save(document, gltfName)
	return nil
}

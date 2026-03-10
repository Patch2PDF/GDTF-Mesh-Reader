package FileHandlers

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"

	Types "github.com/Patch2PDF/GDTF-Mesh-Reader/v2/pkg/MeshTypes"
	"github.com/qmuntal/gltf"
)

type gltfNode struct {
	matrix Types.Matrix
	mesh   *gltf.Mesh
}

var conversion = Types.Matrix{ // coordinate system conversion matrix
	X00: 1, X01: 0, X02: 0, X03: 0,
	X10: 0, X11: 0, X12: -1, X13: 0,
	X20: 0, X21: 1, X22: 0, X23: 0,
	X30: 0, X31: 0, X32: 0, X33: 1,
}

func LoadGLTF(file io.Reader, desiredSize *Types.Vector) (*Types.Mesh, error) {
	var doc gltf.Document
	gltf.NewDecoder(file).Decode(&doc)

	var meshes *Types.Mesh = &Types.Mesh{}

	transformationMatrices := map[int]Types.Matrix{}

	gltfNodes := map[int]gltfNode{}

	for _, scene := range doc.Scenes {
		for _, node := range scene.Nodes {
			handleGLTFNode(gltfNodes, doc, node, transformationMatrices, Types.IdentityMatrix())
		}
	}

	// TODO: refine this as pure rotation breaks bounding box -> possibly rotate all 8 outer points of bounding box and redetermine
	// calculate outer dimensions
	min := Types.Vector{X: math.Inf(1), Y: math.Inf(1), Z: math.Inf(1)}
	max := Types.Vector{X: math.Inf(-1), Y: math.Inf(-1), Z: math.Inf(-1)}
	for node_index, node := range gltfNodes {
		m := node.mesh
		for _, p := range m.Primitives {
			// contains Min and Max attr (for dimension calc)
			posAccessor := doc.Accessors[p.Attributes[gltf.POSITION]]
			// do transformation before getting the outer dimensions
			tempMin := transformationMatrices[node_index].MulPosition(Types.Vector{
				X: posAccessor.Min[0],
				Y: posAccessor.Min[1],
				Z: posAccessor.Min[2],
			})
			tempMax := transformationMatrices[node_index].MulPosition(Types.Vector{
				X: posAccessor.Max[0],
				Y: posAccessor.Max[1],
				Z: posAccessor.Max[2],
			})
			min = min.Min(tempMin)
			max = max.Max(tempMax)
		}
	}

	scaling := Types.Vector{X: 1, Y: 1, Z: 1}
	if desiredSize != nil {
		scaling = desiredSize.Div(max.Sub(min))
		// following is a temporary fix, TODO: remove once bounding box is refined
		scaling.X = math.Abs(scaling.X)
		scaling.Y = math.Abs(scaling.Y)
		scaling.Z = math.Abs(scaling.Z)
	}

	for _, node := range gltfNodes {
		m := node.mesh
		determinant := node.matrix.Determinant()
		for _, p := range m.Primitives {
			posAccessor := doc.Accessors[p.Attributes[gltf.POSITION]]
			positions, err := gltfVec3(&doc, posAccessor, node.matrix, scaling)
			if err != nil {
				return nil, err
			}

			var normals []Types.Vector
			if nIdx, ok := p.Attributes[gltf.NORMAL]; ok {
				normalAccessor := doc.Accessors[nIdx]
				normals, err = gltfVec3(&doc, normalAccessor, node.matrix, Types.Vector{X: 1, Y: 1, Z: 1})
				if err != nil {
					return nil, err
				}
			}

			indexAccessor := doc.Accessors[*p.Indices]
			indices, err := gltfIndices(&doc, indexAccessor)
			if err != nil {
				return nil, err
			}

			var mesh Types.Mesh
			if p.Mode == gltf.PrimitiveTriangles {
				for i := 0; i < len(indices); i += 3 {
					var n0 *Types.Vector = nil
					var n1 *Types.Vector = nil
					var n2 *Types.Vector = nil
					if len(normals) > 0 {
						n0 = &normals[indices[i+0]]
						n1 = &normals[indices[i+1]]
						n2 = &normals[indices[i+2]]
					}

					v0 := positions[indices[i+0]].ToVertex(n0)
					v1 := positions[indices[i+1]].ToVertex(n1)
					v2 := positions[indices[i+2]].ToVertex(n2)

					if determinant < 0 {
						mesh.AddTriangle(Types.Triangle{V0: v0, V1: v2, V2: v1})
					} else {
						mesh.AddTriangle(Types.Triangle{V0: v0, V1: v1, V2: v2})
					}
				}

				meshes.Add(&mesh)
			}
		}
	}
	return meshes, nil
}

func handleGLTFNode(nodes map[int]gltfNode, doc gltf.Document, node_id int, transformationMatrices map[int]Types.Matrix, parentMatrix Types.Matrix) {
	node := doc.Nodes[node_id]
	gltf_matrix := node.MatrixOrDefault()
	var matrix Types.Matrix = parentMatrix
	if gltf_matrix != gltf.DefaultMatrix {
		matrix = matrix.Mul(Types.Matrix{
			X00: gltf_matrix[0], X01: gltf_matrix[4], X02: gltf_matrix[8], X03: gltf_matrix[12],
			X10: gltf_matrix[1], X11: gltf_matrix[5], X12: gltf_matrix[9], X13: gltf_matrix[13],
			X20: gltf_matrix[2], X21: gltf_matrix[6], X22: gltf_matrix[10], X23: gltf_matrix[14],
			X30: gltf_matrix[3], X31: gltf_matrix[7], X32: gltf_matrix[11], X33: gltf_matrix[15],
		})
	} else {
		matrix = matrix.Mul(gltfParseScaleRotationTranslation(node.RotationOrDefault(), node.ScaleOrDefault(), node.TranslationOrDefault()))
	}
	if node.Mesh != nil {
		world_matrix := conversion.Mul(matrix)
		transformationMatrices[node_id] = world_matrix
		nodes[node_id] = gltfNode{
			matrix: world_matrix,
			mesh:   doc.Meshes[*node.Mesh],
		}
	}
	for _, child_node := range node.Children {
		handleGLTFNode(nodes, doc, child_node, transformationMatrices, matrix)
	}
}

func gltfVec3(doc *gltf.Document, acc *gltf.Accessor, transformationMatrix Types.Matrix, scaling Types.Vector) ([]Types.Vector, error) {
	bufView := doc.BufferViews[*acc.BufferView]
	buffer := doc.Buffers[bufView.Buffer]

	start := int(bufView.ByteOffset + acc.ByteOffset)
	end := start + acc.Count*12 // 3 floats * 4 bytes
	raw := buffer.Data[start:end]

	vectors := make([]Types.Vector, acc.Count)
	for i := 0; i < acc.Count; i++ {
		base := i * 12
		vec := Types.Vector{
			X: float64(math.Float32frombits(binary.LittleEndian.Uint32(raw[base+0:]))),
			Y: float64(math.Float32frombits(binary.LittleEndian.Uint32(raw[base+4:]))),
			Z: float64(math.Float32frombits(binary.LittleEndian.Uint32(raw[base+8:]))),
		}
		scaled := transformationMatrix.MulPosition(vec).Mult(scaling)
		vectors[i] = scaled // vec.Mult(scaling)
	}

	return vectors, nil
}

func gltfIndices(doc *gltf.Document, acc *gltf.Accessor) ([]int, error) {
	bufView := doc.BufferViews[*acc.BufferView]
	buffer := doc.Buffers[bufView.Buffer]

	start := int(bufView.ByteOffset + acc.ByteOffset)
	componentSize := acc.ComponentType.ByteSize()
	end := start + acc.Count*componentSize

	raw := buffer.Data[start:end]

	out := make([]int, acc.Count)

	switch acc.ComponentType {
	case gltf.ComponentUshort:
		for i := 0; i < acc.Count; i++ {
			out[i] = int(binary.LittleEndian.Uint16(raw[i*2:]))
		}
	case gltf.ComponentUint:
		for i := 0; i < acc.Count; i++ {
			out[i] = int(binary.LittleEndian.Uint32(raw[i*4:]))
		}
	default:
		return nil, fmt.Errorf("unsupported index type")
	}

	return out, nil
}

func normalizeQuaternion(q [4]float64) [4]float64 {
	// Calculate mag squared
	magSq := q[0]*q[0] + q[1]*q[1] + q[2]*q[2] + q[3]*q[3]

	// Check for near-zero
	if magSq < 1e-12 {
		return [4]float64{0, 0, 0, 1}
	}

	mag := math.Sqrt(magSq)
	return [4]float64{q[0] / mag, q[1] / mag, q[2] / mag, q[3] / mag}
}

func gltfParseScaleRotationTranslation(rotation [4]float64, scale [3]float64, translation [3]float64) Types.Matrix {
	rotation = normalizeQuaternion(rotation)
	return Types.Matrix{
		X00: (1 - 2*(rotation[1]*rotation[1]+rotation[2]*rotation[2])) * scale[0], X01: 2 * (rotation[0]*rotation[1] - rotation[3]*rotation[2]) * scale[1], X02: 2 * (rotation[0]*rotation[2] + rotation[3]*rotation[1]) * scale[2], X03: translation[0],
		X10: 2 * (rotation[0]*rotation[1] + rotation[3]*rotation[2]) * scale[0], X11: (1 - 2*(rotation[0]*rotation[0]+rotation[2]*rotation[2])) * scale[1], X12: 2 * (rotation[1]*rotation[2] - rotation[3]*rotation[0]) * scale[2], X13: translation[1],
		X20: 2 * (rotation[0]*rotation[2] - rotation[3]*rotation[1]) * scale[0], X21: 2 * (rotation[1]*rotation[2] + rotation[3]*rotation[0]) * scale[1], X22: (1 - 2*(rotation[0]*rotation[0]+rotation[1]*rotation[1])) * scale[2], X23: translation[2],
		X30: 0, X31: 0, X32: 0, X33: 1,
	}
}

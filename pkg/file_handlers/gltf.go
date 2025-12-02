package FileHandlers

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"

	Types "github.com/Patch2PDF/GDTF-Mesh-Reader/pkg/MeshTypes"
	"github.com/qmuntal/gltf"
)

func LoadGLTF(file io.Reader, desiredSize *Types.Vector) (*Types.Mesh, error) {
	var doc gltf.Document
	gltf.NewDecoder(file).Decode(&doc)

	var meshes *Types.Mesh = &Types.Mesh{}

	transformationMatrices := map[int]Types.Matrix{}

	for _, node := range doc.Nodes {
		if node.Mesh == nil {
			continue
		}
		matrix := node.MatrixOrDefault()
		transformationMatrices[*node.Mesh] = Types.Matrix{
			X00: matrix[0], X01: matrix[4], X02: matrix[8], X03: matrix[12],
			X10: matrix[1], X11: matrix[5], X12: matrix[9], X13: matrix[13],
			X20: matrix[2], X21: matrix[6], X22: matrix[10], X23: matrix[14],
			X30: matrix[3], X31: matrix[7], X32: matrix[11], X33: matrix[15],
		}
	}

	// calculate outer dimensions
	min := Types.Vector{}
	max := Types.Vector{}
	for meshindex, m := range doc.Meshes {
		for _, p := range m.Primitives {
			// contains Min and Max attr (for dimension calc)
			posAccessor := doc.Accessors[p.Attributes[gltf.POSITION]]
			// do transformation before getting the outer dimensions
			tempMin := transformationMatrices[meshindex].MulPosition(Types.Vector{
				X: posAccessor.Min[0],
				Y: posAccessor.Min[2],
				Z: posAccessor.Min[1],
			})
			tempMax := transformationMatrices[meshindex].MulPosition(Types.Vector{
				X: posAccessor.Max[0],
				Y: posAccessor.Max[2],
				Z: posAccessor.Max[1],
			})
			min = min.Min(&tempMin)
			max = max.Max(&tempMax)
		}
	}

	scaling := Types.Vector{X: 1, Y: 1, Z: 1}
	if desiredSize != nil {
		scaling = desiredSize.Div(max.Sub(min))
	}

	for meshindex, m := range doc.Meshes {
		for _, p := range m.Primitives {
			posAccessor := doc.Accessors[p.Attributes[gltf.POSITION]]
			positions, err := gltfVec3(&doc, posAccessor, transformationMatrices[meshindex], scaling)
			if err != nil {
				return nil, err
			}

			var normals []Types.Vector
			if nIdx, ok := p.Attributes[gltf.NORMAL]; ok {
				normalAccessor := doc.Accessors[nIdx]
				normals, err = gltfVec3(&doc, normalAccessor, transformationMatrices[meshindex], Types.Vector{X: 1, Y: 1, Z: 1})
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

					mesh.AddTriangle(&Types.Triangle{V0: v0, V1: v1, V2: v2})
				}

				meshes.Add(&mesh)
			}
		}
	}
	return meshes, nil
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
		// axes inverted to convert to correct coordinate system
		vec := Types.Vector{
			X: float64(math.Float32frombits(binary.LittleEndian.Uint32(raw[base+0:]))),
			Y: float64(math.Float32frombits(binary.LittleEndian.Uint32(raw[base+4:]))),
			Z: float64(math.Float32frombits(binary.LittleEndian.Uint32(raw[base+8:]))),
		}
		transformed := transformationMatrix.MulPosition(vec)
		scaled := Types.Vector{
			X: transformed.X,
			Y: -transformed.Z,
			Z: transformed.Y,
		}.Mult(scaling)
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

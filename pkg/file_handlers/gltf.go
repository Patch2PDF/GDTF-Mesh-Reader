package FileHandlers

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"

	Types "github.com/Patch2PDF/GDTF-Mesh-Reader/pkg/MeshTypes"
	"github.com/qmuntal/gltf"
)

// Reads content of a '.gltf' or '.glb' file and outputs the mesh it contains
func LoadGLTF(file io.Reader, desiredSize Types.Vector) ([]*Types.Mesh, error) {
	var doc gltf.Document
	gltf.NewDecoder(file).Decode(&doc)

	var meshes []*Types.Mesh

	for _, m := range doc.Meshes {
		for _, p := range m.Primitives {
			// contains Min and Max attr (for dimension calc)
			posAccessor := doc.Accessors[p.Attributes[gltf.POSITION]]
			scaling := Types.Vector{
				// axes inverted to convert to correct coordinate system
				X: desiredSize.X / (posAccessor.Max[0] - posAccessor.Min[0]),
				Y: desiredSize.Y / (posAccessor.Max[2] - posAccessor.Min[2]),
				Z: desiredSize.Z / (posAccessor.Max[1] - posAccessor.Min[1]),
			}
			positions, err := gltfVec3(&doc, posAccessor, scaling)
			if err != nil {
				return nil, err
			}

			var normals []Types.Vector
			if nIdx, ok := p.Attributes[gltf.NORMAL]; ok {
				normalAccessor := doc.Accessors[nIdx]
				normals, err = gltfVec3(&doc, normalAccessor, Types.Vector{X: 1, Y: 1, Z: 1})
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

				meshes = append(meshes, &mesh)
			}
		}
	}
	return meshes, nil
}

func gltfVec3(doc *gltf.Document, acc *gltf.Accessor, scaling Types.Vector) ([]Types.Vector, error) {
	bufView := doc.BufferViews[*acc.BufferView]
	buffer := doc.Buffers[bufView.Buffer]

	start := int(bufView.ByteOffset + acc.ByteOffset)
	end := start + acc.Count*12 // 3 floats * 4 bytes
	raw := buffer.Data[start:end]

	vectors := make([]Types.Vector, acc.Count)
	for i := 0; i < acc.Count; i++ {
		base := i * 12
		// axes inverted to convert to correct coordinate system
		x := math.Float32frombits(binary.LittleEndian.Uint32(raw[base+0:]))
		y := -(math.Float32frombits(binary.LittleEndian.Uint32(raw[base+8:])))
		z := math.Float32frombits(binary.LittleEndian.Uint32(raw[base+4:]))
		vectors[i] = Types.Vector{X: float64(x) * scaling.X, Y: float64(y) * scaling.Y, Z: float64(z) * scaling.Z}
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

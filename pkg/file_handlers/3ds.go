package FileHandlers

// Copyright (c) 2025 Michael Fogleman
// Portions adapted from FauxGL (https://github.com/fogleman/fauxgl)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the “Software”), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

import (
	"bytes"
	"encoding/binary"
	"io"

	MeshTypes "github.com/Patch2PDF/GDTF-Mesh-Reader/pkg/MeshTypes"
)

func Load3DS(fileData *[]byte, desiredSize *MeshTypes.Vector) (*MeshTypes.Mesh, error) {
	type ChunkHeader struct {
		ChunkID uint16
		Length  uint32
	}

	file := bytes.NewReader(*fileData)

	var vertices []*MeshTypes.Vertex
	var faces []*MeshTypes.Triangle
	var triangles []*MeshTypes.Triangle
	for {
		header := ChunkHeader{}
		if err := binary.Read(file, binary.LittleEndian, &header); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		switch header.ChunkID {
		case 0x4D4D:
		case 0x3D3D:
		case 0x4000:
			_, err := readNullTerminatedString(file)
			if err != nil {
				return nil, err
			}
		case 0x4100:
		case 0x4110:
			v, err := readVertexList(file)
			if err != nil {
				return nil, err
			}
			vertices = v
		case 0x4120:
			f, err := readFaceList(file, vertices)
			if err != nil {
				return nil, err
			}
			faces = f
			triangles = append(triangles, faces...)
		case 0x4150:
			err := readSmoothingGroups(file, faces)
			if err != nil {
				return nil, err
			}
		// case 0x4160:
		// 	matrix, err := readLocalAxis(file)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	for i, v := range vertices {
		// 		vertices[i] = matrix.MulPosition(v)
		// 	}
		default:
			file.Seek(int64(header.Length-6), 1)
		}
	}

	mesh := &MeshTypes.Mesh{Triangles: triangles}

	if desiredSize != nil {
		mesh.ScaleToDimensions(desiredSize)
	}

	return mesh, nil
}

func readSmoothingGroups(file *bytes.Reader, triangles []*MeshTypes.Triangle) error {
	groups := make([]uint32, len(triangles))
	if err := binary.Read(file, binary.LittleEndian, &groups); err != nil {
		return err
	}
	var tables [32]map[MeshTypes.Vector][]MeshTypes.Vector
	for i := 0; i < 32; i++ {
		tables[i] = make(map[MeshTypes.Vector][]MeshTypes.Vector)
	}
	for i, g := range groups {
		t := triangles[i]
		n := t.Normal()
		for j := 0; j < 32; j++ {
			if g&1 == 1 {
				tables[j][t.V0.Position] = append(tables[j][t.V0.Position], n)
				tables[j][t.V1.Position] = append(tables[j][t.V1.Position], n)
				tables[j][t.V2.Position] = append(tables[j][t.V2.Position], n)
			}
			g >>= 1
		}
	}
	for i, g := range groups {
		t := triangles[i]
		var n0, n1, n2 MeshTypes.Vector
		for j := 0; j < 32; j++ {
			if g&1 == 1 {
				for _, v := range tables[j][t.V0.Position] {
					n0 = n0.Add(v)
				}
				for _, v := range tables[j][t.V1.Position] {
					n1 = n1.Add(v)
				}
				for _, v := range tables[j][t.V2.Position] {
					n2 = n2.Add(v)
				}
			}
			g >>= 1
		}
		n0 = n0.Normalize()
		n1 = n1.Normalize()
		n2 = n2.Normalize()
		t.V0.Normal = &n0
		t.V1.Normal = &n1
		t.V2.Normal = &n2
	}
	return nil
}

// func readLocalAxis(file *bytes.Reader) (MeshTypes.Matrix, error) {
// 	var m [4][3]float32
// 	if err := binary.Read(file, binary.LittleEndian, &m); err != nil {
// 		return MeshTypes.Matrix{}, err
// 	}
// 	matrix := MeshTypes.Matrix{
// 		X00: float64(m[0][0]), X01: float64(m[0][1]), X02: float64(m[0][2]), X03: float64(m[3][0]),
// 		X10: float64(m[1][0]), X11: float64(m[1][1]), X12: float64(m[1][2]), X13: float64(m[3][1]),
// 		X20: float64(m[2][0]), X21: float64(m[2][1]), X22: float64(m[2][2]), X23: float64(m[3][2]),
// 		X30: 0, X31: 0, X32: 0, X33: 1,
// 	}
// 	return matrix, nil
// }

func readVertexList(file *bytes.Reader) ([]*MeshTypes.Vertex, error) {
	var count uint16
	if err := binary.Read(file, binary.LittleEndian, &count); err != nil {
		return nil, err
	}
	result := make([]*MeshTypes.Vertex, count)
	for i := range result {
		var v [3]float32
		if err := binary.Read(file, binary.LittleEndian, &v); err != nil {
			return nil, err
		}
		result[i] = &MeshTypes.Vertex{
			Position: MeshTypes.Vector{X: float64(v[0]), Y: float64(v[1]), Z: float64(v[2])},
		}
	}
	return result, nil
}

func readFaceList(file *bytes.Reader, vertices []*MeshTypes.Vertex) ([]*MeshTypes.Triangle, error) {
	var count uint16
	if err := binary.Read(file, binary.LittleEndian, &count); err != nil {
		return nil, err
	}
	result := make([]*MeshTypes.Triangle, count)
	for i := range result {
		var v [4]uint16
		if err := binary.Read(file, binary.LittleEndian, &v); err != nil {
			return nil, err
		}
		result[i] = &MeshTypes.Triangle{
			V0: vertices[v[0]], V1: vertices[v[1]], V2: vertices[v[2]],
		}
	}
	return result, nil
}

func readNullTerminatedString(file *bytes.Reader) (string, error) {
	var bytes []byte
	buf := make([]byte, 1)
	for {
		n, err := file.Read(buf)
		if err != nil {
			return "", err
		} else if n == 1 {
			if buf[0] == 0 {
				break
			}
			bytes = append(bytes, buf[0])
		}
	}
	return string(bytes), nil
}

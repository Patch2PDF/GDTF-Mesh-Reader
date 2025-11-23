package Primitives

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

import Types "github.com/Patch2PDF/GDTF-Mesh-Reader/pkg/MeshTypes"

func NewSphere(detail int) Types.Mesh {
	var triangles []*Types.Triangle
	ico := NewIcosahedron()
	for _, t := range ico.Triangles {
		v1 := t.V0.Position
		v2 := t.V1.Position
		v3 := t.V2.Position
		triangles = append(triangles, newSphereHelper(detail, v1, v2, v3)...)
	}
	return Types.Mesh{Triangles: triangles}
}

func newSphereHelper(detail int, v1, v2, v3 Types.Vector) []*Types.Triangle {
	if detail == 0 {
		t := &Types.Triangle{
			V0: v1.ToVertex(nil),
			V1: v2.ToVertex(nil),
			V2: v3.ToVertex(nil),
		}
		return []*Types.Triangle{t}
	}
	var triangles []*Types.Triangle
	v12 := v1.Add(v2).DivScalar(2).Normalize()
	v13 := v1.Add(v3).DivScalar(2).Normalize()
	v23 := v2.Add(v3).DivScalar(2).Normalize()
	triangles = append(triangles, newSphereHelper(detail-1, v1, v12, v13)...)
	triangles = append(triangles, newSphereHelper(detail-1, v2, v23, v12)...)
	triangles = append(triangles, newSphereHelper(detail-1, v3, v13, v23)...)
	triangles = append(triangles, newSphereHelper(detail-1, v12, v23, v13)...)
	return triangles
}

func NewIcosahedron() *Types.Mesh {
	const a = 0.8506507174597755
	const b = 0.5257312591858783
	vertices := []Types.Vector{
		{X: -a, Y: -b, Z: 0},
		{X: -a, Y: b, Z: 0},
		{X: -b, Y: 0, Z: -a},
		{X: -b, Y: 0, Z: a},
		{X: 0, Y: -a, Z: -b},
		{X: 0, Y: -a, Z: b},
		{X: 0, Y: a, Z: -b},
		{X: 0, Y: a, Z: b},
		{X: b, Y: 0, Z: -a},
		{X: b, Y: 0, Z: a},
		{X: a, Y: -b, Z: 0},
		{X: a, Y: b, Z: 0},
	}
	indices := [][3]int{
		{0, 3, 1},
		{1, 3, 7},
		{2, 0, 1},
		{2, 1, 6},
		{4, 0, 2},
		{4, 5, 0},
		{5, 3, 0},
		{6, 1, 7},
		{6, 7, 11},
		{7, 3, 9},
		{8, 2, 6},
		{8, 4, 2},
		{8, 6, 11},
		{8, 10, 4},
		{8, 11, 10},
		{9, 3, 5},
		{10, 5, 4},
		{10, 9, 5},
		{11, 7, 9},
		{11, 9, 10},
	}
	triangles := make([]*Types.Triangle, len(indices))
	for i, idx := range indices {
		p1 := vertices[idx[0]]
		p2 := vertices[idx[1]]
		p3 := vertices[idx[2]]
		triangles[i] = &Types.Triangle{
			V0: p1.ToVertex(nil),
			V1: p2.ToVertex(nil),
			V2: p3.ToVertex(nil),
		}
	}
	return &Types.Mesh{Triangles: triangles}
}

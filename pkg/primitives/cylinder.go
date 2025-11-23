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

import (
	"math"

	Types "github.com/Patch2PDF/GDTF-Mesh-Reader/pkg/MeshTypes"
)

func radians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func NewCylinder(step int, capped bool) Types.Mesh {
	var triangles []*Types.Triangle
	for a0 := 0; a0 < 360; a0 += step {
		a1 := (a0 + step) % 360
		r0 := radians(float64(a0))
		r1 := radians(float64(a1))
		x0 := math.Cos(r0)
		y0 := math.Sin(r0)
		x1 := math.Cos(r1)
		y1 := math.Sin(r1)
		p00 := Types.Vector{X: x0, Y: y0, Z: -0.5}
		p10 := Types.Vector{X: x1, Y: y1, Z: -0.5}
		p11 := Types.Vector{X: x1, Y: y1, Z: 0.5}
		p01 := Types.Vector{X: x0, Y: y0, Z: 0.5}
		t1 := &Types.Triangle{V0: p00.ToVertex(nil), V1: p10.ToVertex(nil), V2: p11.ToVertex(nil)}
		t2 := &Types.Triangle{V0: p00.ToVertex(nil), V1: p11.ToVertex(nil), V2: p01.ToVertex(nil)}
		triangles = append(triangles, t1)
		triangles = append(triangles, t2)
		if capped {
			p0 := Types.Vector{X: 0, Y: 0, Z: -0.5}
			p1 := Types.Vector{X: 0, Y: 0, Z: 0.5}
			t3 := &Types.Triangle{V0: p0.ToVertex(nil), V1: p10.ToVertex(nil), V2: p00.ToVertex(nil)}
			t4 := &Types.Triangle{V0: p1.ToVertex(nil), V1: p01.ToVertex(nil), V2: p11.ToVertex(nil)}
			triangles = append(triangles, t3)
			triangles = append(triangles, t4)
		}
	}
	return Types.Mesh{Triangles: triangles}
}

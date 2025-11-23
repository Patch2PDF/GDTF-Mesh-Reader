package Primitives

import Types "github.com/Patch2PDF/GDTF-Mesh-Reader/pkg/MeshTypes"

func NewCube() Types.Mesh {
	v := []Types.Vector{
		{X: -1, Y: 1, Z: -1}, {X: -1, Y: 1, Z: 1}, {X: 1, Y: 1, Z: 1}, {X: 1, Y: 1, Z: -1},
		{X: -1, Y: -1, Z: -1}, {X: -1, Y: -1, Z: 1}, {X: 1, Y: -1, Z: 1}, {X: 1, Y: -1, Z: -1},
	}
	mesh := Types.Mesh{
		Triangles: []*Types.Triangle{
			// top
			{
				V0: &Types.Vertex{Position: v[1], Normal: nil},
				V1: &Types.Vertex{Position: v[2], Normal: nil},
				V2: &Types.Vertex{Position: v[3], Normal: nil},
			},
			{
				V0: &Types.Vertex{Position: v[0], Normal: nil},
				V1: &Types.Vertex{Position: v[1], Normal: nil},
				V2: &Types.Vertex{Position: v[3], Normal: nil},
			},
			// left
			{
				V0: &Types.Vertex{Position: v[4], Normal: nil},
				V1: &Types.Vertex{Position: v[1], Normal: nil},
				V2: &Types.Vertex{Position: v[0], Normal: nil},
			},
			{
				V0: &Types.Vertex{Position: v[4], Normal: nil},
				V1: &Types.Vertex{Position: v[5], Normal: nil},
				V2: &Types.Vertex{Position: v[1], Normal: nil},
			},
			// front
			{
				V0: &Types.Vertex{Position: v[5], Normal: nil},
				V1: &Types.Vertex{Position: v[2], Normal: nil},
				V2: &Types.Vertex{Position: v[1], Normal: nil},
			},
			{
				V0: &Types.Vertex{Position: v[5], Normal: nil},
				V1: &Types.Vertex{Position: v[6], Normal: nil},
				V2: &Types.Vertex{Position: v[2], Normal: nil},
			},
			// right
			{
				V0: &Types.Vertex{Position: v[7], Normal: nil},
				V1: &Types.Vertex{Position: v[3], Normal: nil},
				V2: &Types.Vertex{Position: v[2], Normal: nil},
			},
			{
				V0: &Types.Vertex{Position: v[7], Normal: nil},
				V1: &Types.Vertex{Position: v[2], Normal: nil},
				V2: &Types.Vertex{Position: v[6], Normal: nil},
			},
			// back
			{
				V0: &Types.Vertex{Position: v[0], Normal: nil},
				V1: &Types.Vertex{Position: v[3], Normal: nil},
				V2: &Types.Vertex{Position: v[4], Normal: nil},
			},
			{
				V0: &Types.Vertex{Position: v[7], Normal: nil},
				V1: &Types.Vertex{Position: v[4], Normal: nil},
				V2: &Types.Vertex{Position: v[3], Normal: nil},
			},
			// bottom
			{
				V0: &Types.Vertex{Position: v[6], Normal: nil},
				V1: &Types.Vertex{Position: v[5], Normal: nil},
				V2: &Types.Vertex{Position: v[4], Normal: nil},
			},
			{
				V0: &Types.Vertex{Position: v[7], Normal: nil},
				V1: &Types.Vertex{Position: v[6], Normal: nil},
				V2: &Types.Vertex{Position: v[4], Normal: nil},
			},
		},
	}
	return mesh
}

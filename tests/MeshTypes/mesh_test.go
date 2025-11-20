package MeshTypes_Test

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/Patch2PDF/GDTF-Mesh-Reader/pkg/MeshTypes"
)

func RandomTriangle() *MeshTypes.Triangle {
	return &MeshTypes.Triangle{
		V0: RandomVertex(),
		V1: RandomVertex(),
		V2: RandomVertex(),
	}
}

func RandomVertex() *MeshTypes.Vertex {
	normal := RandomVector()
	return &MeshTypes.Vertex{
		Position: RandomVector(),
		Normal:   &normal,
	}
}

func RandomVector() MeshTypes.Vector {
	return MeshTypes.Vector{
		X: rand.Float64(),
		Y: rand.Float64(),
		Z: rand.Float64(),
	}
}

func TestAddTriangle(t *testing.T) {
	a := MeshTypes.Mesh{
		Triangles: []*MeshTypes.Triangle{RandomTriangle()},
	}
	b := RandomTriangle()
	copy := a.Copy()
	copy.AddTriangle(b)
	if !reflect.DeepEqual(copy, MeshTypes.Mesh{Triangles: []*MeshTypes.Triangle{a.Triangles[0], b}}) {
		t.Errorf("Mesh AddTriangle() Output does not match expected")
	}
}

func TestMeshCopy(t *testing.T) {
	a := MeshTypes.Mesh{
		Triangles: []*MeshTypes.Triangle{RandomTriangle()},
	}
	copy := a.Copy()
	if !(reflect.DeepEqual(a, copy) && &a != &copy) {
		t.Errorf("Mesh Copy() Output does not match expected")
	}
}

func TestAddMesh(t *testing.T) {
	a := MeshTypes.Mesh{
		Triangles: []*MeshTypes.Triangle{RandomTriangle()},
	}
	b := MeshTypes.Mesh{
		Triangles: []*MeshTypes.Triangle{RandomTriangle()},
	}
	copy := a.Copy()
	result := copy.Add(&b)
	if !reflect.DeepEqual(*result, MeshTypes.Mesh{Triangles: []*MeshTypes.Triangle{a.Triangles[0], b.Triangles[0]}}) {
		t.Errorf("Mesh Add() Output does not match expected")
	}
}

func TestRotateAndTranslate(t *testing.T) {
	a := MeshTypes.Mesh{
		Triangles: []*MeshTypes.Triangle{RandomTriangle()},
	}
	translationMatrix := MeshTypes.Matrix{
		X00: 1, X01: -1, X02: 5, X03: -20,
		X10: 1, X11: 2, X12: -3, X13: 10,
		X20: 1, X21: 5, X22: 3, X23: -5,
		X30: 1, X31: 6, X32: -6, X33: 4,
	}
	result := a.Copy()
	result.RotateAndTranslate(translationMatrix)
	want := a.Copy()
	for _, triangle := range want.Triangles {
		triangle.V0.Position = translationMatrix.MulPosition(triangle.V0.Position) // safe to use as func is tested in another place
		triangle.V1.Position = translationMatrix.MulPosition(triangle.V1.Position)
		triangle.V2.Position = translationMatrix.MulPosition(triangle.V2.Position)
	}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Mesh RotateAndTranslate() Output does not match expected")
	}
}

func TestScaleToDimensions(t *testing.T) {
	a := MeshTypes.Mesh{
		Triangles: []*MeshTypes.Triangle{
			{
				V0: &MeshTypes.Vertex{Position: MeshTypes.Vector{X: -1, Y: 1, Z: -1}},
				V1: &MeshTypes.Vertex{Position: MeshTypes.Vector{X: 1, Y: -3, Z: 1}},
				V2: &MeshTypes.Vertex{Position: MeshTypes.Vector{X: 1, Y: 1, Z: 0}},
			},
		},
	}
	desiredSize := MeshTypes.Vector{
		X: 4, Y: 4, Z: 8,
	}
	result := a.Copy()
	result.ScaleToDimensions(&desiredSize)
	want := MeshTypes.Mesh{
		Triangles: []*MeshTypes.Triangle{
			{
				V0: &MeshTypes.Vertex{Position: MeshTypes.Vector{X: -2, Y: 1, Z: -4}},
				V1: &MeshTypes.Vertex{Position: MeshTypes.Vector{X: 2, Y: -3, Z: 4}},
				V2: &MeshTypes.Vertex{Position: MeshTypes.Vector{X: 2, Y: 1, Z: 0}},
			},
		},
	}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Mesh RotateAndTranslate() Output does not match expected")
	}
}

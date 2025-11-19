package MeshTypes_Test

import (
	"reflect"
	"testing"

	"github.com/Patch2PDF/GDTF-Mesh-Reader/pkg/MeshTypes"
)

func TestTriangleNormal(t *testing.T) {
	a := MeshTypes.Triangle{
		V0: &MeshTypes.Vertex{Position: MeshTypes.Vector{X: 1, Y: 4, Z: 3}, Normal: &MeshTypes.Vector{X: 4, Y: 5, Z: 6}},
		V1: &MeshTypes.Vertex{Position: MeshTypes.Vector{X: 4, Y: 8, Z: 10}, Normal: &MeshTypes.Vector{X: 10, Y: 11, Z: 12}},
		V2: &MeshTypes.Vertex{Position: MeshTypes.Vector{X: 13, Y: 14, Z: 15}, Normal: &MeshTypes.Vector{X: 16, Y: 17, Z: 18}},
	}
	want := MeshTypes.Vector{X: -0.39436910666014113, Y: 0.8604416872584897, Z: -0.3226656327219336}
	result := a.Normal()
	if !reflect.DeepEqual(result, want) {
		t.Error("Triangle Normal() returned value does not match expected value")
	}
}

func TestTriangleCopy(t *testing.T) {
	a := MeshTypes.Triangle{
		V0: &MeshTypes.Vertex{Position: MeshTypes.Vector{X: 1, Y: 2, Z: 3}, Normal: &MeshTypes.Vector{X: 4, Y: 5, Z: 6}},
		V1: &MeshTypes.Vertex{Position: MeshTypes.Vector{X: 7, Y: 8, Z: 9}, Normal: &MeshTypes.Vector{X: 10, Y: 11, Z: 12}},
		V2: &MeshTypes.Vertex{Position: MeshTypes.Vector{X: 13, Y: 14, Z: 15}, Normal: &MeshTypes.Vector{X: 16, Y: 17, Z: 18}},
	}
	copy := a.Copy()
	if !(reflect.DeepEqual(copy, a) && a.V0 != copy.V0 && a.V1 != copy.V1 && a.V2 != copy.V2) {
		t.Error("Triangle Copy() returned value does not match expected value")
	}
}

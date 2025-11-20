package MeshTypes_Test

import (
	"reflect"
	"testing"

	"github.com/Patch2PDF/GDTF-Mesh-Reader/pkg/MeshTypes"
)

func TestVertexCopy(t *testing.T) {
	a := MeshTypes.Vertex{
		Position: MeshTypes.Vector{
			X: 1,
			Y: 2,
			Z: 3,
		},
		Normal: &MeshTypes.Vector{
			X: 4,
			Y: 5,
			Z: 6,
		},
	}
	copy := a.Copy()
	if &a.Position == &copy.Position {
		t.Errorf("Vertex Copy func returns same Position Vector")
	}
	if a.Normal == copy.Normal {
		t.Errorf("Vertex Copy func returns same Normal Vector")
	}
	if !reflect.DeepEqual(a.Position, copy.Position) {
		t.Errorf("Vertex Copy func does not return same Position Vector Values")
	}
	if !reflect.DeepEqual(a.Normal, copy.Normal) {
		t.Errorf("Vertex Copy func does not return same Normal Vector Values")
	}
}

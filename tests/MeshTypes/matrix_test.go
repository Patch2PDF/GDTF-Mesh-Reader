package MeshTypes_Test

import (
	"reflect"
	"testing"

	"github.com/Patch2PDF/GDTF-Mesh-Reader/pkg/MeshTypes"
)

func TestIdentityMatrix(t *testing.T) {
	want := MeshTypes.Matrix{
		X00: 1, X01: 0, X02: 0, X03: 0,
		X10: 0, X11: 1, X12: 0, X13: 0,
		X20: 0, X21: 0, X22: 1, X23: 0,
		X30: 0, X31: 0, X32: 0, X33: 1,
	}
	if !reflect.DeepEqual(MeshTypes.IdentityMatrix(), want) {
		t.Errorf(`IdentityMatrix() Output does not match`)
	}
}

func TestMatrixMul(t *testing.T) {
	a := MeshTypes.Matrix{
		X00: 1, X01: 2, X02: 3, X03: 4,
		X10: 5, X11: 6, X12: 7, X13: 8,
		X20: 9, X21: 10, X22: 11, X23: 12,
		X30: 13, X31: 14, X32: 15, X33: 16,
	}
	b := MeshTypes.Matrix{
		X00: 17, X01: 18, X02: 19, X03: 20,
		X10: 21, X11: 22, X12: 23, X13: 24,
		X20: 25, X21: 26, X22: 27, X23: 28,
		X30: 29, X31: 30, X32: 31, X33: 32,
	}
	want := MeshTypes.Matrix{
		X00: 250, X01: 260, X02: 270, X03: 280,
		X10: 618, X11: 644, X12: 670, X13: 696,
		X20: 986, X21: 1028, X22: 1070, X23: 1112,
		X30: 1354, X31: 1412, X32: 1470, X33: 1528,
	}
	if !reflect.DeepEqual(a.Mul(b), want) {
		t.Errorf(`Matrix Multiplication Output does not match`)
	}
}

func TestMulPosition(t *testing.T) {
	a := MeshTypes.Matrix{
		X00: 1, X01: 2, X02: 3, X03: 4,
		X10: 5, X11: 6, X12: 7, X13: 8,
		X20: 9, X21: 10, X22: 11, X23: 12,
		X30: 13, X31: 14, X32: 15, X33: 16,
	}
	b := MeshTypes.Vector{
		X: 17,
		Y: 18,
		Z: 19,
	}
	want := MeshTypes.Vector{
		X: 114,
		Y: 334,
		Z: 554,
	}

	if !reflect.DeepEqual(a.MulPosition(b), want) {
		t.Errorf(`Matrix Vector Multiplication Output does not match`)
	}
}

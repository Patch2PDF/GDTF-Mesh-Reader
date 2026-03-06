package MeshTypes_Test

import (
	"math"
	"math/rand"
	"reflect"
	"testing"

	"github.com/Patch2PDF/GDTF-Mesh-Reader/v2/pkg/MeshTypes"
)

func MatrixEquals(a MeshTypes.Matrix, b MeshTypes.Matrix) bool {
	// Helper to check individual floats
	isClose := func(a, b float64) bool {
		return math.Abs(a-b) < 0.000000000000001
	}

	return isClose(a.X00, b.X00) && isClose(a.X01, b.X01) &&
		isClose(a.X02, b.X02) && isClose(a.X03, b.X03) &&
		isClose(a.X10, b.X10) && isClose(a.X11, b.X11) &&
		isClose(a.X12, b.X12) && isClose(a.X13, b.X13) &&
		isClose(a.X20, b.X20) && isClose(a.X21, b.X21) &&
		isClose(a.X22, b.X22) && isClose(a.X23, b.X23) &&
		isClose(a.X30, b.X30) && isClose(a.X31, b.X31) &&
		isClose(a.X32, b.X32) && isClose(a.X33, b.X33)
}

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

func TestMulScalar(t *testing.T) {
	a := MeshTypes.Matrix{
		X00: 1, X01: 2, X02: 3, X03: 4,
		X10: 5, X11: 6, X12: 7, X13: 8,
		X20: 9, X21: 10, X22: 11, X23: 12,
		X30: 13, X31: 14, X32: 15, X33: 16,
	}
	want := MeshTypes.Matrix{
		X00: 3, X01: 6, X02: 9, X03: 12,
		X10: 15, X11: 18, X12: 21, X13: 24,
		X20: 27, X21: 30, X22: 33, X23: 36,
		X30: 39, X31: 42, X32: 45, X33: 48,
	}
	result := a.MulScalar(3)
	if !reflect.DeepEqual(result, want) {
		t.Errorf(`Matrix Scalar Multiplication Output does not match`)
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

func TestRotation(t *testing.T) {
	a := MeshTypes.Matrix{
		X00: rand.Float64(), X01: rand.Float64(), X02: rand.Float64(), X03: rand.Float64(),
		X10: rand.Float64(), X11: rand.Float64(), X12: rand.Float64(), X13: rand.Float64(),
		X20: rand.Float64(), X21: rand.Float64(), X22: rand.Float64(), X23: rand.Float64(),
		X30: 0, X31: 0, X32: 0, X33: 1,
	}

	alpha := rand.Float64()
	beta := rand.Float64()
	gamma := rand.Float64()

	rotation := MeshTypes.GenerateRotationMatrix(alpha, beta, gamma)

	if !reflect.DeepEqual(a.Mul(rotation), a.Rotate(alpha, beta, gamma)) {
		t.Errorf(`Matrix Vector Rotation Output does not match`)
	}
}

func TestMatrixRotationReversal(t *testing.T) {
	a := MeshTypes.Matrix{
		X00: rand.Float64(), X01: rand.Float64(), X02: rand.Float64(), X03: rand.Float64(),
		X10: rand.Float64(), X11: rand.Float64(), X12: rand.Float64(), X13: rand.Float64(),
		X20: rand.Float64(), X21: rand.Float64(), X22: rand.Float64(), X23: rand.Float64(),
		X30: 0, X31: 0, X32: 0, X33: 1,
	}

	rotation := MeshTypes.GenerateRotationMatrix(rand.Float64(), rand.Float64(), rand.Float64())

	rotated := a.Mul(rotation)

	back_rotated := rotated.ReverseTransformation(rotation)

	if !MatrixEquals(a, back_rotated) {
		t.Errorf(`Matrix Vector Rotation Reversal Output does not match`)
	}
}

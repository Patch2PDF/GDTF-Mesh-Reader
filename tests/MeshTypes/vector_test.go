package MeshTypes_Test

import (
	"reflect"
	"testing"

	"github.com/Patch2PDF/GDTF-Mesh-Reader/pkg/MeshTypes"
)

func TestVectorDivScalar(t *testing.T) {
	a := MeshTypes.Vector{
		X: 2, Y: 4, Z: 6,
	}
	b := 2.0
	want := MeshTypes.Vector{
		X: 1, Y: 2, Z: 3,
	}
	if !reflect.DeepEqual(a.DivScalar(b), want) {
		t.Error("Vector DivScalar() returned value does not match expected value")
	}
}

func TestVectorAdd(t *testing.T) {
	a := MeshTypes.Vector{
		X: 1, Y: 2, Z: 3,
	}
	b := MeshTypes.Vector{
		X: 4, Y: 5, Z: 6,
	}
	want := MeshTypes.Vector{
		X: 5, Y: 7, Z: 9,
	}
	if !reflect.DeepEqual(a.Add(b), want) {
		t.Error("Vector Add() returned value does not match expected value")
	}
}

func TestVectorSub(t *testing.T) {
	a := MeshTypes.Vector{
		X: 1, Y: 2, Z: 3,
	}
	b := MeshTypes.Vector{
		X: 4, Y: 6, Z: 9,
	}
	want := MeshTypes.Vector{
		X: 3, Y: 4, Z: 6,
	}
	if !reflect.DeepEqual(b.Sub(a), want) {
		t.Error("Vector Sub() returned value does not match expected value")
	}
}

func TestVectorCross(t *testing.T) {
	a := MeshTypes.Vector{
		X: 1, Y: 2, Z: 4,
	}
	b := MeshTypes.Vector{
		X: 5, Y: 6, Z: 7,
	}
	want := MeshTypes.Vector{
		X: -10, Y: 13, Z: -4,
	}
	result := a.Cross(b)
	if !reflect.DeepEqual(result, want) {
		t.Error("Vector Cross() returned value does not match expected value")
	}
}

func TestNormalize(t *testing.T) {
	a := MeshTypes.Vector{
		X: 3, Y: 4, Z: 12,
	}
	want := MeshTypes.Vector{
		X: 3.0 / 13, Y: 4.0 / 13, Z: 12.0 / 13,
	}
	if !reflect.DeepEqual(a.Normalize(), want) {
		t.Error("Vector Normalize() returned value does not match expected value")
	}
}

func TestToVertex(t *testing.T) {
	a := MeshTypes.Vector{
		X: 1, Y: 2, Z: 3,
	}
	b := MeshTypes.Vector{
		X: 4, Y: 5, Z: 6,
	}
	wantFirst := MeshTypes.Vertex{
		Position: a,
		Normal:   nil,
	}
	if !reflect.DeepEqual(*a.ToVertex(nil), wantFirst) {
		t.Error("Vector To Vertex returned value with nil Normal does not match expected value")
	}
	wantSecond := MeshTypes.Vertex{
		Position: a,
		Normal:   &b,
	}
	if !reflect.DeepEqual(*a.ToVertex(&b), wantSecond) || &b != wantSecond.Normal {
		t.Error("Vector To Vertex returned value with Normal does not match expected value")
	}
}

func TestMin(t *testing.T) {
	a := MeshTypes.Vector{
		X: 1, Y: 4, Z: 3,
	}
	b := MeshTypes.Vector{
		X: 4, Y: 2, Z: 6,
	}
	want := MeshTypes.Vector{
		X: 1, Y: 2, Z: 3,
	}
	if !reflect.DeepEqual(a.Min(&b), want) {
		t.Error("Vector Min() returned value does not match expected value")
	}
}

func TestMax(t *testing.T) {
	a := MeshTypes.Vector{
		X: 1, Y: 4, Z: 3,
	}
	b := MeshTypes.Vector{
		X: 3, Y: 2, Z: 6,
	}
	want := MeshTypes.Vector{
		X: 3, Y: 4, Z: 6,
	}
	if !reflect.DeepEqual(a.Max(&b), want) {
		t.Error("Vector Max() returned value does not match expected value")
	}
}

func TestMult(t *testing.T) {
	a := MeshTypes.Vector{
		X: 1, Y: 4, Z: 3,
	}
	b := MeshTypes.Vector{
		X: 3, Y: 2, Z: 6,
	}
	want := MeshTypes.Vector{
		X: 3, Y: 8, Z: 18,
	}
	if !reflect.DeepEqual(a.Mult(b), want) {
		t.Error("Vector Mul() returned value does not match expected value")
	}
}

func TestDiv(t *testing.T) {
	a := MeshTypes.Vector{
		X: 6, Y: 8, Z: 24,
	}
	b := MeshTypes.Vector{
		X: 3, Y: 2, Z: 4,
	}
	want := MeshTypes.Vector{
		X: 2, Y: 4, Z: 6,
	}
	if !reflect.DeepEqual(a.Div(b), want) {
		t.Error("Vector Div() returned value does not match expected value")
	}
}

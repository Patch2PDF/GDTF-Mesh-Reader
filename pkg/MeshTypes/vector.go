package MeshTypes

import "math"

// Vector data type
type Vector struct {
	X float64
	Y float64
	Z float64
}

// Add Vector, returns resulting vector
func (a Vector) Add(b Vector) Vector {
	return Vector{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

// Subtract Vector, returns resulting vector
func (a Vector) Sub(b Vector) Vector {
	return Vector{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

// Calculate Crossproduct, returns resulting vector
func (a Vector) Cross(b Vector) Vector {
	x := a.Y*b.Z - a.Z*b.Y
	y := a.Z*b.X - a.X*b.Z
	z := a.X*b.Y - a.Y*b.X
	return Vector{x, y, z}
}

// Normalize Vector, returns resulting vector
func (a Vector) Normalize() Vector {
	r := 1 / math.Sqrt(a.X*a.X+a.Y*a.Y+a.Z*a.Z)
	return Vector{a.X * r, a.Y * r, a.Z * r}
}

// Convert position vector to vertex
func (obj Vector) ToVertex(normal *Vector) *Vertex {
	return &Vertex{
		Position: obj,
		Normal:   normal,
	}
}

// Get minimum vector
// Compares each element seperately, result can be a mixture of both
func (a Vector) Min(b *Vector) Vector {
	return Vector{
		math.Min(a.X, b.X),
		math.Min(a.Y, b.Y),
		math.Min(a.Z, b.Z),
	}
}

// Get maximum vector
// Compares each element seperately, result can be a mixture of both
func (a Vector) Max(b *Vector) Vector {
	return Vector{
		math.Max(a.X, b.X),
		math.Max(a.Y, b.Y),
		math.Max(a.Z, b.Z),
	}
}

// Multiply vectors, returns resulting vector
func (a Vector) Mult(b Vector) Vector {
	return Vector{
		a.X * b.X,
		a.Y * b.Y,
		a.Z * b.Z,
	}
}

// Divide vectors, returns resulting vector
func (a Vector) Div(b Vector) Vector {
	return Vector{
		a.X / b.X,
		a.Y / b.Y,
		a.Z / b.Z,
	}
}

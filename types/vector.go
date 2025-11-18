package Types

import "math"

type Vector struct {
	X float64
	Y float64
	Z float64
}

func (a Vector) Add(b Vector) Vector {
	return Vector{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

func (a Vector) Sub(b Vector) Vector {
	return Vector{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

func (a Vector) Cross(b Vector) Vector {
	x := a.Y*b.Z - a.Z*b.Y
	y := a.Z*b.X - a.X*b.Z
	z := a.X*b.Y - a.Y*b.X
	return Vector{x, y, z}
}

func (a Vector) Normalize() Vector {
	r := 1 / math.Sqrt(a.X*a.X+a.Y*a.Y+a.Z*a.Z)
	return Vector{a.X * r, a.Y * r, a.Z * r}
}

func (obj Vector) ToVertex(normal *Vector) *Vertex {
	return &Vertex{
		Position: obj,
		Normal:   normal,
	}
}

func (a Vector) Min(b *Vector) Vector {
	return Vector{
		math.Min(a.X, b.X),
		math.Min(a.Y, b.Y),
		math.Min(a.Z, b.Z),
	}
}

func (a Vector) Max(b *Vector) Vector {
	return Vector{
		math.Max(a.X, b.X),
		math.Max(a.Y, b.Y),
		math.Max(a.Z, b.Z),
	}
}

func (a Vector) Mult(b Vector) Vector {
	return Vector{
		a.X * b.X,
		a.Y * b.Y,
		a.Z * b.Z,
	}
}

func (a Vector) Div(b Vector) Vector {
	return Vector{
		a.X / b.X,
		a.Y / b.Y,
		a.Z / b.Z,
	}
}

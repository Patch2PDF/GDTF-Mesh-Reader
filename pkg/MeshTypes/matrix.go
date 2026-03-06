package MeshTypes

import "math"

type Matrix struct {
	X00, X01, X02, X03 float64
	X10, X11, X12, X13 float64
	X20, X21, X22, X23 float64
	X30, X31, X32, X33 float64
}

func IdentityMatrix() Matrix {
	return Matrix{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

func (a Matrix) Mul(b Matrix) Matrix {
	m := Matrix{}
	m.X00 = a.X00*b.X00 + a.X01*b.X10 + a.X02*b.X20 + a.X03*b.X30
	m.X10 = a.X10*b.X00 + a.X11*b.X10 + a.X12*b.X20 + a.X13*b.X30
	m.X20 = a.X20*b.X00 + a.X21*b.X10 + a.X22*b.X20 + a.X23*b.X30
	m.X30 = a.X30*b.X00 + a.X31*b.X10 + a.X32*b.X20 + a.X33*b.X30
	m.X01 = a.X00*b.X01 + a.X01*b.X11 + a.X02*b.X21 + a.X03*b.X31
	m.X11 = a.X10*b.X01 + a.X11*b.X11 + a.X12*b.X21 + a.X13*b.X31
	m.X21 = a.X20*b.X01 + a.X21*b.X11 + a.X22*b.X21 + a.X23*b.X31
	m.X31 = a.X30*b.X01 + a.X31*b.X11 + a.X32*b.X21 + a.X33*b.X31
	m.X02 = a.X00*b.X02 + a.X01*b.X12 + a.X02*b.X22 + a.X03*b.X32
	m.X12 = a.X10*b.X02 + a.X11*b.X12 + a.X12*b.X22 + a.X13*b.X32
	m.X22 = a.X20*b.X02 + a.X21*b.X12 + a.X22*b.X22 + a.X23*b.X32
	m.X32 = a.X30*b.X02 + a.X31*b.X12 + a.X32*b.X22 + a.X33*b.X32
	m.X03 = a.X00*b.X03 + a.X01*b.X13 + a.X02*b.X23 + a.X03*b.X33
	m.X13 = a.X10*b.X03 + a.X11*b.X13 + a.X12*b.X23 + a.X13*b.X33
	m.X23 = a.X20*b.X03 + a.X21*b.X13 + a.X22*b.X23 + a.X23*b.X33
	m.X33 = a.X30*b.X03 + a.X31*b.X13 + a.X32*b.X23 + a.X33*b.X33
	return m
}

func (a Matrix) MulScalar(b float64) Matrix {
	return Matrix{
		a.X00 * b, a.X01 * b, a.X02 * b, a.X03 * b,
		a.X10 * b, a.X11 * b, a.X12 * b, a.X13 * b,
		a.X20 * b, a.X21 * b, a.X22 * b, a.X23 * b,
		a.X30 * b, a.X31 * b, a.X32 * b, a.X33 * b,
	}
}

func (a Matrix) MulPosition(b Vector) Vector {
	x := a.X00*b.X + a.X01*b.Y + a.X02*b.Z + a.X03
	y := a.X10*b.X + a.X11*b.Y + a.X12*b.Z + a.X13
	z := a.X20*b.X + a.X21*b.Y + a.X22*b.Z + a.X23
	return Vector{x, y, z}
}

func (a Matrix) MulDirection(b Vector) Vector {
	x := a.X00*b.X + a.X01*b.Y + a.X02*b.Z
	y := a.X10*b.X + a.X11*b.Y + a.X12*b.Z
	z := a.X20*b.X + a.X21*b.Y + a.X22*b.Z
	return Vector{x, y, z}.Normalize()
}

func GenerateRotationMatrix(alpha float64, beta float64, gamma float64) Matrix {
	alphaSin := math.Sin(alpha / 180 * math.Pi)
	alphaCos := math.Cos(alpha / 180 * math.Pi)
	betaSin := math.Sin(beta / 180 * math.Pi)
	betaCos := math.Cos(beta / 180 * math.Pi)
	gammaSin := math.Sin(gamma / 180 * math.Pi)
	gammaCos := math.Cos(gamma / 180 * math.Pi)

	return Matrix{
		X00: betaCos * gammaCos, X01: -betaCos * gammaSin, X02: betaSin, X03: 0,
		X10: alphaCos*gammaSin + alphaSin*betaSin*gammaCos, X11: alphaCos*gammaCos - alphaSin*betaSin*gammaSin, X12: -alphaSin * betaCos, X13: 0,
		X20: alphaSin*gammaSin - alphaCos*betaSin*gammaCos, X21: alphaSin*gammaCos + alphaCos*betaSin*gammaSin, X22: alphaCos * betaCos, X23: 0,
		X30: 0, X31: 0, X32: 0, X33: 1,
	}
}

func (a Matrix) Rotate(alpha float64, beta float64, gamma float64) Matrix {
	return a.Mul(GenerateRotationMatrix(alpha, beta, gamma))
}

func (a Matrix) ReverseTransformation(previousRotationMatrix Matrix) Matrix {
	inv := Matrix{
		X00: previousRotationMatrix.X00, X01: previousRotationMatrix.X10, X02: previousRotationMatrix.X20, X03: 0,
		X10: previousRotationMatrix.X01, X11: previousRotationMatrix.X11, X12: previousRotationMatrix.X21, X13: 0,
		X20: previousRotationMatrix.X02, X21: previousRotationMatrix.X12, X22: previousRotationMatrix.X22, X23: 0,
		X30: 0, X31: 0, X32: 0, X33: 1,
	}

	inv.X03 = -(inv.X00*previousRotationMatrix.X03 + inv.X01*previousRotationMatrix.X13 + inv.X02*previousRotationMatrix.X23)
	inv.X13 = -(inv.X10*previousRotationMatrix.X03 + inv.X11*previousRotationMatrix.X13 + inv.X12*previousRotationMatrix.X23)
	inv.X23 = -(inv.X20*previousRotationMatrix.X03 + inv.X21*previousRotationMatrix.X13 + inv.X22*previousRotationMatrix.X23)

	return a.Mul(inv)
}

func (a Matrix) Transpose() Matrix {
	return Matrix{
		a.X00, a.X10, a.X20, a.X30,
		a.X01, a.X11, a.X21, a.X31,
		a.X02, a.X12, a.X22, a.X32,
		a.X03, a.X13, a.X23, a.X33,
	}
}

func (a Matrix) Determinant() float64 {
	return (a.X00*a.X11*a.X22*a.X33 - a.X00*a.X11*a.X23*a.X32 +
		a.X00*a.X12*a.X23*a.X31 - a.X00*a.X12*a.X21*a.X33 +
		a.X00*a.X13*a.X21*a.X32 - a.X00*a.X13*a.X22*a.X31 -
		a.X01*a.X12*a.X23*a.X30 + a.X01*a.X12*a.X20*a.X33 -
		a.X01*a.X13*a.X20*a.X32 + a.X01*a.X13*a.X22*a.X30 -
		a.X01*a.X10*a.X22*a.X33 + a.X01*a.X10*a.X23*a.X32 +
		a.X02*a.X13*a.X20*a.X31 - a.X02*a.X13*a.X21*a.X30 +
		a.X02*a.X10*a.X21*a.X33 - a.X02*a.X10*a.X23*a.X31 +
		a.X02*a.X11*a.X23*a.X30 - a.X02*a.X11*a.X20*a.X33 -
		a.X03*a.X10*a.X21*a.X32 + a.X03*a.X10*a.X22*a.X31 -
		a.X03*a.X11*a.X22*a.X30 + a.X03*a.X11*a.X20*a.X32 -
		a.X03*a.X12*a.X20*a.X31 + a.X03*a.X12*a.X21*a.X30)
}

func (a Matrix) Inverse() Matrix {
	m := Matrix{}
	d := a.Determinant()
	m.X00 = (a.X12*a.X23*a.X31 - a.X13*a.X22*a.X31 + a.X13*a.X21*a.X32 - a.X11*a.X23*a.X32 - a.X12*a.X21*a.X33 + a.X11*a.X22*a.X33) / d
	m.X01 = (a.X03*a.X22*a.X31 - a.X02*a.X23*a.X31 - a.X03*a.X21*a.X32 + a.X01*a.X23*a.X32 + a.X02*a.X21*a.X33 - a.X01*a.X22*a.X33) / d
	m.X02 = (a.X02*a.X13*a.X31 - a.X03*a.X12*a.X31 + a.X03*a.X11*a.X32 - a.X01*a.X13*a.X32 - a.X02*a.X11*a.X33 + a.X01*a.X12*a.X33) / d
	m.X03 = (a.X03*a.X12*a.X21 - a.X02*a.X13*a.X21 - a.X03*a.X11*a.X22 + a.X01*a.X13*a.X22 + a.X02*a.X11*a.X23 - a.X01*a.X12*a.X23) / d
	m.X10 = (a.X13*a.X22*a.X30 - a.X12*a.X23*a.X30 - a.X13*a.X20*a.X32 + a.X10*a.X23*a.X32 + a.X12*a.X20*a.X33 - a.X10*a.X22*a.X33) / d
	m.X11 = (a.X02*a.X23*a.X30 - a.X03*a.X22*a.X30 + a.X03*a.X20*a.X32 - a.X00*a.X23*a.X32 - a.X02*a.X20*a.X33 + a.X00*a.X22*a.X33) / d
	m.X12 = (a.X03*a.X12*a.X30 - a.X02*a.X13*a.X30 - a.X03*a.X10*a.X32 + a.X00*a.X13*a.X32 + a.X02*a.X10*a.X33 - a.X00*a.X12*a.X33) / d
	m.X13 = (a.X02*a.X13*a.X20 - a.X03*a.X12*a.X20 + a.X03*a.X10*a.X22 - a.X00*a.X13*a.X22 - a.X02*a.X10*a.X23 + a.X00*a.X12*a.X23) / d
	m.X20 = (a.X11*a.X23*a.X30 - a.X13*a.X21*a.X30 + a.X13*a.X20*a.X31 - a.X10*a.X23*a.X31 - a.X11*a.X20*a.X33 + a.X10*a.X21*a.X33) / d
	m.X21 = (a.X03*a.X21*a.X30 - a.X01*a.X23*a.X30 - a.X03*a.X20*a.X31 + a.X00*a.X23*a.X31 + a.X01*a.X20*a.X33 - a.X00*a.X21*a.X33) / d
	m.X22 = (a.X01*a.X13*a.X30 - a.X03*a.X11*a.X30 + a.X03*a.X10*a.X31 - a.X00*a.X13*a.X31 - a.X01*a.X10*a.X33 + a.X00*a.X11*a.X33) / d
	m.X23 = (a.X03*a.X11*a.X20 - a.X01*a.X13*a.X20 - a.X03*a.X10*a.X21 + a.X00*a.X13*a.X21 + a.X01*a.X10*a.X23 - a.X00*a.X11*a.X23) / d
	m.X30 = (a.X12*a.X21*a.X30 - a.X11*a.X22*a.X30 - a.X12*a.X20*a.X31 + a.X10*a.X22*a.X31 + a.X11*a.X20*a.X32 - a.X10*a.X21*a.X32) / d
	m.X31 = (a.X01*a.X22*a.X30 - a.X02*a.X21*a.X30 + a.X02*a.X20*a.X31 - a.X00*a.X22*a.X31 - a.X01*a.X20*a.X32 + a.X00*a.X21*a.X32) / d
	m.X32 = (a.X02*a.X11*a.X30 - a.X01*a.X12*a.X30 - a.X02*a.X10*a.X31 + a.X00*a.X12*a.X31 + a.X01*a.X10*a.X32 - a.X00*a.X11*a.X32) / d
	m.X33 = (a.X01*a.X12*a.X20 - a.X02*a.X11*a.X20 + a.X02*a.X10*a.X21 - a.X00*a.X12*a.X21 - a.X01*a.X10*a.X22 + a.X00*a.X11*a.X22) / d
	return m
}

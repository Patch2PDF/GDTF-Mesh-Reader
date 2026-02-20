package MeshTypes

import (
	"fmt"
)

type Mesh struct {
	Triangles []Triangle
}

func (obj *Mesh) AddTriangle(triangle Triangle) {
	obj.Triangles = append(obj.Triangles, triangle)
}

func (obj *Mesh) Copy() Mesh {
	triangles := make([]Triangle, len(obj.Triangles))
	copy(triangles, obj.Triangles)
	return Mesh{Triangles: triangles}
}

func (obj *Mesh) Add(mesh ...*Mesh) *Mesh {
	for _, element := range mesh {
		if element != nil {
			obj.Triangles = append(obj.Triangles, element.Triangles...)
		}
	}
	return obj
}

func (obj *Mesh) RotateAndTranslate(translationMatrix Matrix) {
	for triangle_id := range obj.Triangles {
		obj.Triangles[triangle_id].V0.Position = translationMatrix.MulPosition(obj.Triangles[triangle_id].V0.Position)
		obj.Triangles[triangle_id].V1.Position = translationMatrix.MulPosition(obj.Triangles[triangle_id].V1.Position)
		obj.Triangles[triangle_id].V2.Position = translationMatrix.MulPosition(obj.Triangles[triangle_id].V2.Position)
		if obj.Triangles[triangle_id].V0.Normal != nil {
			n0 := translationMatrix.MulDirection(*obj.Triangles[triangle_id].V0.Normal)
			obj.Triangles[triangle_id].V0.Normal = &n0
		}
		if obj.Triangles[triangle_id].V1.Normal != nil {
			n1 := translationMatrix.MulDirection(*obj.Triangles[triangle_id].V1.Normal)
			obj.Triangles[triangle_id].V1.Normal = &n1
		}
		if obj.Triangles[triangle_id].V2.Normal != nil {
			n2 := translationMatrix.MulDirection(*obj.Triangles[triangle_id].V2.Normal)
			obj.Triangles[triangle_id].V2.Normal = &n2
		}
	}
}

func (obj *Mesh) CalculateBoundingBox() Vector {
	// init with first triangle to prohibit 0 values being min or max
	if len(obj.Triangles) == 0 {
		return Vector{}
	}
	min := obj.Triangles[0].V0.Position
	max := obj.Triangles[0].V0.Position

	for _, triangle := range obj.Triangles {
		min = triangle.V0.Position.Min(min)
		max = triangle.V0.Position.Max(max)

		min = triangle.V1.Position.Min(min)
		max = triangle.V1.Position.Max(max)

		min = triangle.V2.Position.Min(min)
		max = triangle.V2.Position.Max(max)
	}
	return Vector{
		X: max.X - min.X,
		Y: max.Y - min.Y,
		Z: max.Z - min.Z,
	}
}

func (obj *Mesh) Scale(scaling Vector) error {
	scaledVectors := make(map[Vertex]struct{})
	for triangle_id, triangle := range obj.Triangles {
		if _, exists := scaledVectors[triangle.V0]; !exists {
			obj.Triangles[triangle_id].V0.Position = triangle.V0.Position.Mult(scaling)
			scaledVectors[obj.Triangles[triangle_id].V0] = struct{}{}
		}
		if _, exists := scaledVectors[triangle.V1]; !exists {
			obj.Triangles[triangle_id].V1.Position = triangle.V1.Position.Mult(scaling)
			scaledVectors[obj.Triangles[triangle_id].V1] = struct{}{}
		}
		if _, exists := scaledVectors[triangle.V2]; !exists {
			obj.Triangles[triangle_id].V2.Position = triangle.V2.Position.Mult(scaling)
			scaledVectors[obj.Triangles[triangle_id].V2] = struct{}{}
		}
	}
	return nil
}

func (obj *Mesh) ScaleToDimensions(desiredSize *Vector) error {
	actual := obj.CalculateBoundingBox()
	if actual.X == 0 && actual.Y == 0 && actual.Z == 0 {
		return fmt.Errorf("invalid Mesh with 0 dimension")
	}
	scaling := desiredSize.Div(actual)
	return obj.Scale(scaling)
}

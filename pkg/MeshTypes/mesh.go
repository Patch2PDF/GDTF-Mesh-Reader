package MeshTypes

type Mesh struct {
	Triangles []*Triangle
}

func (obj *Mesh) AddTriangle(triangle *Triangle) {
	obj.Triangles = append(obj.Triangles, triangle)
}

func (obj *Mesh) Copy() Mesh {
	triangles := make([]*Triangle, len(obj.Triangles))
	for index, triangle := range obj.Triangles {
		temp := triangle.Copy()
		triangles[index] = &temp
	}
	return Mesh{Triangles: triangles}
}

func (obj *Mesh) Add(mesh *Mesh) *Mesh {
	obj.Triangles = append(obj.Triangles, mesh.Triangles...)
	return obj
}

func (obj *Mesh) RotateAndTranslate(translationMatrix Matrix) {
	for _, triangle := range obj.Triangles {
		triangle.V0.Position = translationMatrix.MulPosition(triangle.V0.Position)
		triangle.V1.Position = translationMatrix.MulPosition(triangle.V1.Position)
		triangle.V2.Position = translationMatrix.MulPosition(triangle.V2.Position)
	}
}

func (obj *Mesh) calculateBoundingBox() Vector {
	// init with first triangle to prohibit 0 values being min or max
	if obj.Triangles[0] == nil || obj.Triangles[0].V0 == nil {
		return Vector{}
	}
	min := obj.Triangles[0].V0.Position
	max := obj.Triangles[0].V0.Position

	for _, triangle := range obj.Triangles {
		min = triangle.V0.Position.Min(&min)
		max = triangle.V0.Position.Max(&max)

		min = triangle.V1.Position.Min(&min)
		max = triangle.V1.Position.Max(&max)

		min = triangle.V2.Position.Min(&min)
		max = triangle.V2.Position.Max(&max)
	}
	return Vector{
		X: max.X - min.X,
		Y: max.Y - min.Y,
		Z: max.Z - min.Z,
	}
}

func (obj *Mesh) ScaleToDimensions(desiredSize *Vector) {
	actual := obj.calculateBoundingBox()
	scaling := desiredSize.Div(actual)
	scaledVectors := make(map[*Vertex]struct{})
	for _, triangle := range obj.Triangles {
		if _, exists := scaledVectors[triangle.V0]; !exists {
			triangle.V0.Position = triangle.V0.Position.Mult(scaling)
			scaledVectors[triangle.V0] = struct{}{}
		}
		if _, exists := scaledVectors[triangle.V1]; !exists {
			triangle.V1.Position = triangle.V1.Position.Mult(scaling)
			scaledVectors[triangle.V1] = struct{}{}
		}
		if _, exists := scaledVectors[triangle.V2]; !exists {
			triangle.V2.Position = triangle.V2.Position.Mult(scaling)
			scaledVectors[triangle.V2] = struct{}{}
		}
	}
}

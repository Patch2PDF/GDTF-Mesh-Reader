package Types

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

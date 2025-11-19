package MeshTypes

// Datatype for Model Meshes
type Mesh struct {
	Triangles []*Triangle
}

// add Triangle to Mesh
func (obj *Mesh) AddTriangle(triangle *Triangle) {
	obj.Triangles = append(obj.Triangles, triangle)
}

// create a deep copy of Mesh object
func (obj *Mesh) Copy() Mesh {
	triangles := make([]*Triangle, len(obj.Triangles))
	for index, triangle := range obj.Triangles {
		temp := triangle.Copy()
		triangles[index] = &temp
	}
	return Mesh{Triangles: triangles}
}

// add contents from Mesh to this mesh, additionally returns the mesh pointer
func (obj *Mesh) Add(mesh *Mesh) *Mesh {
	obj.Triangles = append(obj.Triangles, mesh.Triangles...)
	return obj
}

// rotate and translate an entire mesh object
func (obj *Mesh) RotateAndTranslate(translationMatrix Matrix) {
	for _, triangle := range obj.Triangles {
		triangle.V0.Position = translationMatrix.MulPosition(triangle.V0.Position)
		triangle.V1.Position = translationMatrix.MulPosition(triangle.V1.Position)
		triangle.V2.Position = translationMatrix.MulPosition(triangle.V2.Position)
	}
}

// calculate Mesh dimension
func (obj *Mesh) calculateDimension() Vector {
	min := Vector{}
	max := Vector{}
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

// Scale mesh to desired size
func (obj *Mesh) ScaleToDimensions(desiredSize *Vector) {
	actual := obj.calculateDimension()
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

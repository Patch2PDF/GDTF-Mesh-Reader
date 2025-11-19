package MeshTypes

// Triangle data type
type Triangle struct {
	V0 *Vertex
	V1 *Vertex
	V2 *Vertex
}

// calculate normal of Triangle
func (t *Triangle) Normal() Vector {
	e1 := t.V1.Position.Sub(t.V0.Position)
	e2 := t.V2.Position.Sub(t.V0.Position)
	return e1.Cross(e2).Normalize()
}

// Create deep copy of Triangle object
func (obj *Triangle) Copy() Triangle {
	V0 := obj.V0.Copy()
	V1 := obj.V1.Copy()
	V2 := obj.V2.Copy()
	return Triangle{V0: &V0, V1: &V1, V2: &V2}
}

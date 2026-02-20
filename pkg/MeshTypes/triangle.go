package MeshTypes

type Triangle struct {
	V0 Vertex
	V1 Vertex
	V2 Vertex
}

func (t *Triangle) Normal() Vector {
	e1 := t.V1.Position.Sub(t.V0.Position)
	e2 := t.V2.Position.Sub(t.V0.Position)
	return e1.Cross(e2).Normalize()
}

func (obj *Triangle) Copy() Triangle {
	return Triangle{V0: obj.V0, V1: obj.V1, V2: obj.V2}
}

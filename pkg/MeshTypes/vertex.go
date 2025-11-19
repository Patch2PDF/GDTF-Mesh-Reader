package MeshTypes

// Vertex data type
type Vertex struct {
	Position Vector
	Normal   *Vector
}

// create deep copy of Vertex
func (obj *Vertex) Copy() Vertex {
	var normalPtr *Vector = nil
	if obj.Normal != nil {
		normalCopy := *obj.Normal
		normalPtr = &normalCopy
	}
	return Vertex{
		Position: obj.Position,
		Normal:   normalPtr,
	}
}

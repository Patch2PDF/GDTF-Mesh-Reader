package MeshTypes

type Vertex struct {
	Position Vector
	Normal   *Vector
}

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

package Types

type Vertex struct {
	Position Vector
	Normal   *Vector
}

func (obj *Vertex) Copy() Vertex {
	normalCopy := *obj.Normal
	return Vertex{
		Position: obj.Position,
		Normal:   &normalCopy,
	}
}

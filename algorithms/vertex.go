package algorithms

type Vertex struct {
	Value    string
	Children []*Vertex
	visited  bool
}

func (v *Vertex) BreadthsFirstSearch(callback func(string)) {
	queue := []*Vertex{v}
	for len(queue) != 0 {
		vertex := queue[0]
		queue = queue[1:]

		callback(vertex.Value)
		vertex.visited = true

		for _, child := range vertex.Children {
			queue = append(queue, child)
		}
	}
}

func (v *Vertex) DepthFirstSearch(callback func(string)) {
	callback(v.Value)
	v.visited = true

	for _, child := range v.Children {
		child.DepthFirstSearch(callback)
	}
}

package graph

type Edge interface {
	Weight() int
}

type edge int

func NewEdge(weight int) Edge {
	return edge(weight)
}

func (e edge) Weight() int {
	return int(e)
}

var ZeroEdge = edge(0)

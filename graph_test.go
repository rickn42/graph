package graph_test

import (
	"testing"

	"github.com/rickn42/graph"
)

func TestGraph_GetNode(t *testing.T) {
	g := graph.NewGraph()

	id1 := graph.NewId(1)
	n1 := graph.NewNode(id1)

	if !n1.Invalid() {
		t.Error("Invalid node before Graph.AddNode()")
	}

	if err := g.AddNode(n1); err != nil {
		t.Error(err)
	}

	if n1.Invalid() {
		t.Error("Valid node after Graph.AddNode()")
	}

	if n, ok := g.GetNode(graph.NewId(1)); !ok || n != n1 {
		t.Errorf("GetNode failed.")
	}
}

func TestGraph_Connect(t *testing.T) {
	g := graph.NewGraph()

	// node1
	id1 := graph.NewId(1)
	n1 := graph.NewNode(id1)
	g.AddNode(n1)

	// node2
	id2 := graph.NewId(2)
	n2 := graph.NewNode(id2)
	g.AddNode(n2)

	// add edge
	err := g.AddEdge(n1.Id(), n2.Id())
	if err != nil {
		t.Error(err)
	}

	if ns := n1.Neighbors(); len(ns) != 1 || ns[0].Node != n2 {
		t.Error("AddEdge failed")
	}
	if ns := n2.Neighbors(); len(ns) != 1 || ns[0].Node != n1 {
		t.Error("AddEdge failed")
	}

	g.RemoveEdge(n1.Id(), n2.Id())
	if ns := n1.Neighbors(); len(ns) != 0 {
		t.Error("RemoveEdge failed")
	}
	if ns := n2.Neighbors(); len(ns) != 0 {
		t.Error("RemoveEdge failed")
	}

	g.AddEdge(n1.Id(), n2.Id())

	g.RemoveNode(n2.Id())
	if ns := n1.Neighbors(); len(ns) != 0 {
		t.Error("RemoveNode failed")
	}

	if !n2.Invalid() {
		t.Error("Invalid check after RemoveNode failed")
	}
}

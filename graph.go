package graph

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
)

var (
	DevTools        bool
	AlreadyIdExists = errors.New("Id exists already")
	NodeNotExists   = errors.New("Node not exists")
)

// Graph is a uni-edge, non-directed simple graph.
//
// This use interface-type node.
// So, easily make custom node by embedding node implementation.
type Graph struct {
	ns sync.Map
	mu sync.Mutex
	n  int // Last node lock order number
}

func NewGraph() *Graph {
	return &Graph{
		ns: sync.Map{},
	}
}

func (g *Graph) GetNode(id interface{}) (n Node, ok bool) {
	var key interface{}
	if i, ok := id.(Id); ok {
		key = i.Value()
	} else {
		key = id
	}
	val, ok := g.ns.Load(key)
	return val.(Node), ok
}

func (g *Graph) AddNode(n Node) error {
	n.lock()
	defer n.unlock()

	err := n.init(g, g.genLockNum())
	if err != nil {
		return err
	}

	_, loaded := g.ns.LoadOrStore(n.Id().Value(), n)
	if loaded {
		n.resetInit()
		return AlreadyIdExists
	}

	if DevTools {
		fmt.Println("Added", n.toString())
	}
	return nil
}

func (g *Graph) genLockNum() (n int) {
	g.mu.Lock()
	n = g.n + 1
	g.n = n
	g.mu.Unlock()
	return n
}

func (g *Graph) RemoveNode(id Id) (err error) {
	key := id.Value()

	val, ok := g.ns.Load(key)
	if !ok {
		return NodeNotExists
	}

	g.ns.Delete(key)

	n := val.(Node)
	err = n.Remove()

	if DevTools {
		fmt.Println("Removed", n.toString())
	}
	return err
}

func (g *Graph) AddEdge(id1, id2 interface{}, e ...Edge) error {

	var key interface{}
	if id, ok := id1.(Id); ok {
		key = id.Value()
	} else {
		key = id1
	}

	v1, ok := g.ns.Load(key)
	if !ok {
		return NodeNotExists
	}

	if id, ok := id2.(Id); ok {
		key = id.Value()
	} else {
		key = id2
	}
	v2, ok := g.ns.Load(key)
	if !ok {
		return NodeNotExists
	}

	var edge Edge
	if len(e) == 0 {
		edge = ZeroEdge
	} else {
		edge = e[len(e)-1]
	}
	return v1.(Node).AddEdge(v2.(Node), edge)
}

func (g *Graph) RemoveEdge(id1, id2 interface{}) (err error) {

	var key interface{}
	if id, ok := id1.(Id); ok {
		key = id.Value()
	} else {
		key = id1
	}
	v1, ok := g.ns.Load(key)
	if !ok {
		return NodeNotExists
	}

	if id, ok := id2.(Id); ok {
		key = id.Value()
	} else {
		key = id2
	}
	v2, ok := g.ns.Load(key)
	if !ok {
		return NodeNotExists
	}

	return v1.(Node).RemoveEdge(v2.(Node))
}

func (g *Graph) String() string {
	var buf bytes.Buffer
	var cnt int

	g.ns.Range(func(_, val interface{}) bool {
		n := val.(Node)
		buf.WriteString("\n\t")
		buf.WriteString(n.String())
		cnt++
		return true
	})

	return fmt.Sprintf("Graph (Node cnt = %d)%s", cnt, buf.String())
}

package graph

import (
	"errors"
	"fmt"
	"sync"
)

var (
	InvalidNode        = errors.New("Invalid node")
	AlreadyNodeRemoved = errors.New("already removed node")
	AlreadyNodeAdded   = errors.New("Node already added")
	DifferentGraphNode = errors.New("Node graph is different")
	DifferentGraph     = errors.New("different graph")
)

type Neighbor struct {
	Node
	Edge
}

func (n Neighbor) String() string {
	return fmt.Sprintf("Neighbor %s Edge(%d)", n.Node.String(), n.Edge.Weight())
}

type Node interface {
	lock()
	unlock()
	lockNum() int
	init(g *Graph, lockOrderNum int) error
	resetInit()
	graph() *Graph
	addNeighbor(Node, Edge)
	removeNeighbor(Node)
	toString() string

	Id() Id
	Invalid() bool
	AddEdge(Node, Edge) error
	RemoveEdge(Node) error
	Neighbors() []Neighbor
	IsNeighbor(Node) bool
	Remove() error
	String() string
}

var _ Node = new(nodeImpl)

type nodeImpl struct {
	id Id
	ns sync.Map
	mu sync.RWMutex
	g  *Graph
	n  int // node lock order number
}

func NewNode(id Id) *nodeImpl {
	return &nodeImpl{
		id: id,
		ns: sync.Map{},
	}
}

func (n *nodeImpl) lock() {
	n.mu.Lock()
}

func (n *nodeImpl) unlock() {
	n.mu.Unlock()
}

func (n *nodeImpl) init(g *Graph, num int) error {
	if n.g != nil {
		if n.g != g {
			return AlreadyNodeAdded
		} else {
			return DifferentGraphNode
		}
	}
	n.g = g
	n.n = num

	return nil
}

func (n *nodeImpl) resetInit() {
	n.g = nil
	n.ns = sync.Map{}
}

func (n *nodeImpl) lockNum() int {
	return n.n
}

func (n *nodeImpl) graph() *Graph {
	return n.g
}

func (n *nodeImpl) addNeighbor(n2 Node, e Edge) {
	n.ns.Store(n2, e)
}

func (n *nodeImpl) removeNeighbor(n2 Node) {
	n.ns.Delete(n2)
}

func (n *nodeImpl) toString() string {
	if n.g == nil {
		return "Node(invalid)"
	}

	return fmt.Sprintf("Node(%s)", n.id)
}

func (n *nodeImpl) Id() Id {
	return n.id
}

func (n *nodeImpl) Invalid() (invalid bool) {
	n.mu.RLock()
	invalid = n.g == nil
	n.mu.RUnlock()
	return invalid
}

func (n *nodeImpl) AddEdge(n2 Node, e Edge) error {
	lockNodes(n, n2)
	defer unlockNodes(n, n2)

	if n.graph() == nil || n2.graph() == nil {
		return InvalidNode
	}

	if n.graph() != n2.graph() {
		return DifferentGraph
	}

	n.addNeighbor(n2, e)
	n2.addNeighbor(n, e)

	if DevTools {
		fmt.Printf("Add Edge: -*-\n\t%-6s%s\n\t%-6s%s)\n", "node1:", n.toString(), "node2:", n2.toString())
	}
	return nil
}

func (n *nodeImpl) RemoveEdge(n2 Node) error {
	lockNodes(n, n2)
	defer unlockNodes(n, n2)

	n.removeNeighbor(n2)
	n2.removeNeighbor(n)

	if DevTools {
		fmt.Printf("Disonnected: -/-\n\t%-6s%s\n\t%-6s%s)\n", "from:", n.toString(), "to:", n2.toString())
	}
	return nil
}

func (n *nodeImpl) Neighbors() (ns []Neighbor) {
	n.ns.Range(func(key, value interface{}) bool {
		ns = append(ns, Neighbor{key.(Node), value.(Edge)})
		return true
	})
	return ns
}

func (n *nodeImpl) IsNeighbor(n2 Node) (ok bool) {
	_, ok = n.ns.Load(n2)
	return ok
}

func (n *nodeImpl) Remove() error {
	n.ns.Range(func(key, value interface{}) bool {
		n.RemoveEdge(key.(Node))
		return true
	})

	n.mu.Lock()
	n.g = nil
	n.mu.Unlock()

	return nil
}

func (n *nodeImpl) String() string {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return n.toString()
}

func (n *nodeImpl) StringWithNeighbors() string {
	n.mu.RLock()
	defer n.mu.RUnlock()

	if n.g == nil {
		return n.String()
	}

	var s string
	var cnt int
	n.ns.Range(func(key, value interface{}) bool {
		s += "\tnode: " + key.(Node).String() + "\n"
		cnt++
		return true
	})

	return fmt.Sprintf("%s DownNodes: %d\n%s", n, cnt, s)
}

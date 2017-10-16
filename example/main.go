package main

import (
	"fmt"

	"github.com/rickn42/graph"
)

func main() {
	graph.DevTools = true

	cities := graph.NewGraph()
	cities.AddNode(NewCity("seoul", 605.21, 10204057))
	cities.AddNode(NewCity("pusan", 765.82, 3357716))
	cities.AddEdge("seoul", "pusan", graph.NewEdge(300))

	seoul, _ := cities.GetNode("seoul")
	neighbor := seoul.Neighbors()[0]

	// need to convert graph.Node to custom-type node
	city := neighbor.Node.(*City)
	distance := neighbor.Edge.Weight()

	fmt.Println(distance)          // 300
	fmt.Println(city.Id().Value()) // pusan
	fmt.Println(city.Size)         // 765.82 (pusan.Size)
	fmt.Println(city.Population)   // 3357716 (pusan.Population)
}

// City is custom node made by embedding graph.Node
type City struct {
	graph.Node
	Size       float32
	Population int
}

func NewCity(id interface{}, size float32, population int) *City {
	return &City{
		Node:       graph.NewNode(graph.NewId(id)),
		Size:       size,
		Population: population,
	}
}

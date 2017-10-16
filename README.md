# Graph

simple graph (uni-edge, non-directed)

### example

custom-node made by embedding
```go
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
```

use custom-node
```go
cities := graph.NewGraph()
cities.AddNode(NewCity("seoul", 605.21, 10204057))
cities.AddNode(NewCity("pusan", 765.82, 3357716))
cities.AddEdge("seoul", "pusan", graph.NewEdge(300))

seoul, _ := cities.GetNode("seoul")
neighbor := seoul.Neighbors()[0]

// need to convert graph.Node to custom-type node
city := neighbor.Node.(*City)
distance := neighbor.Edge.Weight()
```

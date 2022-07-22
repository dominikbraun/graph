package draw

import (
	"testing"

	"github.com/dominikbraun/graph"
)

func TestGraph(t *testing.T) {
	g := graph.New(graph.IntHash)
	g.Vertex(1)
	g.Vertex(2)
	g.Vertex(3)
	g.Vertex(4)

	g.Edge(1, 2)
	g.Edge(2, 3)
	g.Edge(2, 4)

	Graph(g)
}

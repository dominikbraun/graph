package graph

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestDirectedTopologicalSort(t *testing.T) {
	tests := map[string]struct {
		vertices      []int
		edges         []Edge[int]
		expectedOrder []int
		shouldFail    bool
	}{
		"graph with 5 vertices": {
			vertices: []int{1, 2, 3, 4, 5},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
				{Source: 2, Target: 3},
				{Source: 2, Target: 4},
				{Source: 2, Target: 5},
				{Source: 3, Target: 4},
				{Source: 4, Target: 5},
			},
			expectedOrder: []int{1, 2, 3, 4, 5},
		},
		"graph with many possible topological orders": {
			vertices: []int{1, 2, 3, 4, 5, 6, 10, 20, 30, 40, 50, 60},
			edges: []Edge[int]{
				{Source: 1, Target: 10},
				{Source: 2, Target: 20},
				{Source: 3, Target: 30},
				{Source: 4, Target: 40},
				{Source: 5, Target: 50},
				{Source: 6, Target: 60},
			},
		},
		"graph with cycle": {
			vertices: []int{1, 2, 3},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 2, Target: 3},
				{Source: 3, Target: 1},
			},
			shouldFail: true,
		},
	}

	for name, test := range tests {
		graph := New(IntHash, Directed())

		err := buildGraph(&graph, test.vertices, test.edges)
		if err != nil {
			t.Fatalf("%s: failed to construct graph: %s", name, err.Error())
		}

		order, err := TopologicalSort(graph)

		if test.shouldFail != (err != nil) {
			t.Errorf("%s: error expectancy doesn't match: expected %v, got %v (error: %v)", name, test.shouldFail, err != nil, err)
		}

		if test.shouldFail {
			continue
		}

		if len(order) != len(test.vertices) {
			t.Errorf("%s: order length expectancy doesn't match: expected %v, got %v", name, len(test.vertices), len(order))
		}

		if len(test.expectedOrder) <= 0 {

			fmt.Println("topological sort", order)

			if err := verifyTopologicalSort(graph, order); err != nil {
				t.Errorf("%s: invalid topological sort - %v", name, err)
			}
		}

		for i, expectedVertex := range test.expectedOrder {
			if expectedVertex != order[i] {
				t.Errorf("%s: order expectancy doesn't match: expected %v at %d, got %v", name, expectedVertex, i, order[i])
			}
		}
	}
}

func TestUndirectedTopologicalSort(t *testing.T) {
	tests := map[string]struct {
		expectedOrder []int
		shouldFail    bool
	}{
		"return error": {
			expectedOrder: nil,
			shouldFail:    true,
		},
	}

	for name, test := range tests {
		graph := New(IntHash)

		order, err := TopologicalSort(graph)

		if test.shouldFail != (err != nil) {
			t.Errorf("%s: error expectancy doesn't match: expected %v, got %v (error: %v)", name, test.shouldFail, err != nil, err)
		}

		if test.expectedOrder == nil && order != nil {
			t.Errorf("%s: order expectancy doesn't match: expcted %v, got %v", name, test.expectedOrder, order)
		}
	}
}

func TestDirectedStableTopologicalSort(t *testing.T) {
	tests := map[string]struct {
		vertices      []int
		edges         []Edge[int]
		expectedOrder []int
		shouldFail    bool
	}{
		"graph with 5 vertices": {
			vertices: []int{1, 2, 3, 4, 5},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
				{Source: 2, Target: 3},
				{Source: 2, Target: 4},
				{Source: 2, Target: 5},
				{Source: 3, Target: 4},
				{Source: 4, Target: 5},
			},
			expectedOrder: []int{1, 2, 3, 4, 5},
		},
		"graph with many possible topological orders": {
			vertices: []int{1, 2, 3, 4, 5, 6, 10, 20, 30, 40, 50, 60},
			edges: []Edge[int]{
				{Source: 1, Target: 10},
				{Source: 2, Target: 20},
				{Source: 3, Target: 30},
				{Source: 4, Target: 40},
				{Source: 5, Target: 50},
				{Source: 6, Target: 60},
			},
			expectedOrder: []int{1, 2, 3, 4, 5, 6, 10, 20, 30, 40, 50, 60},
		},
		"graph with cycle": {
			vertices: []int{1, 2, 3},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 2, Target: 3},
				{Source: 3, Target: 1},
			},
			shouldFail: true,
		},
	}

	for name, test := range tests {
		graph := New(IntHash, Directed())

		err := buildGraph(&graph, test.vertices, test.edges)
		if err != nil {
			t.Fatalf("%s: failed to construct graph: %s", name, err.Error())
		}

		order, err := StableTopologicalSort(graph, func(a, b int) bool {
			return a < b
		})

		if test.shouldFail != (err != nil) {
			t.Errorf("%s: error expectancy doesn't match: expected %v, got %v (error: %v)", name, test.shouldFail, err != nil, err)
		}

		if test.shouldFail {
			continue
		}

		if len(order) != len(test.expectedOrder) {
			t.Errorf("%s: order length expectancy doesn't match: expected %v, got %v", name, len(test.expectedOrder), len(order))
		}

		fmt.Println("expected", test.expectedOrder)
		fmt.Println("actual", order)

		for i, expectedVertex := range test.expectedOrder {
			if expectedVertex != order[i] {
				t.Errorf("%s: order expectancy doesn't match: expected %v at %d, got %v", name, expectedVertex, i, order[i])
			}
		}
	}
}

func TestDirectedTransitiveReduction(t *testing.T) {
	tests := map[string]struct {
		vertices      []string
		edges         []Edge[string]
		expectedEdges []Edge[string]
		shouldFail    bool
	}{
		"graph as on img/transitive-reduction-before.svg": {
			vertices: []string{"A", "B", "C", "D", "E"},
			edges: []Edge[string]{
				{Source: "A", Target: "B"},
				{Source: "A", Target: "C"},
				{Source: "A", Target: "D"},
				{Source: "A", Target: "E"},
				{Source: "B", Target: "D"},
				{Source: "C", Target: "D"},
				{Source: "C", Target: "E"},
				{Source: "D", Target: "E"},
			},
			expectedEdges: []Edge[string]{
				{Source: "A", Target: "B"},
				{Source: "A", Target: "C"},
				{Source: "B", Target: "D"},
				{Source: "C", Target: "D"},
				{Source: "D", Target: "E"},
			},
		},
		"graph with cycle": {
			vertices: []string{"A", "B", "C"},
			edges: []Edge[string]{
				{Source: "A", Target: "B"},
				{Source: "B", Target: "C"},
				{Source: "C", Target: "A"},
			},
			shouldFail: true,
		},
		"graph from issue 83": {
			vertices: []string{"_root", "A", "B", "C", "D", "E", "F"},
			edges: []Edge[string]{
				{Source: "_root", Target: "A"},
				{Source: "_root", Target: "B"},
				{Source: "_root", Target: "C"},
				{Source: "_root", Target: "D"},
				{Source: "_root", Target: "E"},
				{Source: "_root", Target: "F"},
				{Source: "E", Target: "C"},
				{Source: "F", Target: "D"},
				{Source: "F", Target: "C"},
				{Source: "F", Target: "E"},
				{Source: "C", Target: "A"},
				{Source: "C", Target: "B"},
			},
			expectedEdges: []Edge[string]{
				{Source: "_root", Target: "F"},
				{Source: "F", Target: "D"},
				{Source: "F", Target: "E"},
				{Source: "E", Target: "C"},
				{Source: "C", Target: "A"},
				{Source: "C", Target: "B"},
			},
		},
	}

	for name, test := range tests {
		graph := New(StringHash, Directed())

		err := buildGraph(&graph, test.vertices, test.edges)
		if err != nil {
			t.Fatalf("%s: failed to construct graph: %s", name, err.Error())
		}

		reduction, err := TransitiveReduction(graph)

		if test.shouldFail != (err != nil) {
			t.Errorf("%s: error expectancy doesn't match: expected %v, got %v (error: %v)", name, test.shouldFail, err != nil, err)
		}

		if test.shouldFail {
			continue
		}

		actualEdges := make([]Edge[string], 0)
		adjacencyMap, _ := reduction.AdjacencyMap()

		for _, adjacencies := range adjacencyMap {
			for _, edge := range adjacencies {
				actualEdges = append(actualEdges, edge)
			}
		}

		equalsFunc := reduction.(*directed[string, string]).edgesAreEqual

		if !slicesAreEqualWithFunc(actualEdges, test.expectedEdges, equalsFunc) {
			t.Errorf("%s: edge expectancy doesn't match: expected %v, got %v", name, test.expectedEdges, actualEdges)
		}
	}
}

func TestUndirectedTransitiveReduction(t *testing.T) {
	tests := map[string]struct {
		shouldFail bool
	}{
		"return error": {
			shouldFail: true,
		},
	}

	for name, test := range tests {
		graph := New(StringHash)

		_, err := TransitiveReduction(graph)

		if test.shouldFail != (err != nil) {
			t.Errorf("%s: error expectancy doesn't match: expected %v, got %v (error: %v)", name, test.shouldFail, err != nil, err)
		}
	}
}

func TestVerifyTopologicalSort(t *testing.T) {
	tests := map[string]struct {
		vertices      []int
		edges         []Edge[int]
		invalidOrder []int
	}{
		"graph with 2 vertices": {
			vertices: []int{1, 2},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
			},
		},
		"graph with 2 vertices - reversed": {
			vertices: []int{1, 2},
			edges: []Edge[int]{
				{Source: 2, Target: 1},
			},
		},
		"graph with 2 vertices - invalid": {
			vertices: []int{1, 2},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
			},
			invalidOrder: []int{2, 1},
		},
		"graph with 3 vertices": {
			vertices: []int{1, 2, 3},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
				{Source: 2, Target: 3},
			},
		},
		"graph with 3 vertices - invalid": {
			vertices: []int{1, 2, 3},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
				{Source: 2, Target: 3},
			},
			invalidOrder: []int{1, 3, 2},
		},
		"graph with 5 vertices": {
			vertices: []int{1, 2, 3, 4, 5},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
				{Source: 2, Target: 3},
				{Source: 2, Target: 4},
				{Source: 2, Target: 5},
				{Source: 3, Target: 4},
				{Source: 4, Target: 5},
			},
		},
		"graph with many possible topological orders": {
			vertices: []int{1, 2, 3, 4, 5, 6, 10, 20, 30, 40, 50, 60},
			edges: []Edge[int]{
				{Source: 1, Target: 10},
				{Source: 2, Target: 20},
				{Source: 3, Target: 30},
				{Source: 4, Target: 40},
				{Source: 5, Target: 50},
				{Source: 6, Target: 60},
			},
			invalidOrder: []int{2, 3, 4, 5, 6, 10, 1, 20, 30, 40, 50, 60},
		},
	}

	for name, test := range tests {
		graph := New[int, int](IntHash, Directed())

		err := buildGraph(&graph, test.vertices, test.edges)
		if err != nil {
			t.Fatalf("%s: failed to construct graph: %s", name, err.Error())
		}

		var order[] int

		if len(test.invalidOrder) > 0 {
			order = test.invalidOrder
		} else {
			order, err = TopologicalSort(graph)
			if err != nil {
				t.Fatalf("%s: error failed to produce topological sort: %v)", name, err)
			}
		}

		err = verifyTopologicalSort(graph, order)

		shouldFail := len(test.invalidOrder) > 0

		if shouldFail != (err != nil) {
			t.Errorf("%s: error expectancy doesn't match: expected %v, got %v (error: %v)", name, shouldFail, err != nil, err)
		}
	}
}

func slicesAreEqualWithFunc[T any](a, b []T, equals func(a, b T) bool) bool {
	if len(a) != len(b) {
		return false
	}

	for _, aValue := range a {
		found := false
		for _, bValue := range b {
			if equals(aValue, bValue) {
				found = true
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// Please note that this call is destructive.  Make a clone of your graph before calling if you
// wish to preserve the graph.
func verifyTopologicalSort[K comparable, T any](graph Graph[K, T], order []K) error {

	adjacencyMap, err := graph.AdjacencyMap()
	if err != nil {
		return fmt.Errorf("failed to get adjacency map: %v", err)
	}

	for i := range order {

		for _, edge := range adjacencyMap[order[i]] {
			err = graph.RemoveEdge(edge.Source, edge.Target)
			if err != nil {
				return fmt.Errorf("failed to remove edge: %v -> %v : %v", edge.Source, edge.Target, err)
			}
		}

		err = graph.RemoveVertex(order[i])
		if err != nil {
			return fmt.Errorf("failed to remove vertex: %v at index %d: %v", order[i], i, err)
		}
	}

	return nil
}

// randomizes the ordering of the edges and vertices to help ferret out any potential bugs
// related to ordering
func buildGraph[K comparable, T any](g *Graph[K, T], vertices []T, edges []Edge[K]) error {

	if g == nil {
		return fmt.Errorf("graph must be initialized")
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(vertices), func(i, j int) { vertices[i], vertices[j] = vertices[j], vertices[i] })

	for _, vertex := range vertices {
		_ = (*g).AddVertex(vertex)
	}

	rand.Shuffle(len(edges), func(i, j int) { edges[i], edges[j] = edges[j], edges[i] })

	for _, edge := range edges {
		if err := (*g).AddEdge(edge.Source, edge.Target, EdgeWeight(edge.Properties.Weight)); err != nil {
			return err
		}
	}

	return nil
}

package graph

import (
	"fmt"
	"testing"
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

		for _, vertex := range test.vertices {
			_ = graph.AddVertex(vertex)
		}

		for _, edge := range test.edges {
			if err := graph.AddEdge(edge.Source, edge.Target, EdgeWeight(edge.Properties.Weight)); err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		order, err := TopologicalSort(graph)

		if test.shouldFail != (err != nil) {
			t.Errorf("%s: error expectancy doesn't match: expected %v, got %v (error: %v)", name, test.shouldFail, err != nil, err)
		}

		if test.shouldFail {
			continue
		}

		if len(order) != len(test.expectedOrder) {
			t.Errorf("%s: order length expectancy doesn't match: expected %v, got %v", name, len(test.expectedOrder), len(order))
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

		for _, vertex := range test.vertices {
			_ = graph.AddVertex(vertex)
		}

		for _, edge := range test.edges {
			if err := graph.AddEdge(edge.Source, edge.Target, EdgeWeight(edge.Properties.Weight)); err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
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

		for _, vertex := range test.vertices {
			_ = graph.AddVertex(vertex)
		}

		for _, edge := range test.edges {
			if err := graph.AddEdge(edge.Source, edge.Target, EdgeWeight(edge.Properties.Weight)); err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
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

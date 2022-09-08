package graph

import "testing"

func TestDirectedTopologicalSort(t *testing.T) {
	tests := map[string]struct {
		vertices      []int
		edges         []Edge[int]
		expectedOrder []int
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
	}

	for name, test := range tests {
		graph := New(IntHash, Directed(), Acyclic())

		for _, vertex := range test.vertices {
			_ = graph.AddVertex(vertex)
		}

		for _, edge := range test.edges {
			if err := graph.AddEdge(edge.Source, edge.Target, EdgeWeight(edge.Properties.Weight)); err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		order, err := TopologicalSort(graph)
		if err != nil {
			t.Fatalf("%s: failed to add edge: %s", name, err.Error())
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

func TestDirectedTransitiveReduction(t *testing.T) {
	tests := map[string]struct {
		vertices      []string
		edges         []Edge[string]
		expectedEdges []Edge[string]
	}{
		"graph as on img/transitive-reduction.svg": {
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
	}

	for name, test := range tests {
		graph := New(StringHash, Directed(), Acyclic())

		for _, vertex := range test.vertices {
			_ = graph.AddVertex(vertex)
		}

		for _, edge := range test.edges {
			if err := graph.AddEdge(edge.Source, edge.Target, EdgeWeight(edge.Properties.Weight)); err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		if err := TransitiveReduction(graph); err != nil {
			t.Fatalf("%s: failed to add edge: %s", name, err.Error())
		}

		actualEdges := make([]Edge[string], 0)
		adjacencyMap, _ := graph.AdjacencyMap()

		for _, adjacencies := range adjacencyMap {
			for _, edge := range adjacencies {
				actualEdges = append(actualEdges, edge)
			}
		}

		equalsFunc := graph.(*directed[string, string]).edgesAreEqual

		if !slicesAreEqualWithFunc(actualEdges, test.expectedEdges, equalsFunc) {
			t.Errorf("%s: edge expectancy doesn't match: expected %v, got %v", name, test.expectedEdges, actualEdges)
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

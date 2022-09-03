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

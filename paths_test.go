package graph

import "testing"

func TestDirectedCreatesCycle(t *testing.T) {
	tests := map[string]struct {
		vertices     []int
		edges        []Edge[int]
		sourceHash   int
		targetHash   int
		createsCycle bool
	}{
		"directed 2-4-7-5 cycle": {
			vertices: []int{1, 2, 3, 4, 5, 6, 7},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
				{Source: 2, Target: 4},
				{Source: 3, Target: 6},
				{Source: 4, Target: 7},
				{Source: 5, Target: 2},
			},
			sourceHash:   7,
			targetHash:   5,
			createsCycle: true,
		},
		"undirected 2-4-7-5 'cycle'": {
			vertices: []int{1, 2, 3, 4, 5, 6, 7},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
				{Source: 2, Target: 4},
				{Source: 3, Target: 6},
				{Source: 4, Target: 7},
				{Source: 5, Target: 2},
			},
			sourceHash: 5,
			targetHash: 7,
			// The direction of the edge (57 instead of 75) doesn't create a directed cycle.
			createsCycle: false,
		},
		"no cycle": {
			vertices: []int{1, 2, 3, 4, 5, 6, 7},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
				{Source: 2, Target: 4},
				{Source: 3, Target: 6},
				{Source: 4, Target: 7},
				{Source: 5, Target: 2},
			},
			sourceHash:   5,
			targetHash:   6,
			createsCycle: false,
		},
	}

	for name, test := range tests {
		graph := New(IntHash, Directed())

		for _, vertex := range test.vertices {
			graph.AddVertex(vertex)
		}

		for _, edge := range test.edges {
			if err := graph.AddEdge(edge.Source, edge.Target); err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		createsCycle, err := CreatesCycle(graph, test.sourceHash, test.targetHash)
		if err != nil {
			t.Fatalf("%s: failed to add edge: %s", name, err.Error())
		}

		if createsCycle != test.createsCycle {
			t.Errorf("%s: cycle expectancy doesn't match: expected %v, got %v", name, test.createsCycle, createsCycle)
		}
	}
}

func TestUndirectedCreatesCycle(t *testing.T) {
	tests := map[string]struct {
		vertices     []int
		edges        []Edge[int]
		sourceHash   int
		targetHash   int
		createsCycle bool
	}{
		"undirected 2-4-7-5 cycle": {
			vertices: []int{1, 2, 3, 4, 5, 6, 7},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
				{Source: 2, Target: 4},
				{Source: 3, Target: 6},
				{Source: 4, Target: 7},
				{Source: 5, Target: 7},
			},
			sourceHash:   2,
			targetHash:   5,
			createsCycle: true,
		},
		"undirected 5-6-3-1-2-7 cycle": {
			vertices: []int{1, 2, 3, 4, 5, 6, 7},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
				{Source: 2, Target: 4},
				{Source: 3, Target: 6},
				{Source: 4, Target: 7},
				{Source: 5, Target: 7},
			},
			sourceHash:   5,
			targetHash:   6,
			createsCycle: true,
		},
		"no cycle": {
			vertices: []int{1, 2, 3, 4, 5, 6, 7},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
				{Source: 2, Target: 4},
				{Source: 3, Target: 6},
				{Source: 4, Target: 7},
			},
			sourceHash:   5,
			targetHash:   7,
			createsCycle: false,
		},
	}

	for name, test := range tests {
		graph := New(IntHash)

		for _, vertex := range test.vertices {
			graph.AddVertex(vertex)
		}

		for _, edge := range test.edges {
			if err := graph.AddEdge(edge.Source, edge.Target); err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		createsCycle, err := CreatesCycle(graph, test.sourceHash, test.targetHash)
		if err != nil {
			t.Fatalf("%s: failed to add edge: %s", name, err.Error())
		}

		if createsCycle != test.createsCycle {
			t.Errorf("%s: cycle expectancy doesn't match: expected %v, got %v", name, test.createsCycle, createsCycle)
		}
	}
}

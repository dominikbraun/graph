package graph

import (
	"testing"
)

func TestDirectedMinimumSpanningTree(t *testing.T) {
	tests := map[string]struct {
		shouldFail bool
	}{
		"returns error": {
			shouldFail: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			graph := New(IntHash, Directed())

			_, err := MinimumSpanningTree(graph)

			if test.shouldFail != (err != nil) {
				t.Errorf("expected error == %v, got %v", test.shouldFail, err)
			}
		})
	}
}

func TestUndirectedMinimumSpanningTree(t *testing.T) {
	tests := map[string]struct {
		vertices                []string
		edges                   []Edge[string]
		expectedErr             error
		expectedMSTAdjacencyMap map[string]map[string]Edge[string]
	}{
		"graph from img/mst.svg": {
			vertices: []string{"A", "B", "C", "D"},
			edges: []Edge[string]{
				{Source: "A", Target: "B", Properties: EdgeProperties{Weight: 2}},
				{Source: "A", Target: "C", Properties: EdgeProperties{Weight: 4}},
				{Source: "A", Target: "D", Properties: EdgeProperties{Weight: 3}},
				{Source: "B", Target: "C", Properties: EdgeProperties{Weight: 4}},
				{Source: "B", Target: "D", Properties: EdgeProperties{Weight: 1}},
				{Source: "C", Target: "D", Properties: EdgeProperties{Weight: 3}},
			},
			expectedErr: nil,
			expectedMSTAdjacencyMap: map[string]map[string]Edge[string]{
				"A": {
					"B": {Source: "A", Target: "B", Properties: EdgeProperties{Weight: 2}},
				},
				"B": {
					"D": {Source: "B", Target: "D", Properties: EdgeProperties{Weight: 1}},
					"A": {Source: "B", Target: "A", Properties: EdgeProperties{Weight: 2}},
				},
				"C": {
					"D": {Source: "C", Target: "D", Properties: EdgeProperties{Weight: 3}},
				},
				"D": {
					"B": {Source: "D", Target: "B", Properties: EdgeProperties{Weight: 1}},
					"C": {Source: "D", Target: "C", Properties: EdgeProperties{Weight: 3}},
				},
			},
		},
		"two trees for a disconnected graph": {
			vertices: []string{"A", "B", "C", "D"},
			edges: []Edge[string]{
				{Source: "A", Target: "B", Properties: EdgeProperties{Weight: 2}},
				{Source: "C", Target: "D", Properties: EdgeProperties{Weight: 4}},
			},
			expectedErr: nil,
			expectedMSTAdjacencyMap: map[string]map[string]Edge[string]{
				"A": {
					"B": {Source: "A", Target: "B", Properties: EdgeProperties{Weight: 2}},
				},
				"B": {
					"A": {Source: "B", Target: "A", Properties: EdgeProperties{Weight: 2}},
				},
				"C": {
					"D": {Source: "C", Target: "D", Properties: EdgeProperties{Weight: 4}},
				},
				"D": {
					"C": {Source: "D", Target: "C", Properties: EdgeProperties{Weight: 4}},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			g := New(StringHash)

			for _, vertex := range test.vertices {
				_ = g.AddVertex(vertex)
			}

			for _, edge := range test.edges {
				_ = g.AddEdge(copyEdge(edge))
			}

			mst, _ := MinimumSpanningTree(g)
			adjacencyMap, _ := mst.AdjacencyMap()

			edgesAreEqual := g.(*undirected[string, string]).edgesAreEqual

			if !adjacencyMapsAreEqual(test.expectedMSTAdjacencyMap, adjacencyMap, edgesAreEqual) {
				t.Fatalf("expected adjacency map %v, got %v", test.expectedMSTAdjacencyMap, adjacencyMap)
			}
		})
	}
}

func TestDirectedMaximumSpanningTree(t *testing.T) {
	tests := map[string]struct {
		shouldFail bool
	}{
		"returns error": {
			shouldFail: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			graph := New(IntHash, Directed())

			_, err := MaximumSpanningTree(graph)

			if test.shouldFail != (err != nil) {
				t.Errorf("expected error == %v, got %v", test.shouldFail, err)
			}
		})
	}
}

func TestUndirectedMaximumSpanningTree(t *testing.T) {
	tests := map[string]struct {
		vertices                []string
		edges                   []Edge[string]
		expectedErr             error
		expectedMSTAdjacencyMap map[string]map[string]Edge[string]
	}{
		"graph from img/mst.svg with higher weights": {
			vertices: []string{"A", "B", "C", "D"},
			edges: []Edge[string]{
				{Source: "A", Target: "B", Properties: EdgeProperties{Weight: 20}},
				{Source: "A", Target: "C", Properties: EdgeProperties{Weight: 4}},
				{Source: "A", Target: "D", Properties: EdgeProperties{Weight: 3}},
				{Source: "B", Target: "C", Properties: EdgeProperties{Weight: 4}},
				{Source: "B", Target: "D", Properties: EdgeProperties{Weight: 10}},
				{Source: "C", Target: "D", Properties: EdgeProperties{Weight: 30}},
			},
			expectedErr: nil,
			expectedMSTAdjacencyMap: map[string]map[string]Edge[string]{
				"A": {
					"B": {Source: "A", Target: "B", Properties: EdgeProperties{Weight: 20}},
				},
				"B": {
					"D": {Source: "B", Target: "D", Properties: EdgeProperties{Weight: 10}},
					"A": {Source: "B", Target: "A", Properties: EdgeProperties{Weight: 20}},
				},
				"C": {
					"D": {Source: "C", Target: "D", Properties: EdgeProperties{Weight: 30}},
				},
				"D": {
					"B": {Source: "D", Target: "B", Properties: EdgeProperties{Weight: 10}},
					"C": {Source: "D", Target: "C", Properties: EdgeProperties{Weight: 30}},
				},
			},
		},
		"two trees for a disconnected graph": {
			vertices: []string{"A", "B", "C", "D"},
			edges: []Edge[string]{
				{Source: "A", Target: "B", Properties: EdgeProperties{Weight: 2}},
				{Source: "C", Target: "D", Properties: EdgeProperties{Weight: 4}},
			},
			expectedErr: nil,
			expectedMSTAdjacencyMap: map[string]map[string]Edge[string]{
				"A": {
					"B": {Source: "A", Target: "B", Properties: EdgeProperties{Weight: 2}},
				},
				"B": {
					"A": {Source: "B", Target: "A", Properties: EdgeProperties{Weight: 2}},
				},
				"C": {
					"D": {Source: "C", Target: "D", Properties: EdgeProperties{Weight: 4}},
				},
				"D": {
					"C": {Source: "D", Target: "C", Properties: EdgeProperties{Weight: 4}},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			g := New(StringHash)

			for _, vertex := range test.vertices {
				_ = g.AddVertex(vertex)
			}

			for _, edge := range test.edges {
				_ = g.AddEdge(copyEdge(edge))
			}

			mst, _ := MaximumSpanningTree(g)
			adjacencyMap, _ := mst.AdjacencyMap()

			edgesAreEqual := g.(*undirected[string, string]).edgesAreEqual

			if !adjacencyMapsAreEqual(test.expectedMSTAdjacencyMap, adjacencyMap, edgesAreEqual) {
				t.Fatalf("expected adjacency map %v, got %v", test.expectedMSTAdjacencyMap, adjacencyMap)
			}
		})
	}
}

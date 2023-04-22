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
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			g := New(StringHash)

			for _, vertex := range test.vertices {
				_ = g.AddVertex(vertex)
			}

			for _, edge := range test.edges {
				_ = g.AddEdge(edge.Source, edge.Target, copyEdgeProperties(edge.Properties))
			}

			mst, _ := MinimumSpanningTree(g)
			adjacencyMap, _ := mst.AdjacencyMap()

			for expectedHash, expectedAdjacencies := range test.expectedMSTAdjacencyMap {
				actualAdjacencies, ok := adjacencyMap[expectedHash]
				if !ok {
					t.Errorf("%s: key %v doesn't exist in adjacency map", name, expectedHash)
					continue
				}

				for expectedAdjacency, expectedEdge := range expectedAdjacencies {
					actualEdge, ok := actualAdjacencies[expectedAdjacency]
					if !ok {
						t.Errorf("%s: key %v doesn't exist in adjacencies of %v", name, expectedAdjacency, expectedHash)
						continue
					}

					if !mst.(*undirected[string, string]).edgesAreEqual(expectedEdge, actualEdge) {
						t.Errorf("%s: expected edge %v, got %v at AdjacencyMap[%v][%v]", name, expectedEdge, actualEdge, expectedHash, expectedAdjacency)
					}

					for expectedKey, expectedValue := range expectedEdge.Properties.Attributes {
						actualValue, ok := actualEdge.Properties.Attributes[expectedKey]
						if !ok {
							t.Errorf("%s: expected attribute %v to exist in edge %v", name, expectedKey, actualEdge)
						}
						if actualValue != expectedValue {
							t.Errorf("%s: expected value %v for key %v in edge %v, got %v", name, expectedValue, expectedKey, expectedEdge, actualValue)
						}
					}

					if actualEdge.Properties.Weight != expectedEdge.Properties.Weight {
						t.Errorf("%s: expected weight %v for edge %v, got %v", name, expectedEdge.Properties.Weight, expectedEdge, actualEdge.Properties.Weight)
					}
				}
			}

			for actualHash := range adjacencyMap {
				if _, ok := test.expectedMSTAdjacencyMap[actualHash]; !ok {
					t.Errorf("%s: unexpected key %v in union adjacency map", name, actualHash)
				}
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
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			g := New(StringHash)

			for _, vertex := range test.vertices {
				_ = g.AddVertex(vertex)
			}

			for _, edge := range test.edges {
				_ = g.AddEdge(edge.Source, edge.Target, copyEdgeProperties(edge.Properties))
			}

			mst, _ := MaximumSpanningTree(g)
			adjacencyMap, _ := mst.AdjacencyMap()

			for expectedHash, expectedAdjacencies := range test.expectedMSTAdjacencyMap {
				actualAdjacencies, ok := adjacencyMap[expectedHash]
				if !ok {
					t.Errorf("%s: key %v doesn't exist in adjacency map", name, expectedHash)
					continue
				}

				for expectedAdjacency, expectedEdge := range expectedAdjacencies {
					actualEdge, ok := actualAdjacencies[expectedAdjacency]
					if !ok {
						t.Errorf("%s: key %v doesn't exist in adjacencies of %v", name, expectedAdjacency, expectedHash)
						continue
					}

					if !mst.(*undirected[string, string]).edgesAreEqual(expectedEdge, actualEdge) {
						t.Errorf("%s: expected edge %v, got %v at AdjacencyMap[%v][%v]", name, expectedEdge, actualEdge, expectedHash, expectedAdjacency)
					}

					for expectedKey, expectedValue := range expectedEdge.Properties.Attributes {
						actualValue, ok := actualEdge.Properties.Attributes[expectedKey]
						if !ok {
							t.Errorf("%s: expected attribute %v to exist in edge %v", name, expectedKey, actualEdge)
						}
						if actualValue != expectedValue {
							t.Errorf("%s: expected value %v for key %v in edge %v, got %v", name, expectedValue, expectedKey, expectedEdge, actualValue)
						}
					}

					if actualEdge.Properties.Weight != expectedEdge.Properties.Weight {
						t.Errorf("%s: expected weight %v for edge %v, got %v", name, expectedEdge.Properties.Weight, expectedEdge, actualEdge.Properties.Weight)
					}
				}
			}

			for actualHash := range adjacencyMap {
				if _, ok := test.expectedMSTAdjacencyMap[actualHash]; !ok {
					t.Errorf("%s: unexpected key %v in union adjacency map", name, actualHash)
				}
			}
		})
	}
}

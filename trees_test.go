package graph

import (
	"testing"
)

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

func TestUnionFind_add(t *testing.T) {
	tests := map[string]struct {
		vertex         int
		expectedParent int
	}{
		"add vertex 1": {
			vertex:         1,
			expectedParent: 1,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			u := newUnionFind[int]()
			u.add(test.vertex)

			if u.parents[test.vertex] != test.expectedParent {
				t.Errorf("expected parent %v, got %v", test.expectedParent, u.parents[test.vertex])
			}
		})
	}
}

func TestUnionFind_union(t *testing.T) {
	tests := map[string]struct {
		vertices        []int
		args            [2]int
		expectedParents map[int]int
	}{
		"merge 1 and 2": {
			vertices: []int{1, 2, 3, 4},
			args:     [2]int{1, 2},
			expectedParents: map[int]int{
				1: 1,
				2: 1,
				3: 3,
				4: 4,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			u := newUnionFind(test.vertices...)
			u.union(test.args[0], test.args[1])

			for expectedKey, expectedValue := range test.expectedParents {
				actualValue, ok := u.parents[expectedKey]
				if !ok {
					t.Fatalf("expected key %v not found", expectedKey)
				}
				if actualValue != expectedValue {
					t.Fatalf("expected value %v for %v, got %v", expectedValue, expectedKey, actualValue)
				}
			}
		})
	}
}

func TestUnionFind_find(t *testing.T) {
	tests := map[string]struct {
		parents  map[int]int
		arg      int
		expected int
	}{
		"find 1 in sets 1, 2, 3": {
			parents: map[int]int{
				1: 1,
				2: 2,
				3: 3,
			},
			arg:      1,
			expected: 1,
		},
		"find 3 in sets 1-2, 3-4": {
			parents: map[int]int{
				1: 1,
				2: 1,
				3: 3,
				4: 3,
			},
			arg:      3,
			expected: 3,
		},
		"find 4 in sets 1-2, 3-4": {
			parents: map[int]int{
				1: 1,
				2: 1,
				3: 3,
				4: 3,
			},
			arg:      4,
			expected: 3,
		},
		"find 3 in set 1-2-3": {
			parents: map[int]int{
				1: 1,
				2: 1,
				3: 2,
			},
			arg:      3,
			expected: 1,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			u := unionFind[int]{
				parents: test.parents,
			}

			actual := u.find(test.arg)

			if actual != test.expected {
				t.Errorf("expected %v, got %v", test.expected, actual)
			}
		})
	}
}

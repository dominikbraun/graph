package graph

import (
	"testing"
)

func TestDirectedUnion(t *testing.T) {
	tests := map[string]struct {
		gVertices            []int
		gVertexProperties    map[int]VertexProperties
		gEdges               []Edge[int]
		hVertices            []int
		hVertexProperties    map[int]VertexProperties
		hEdges               []Edge[int]
		expectedAdjacencyMap map[int]map[int]Edge[int]
	}{
		"two 3-vertices directed graphs": {
			gVertices:         []int{1, 2, 3},
			gVertexProperties: map[int]VertexProperties{},
			gEdges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 2, Target: 3},
			},
			hVertices:         []int{4, 5, 6},
			hVertexProperties: map[int]VertexProperties{},
			hEdges: []Edge[int]{
				{Source: 4, Target: 5},
				{Source: 5, Target: 6},
			},
			expectedAdjacencyMap: map[int]map[int]Edge[int]{
				1: {
					2: {Source: 1, Target: 2},
				},
				2: {
					3: {Source: 2, Target: 3},
				},
				3: {},
				4: {
					5: {Source: 4, Target: 5},
				},
				5: {
					6: {Source: 5, Target: 6},
				},
				6: {},
			},
		},
		"vertices and edges with properties": {
			gVertices: []int{1, 2},
			gVertexProperties: map[int]VertexProperties{
				1: {
					Attributes: map[string]string{
						"color": "red",
					},
					Weight: 10,
				},
				2: {
					Attributes: map[string]string{},
					Weight:     20,
				},
			},
			gEdges: []Edge[int]{
				{
					Source: 1,
					Target: 2,
					Properties: EdgeProperties{
						Attributes: map[string]string{
							"label": "my-edge",
						},
						Weight: 42,
						Data:   "edge data #1",
					},
				},
			},
			hVertices: []int{3, 4},
			hVertexProperties: map[int]VertexProperties{
				3: {
					Attributes: map[string]string{
						"color": "blue",
					},
					Weight: 15,
				},
			},
			hEdges: []Edge[int]{
				{
					Source: 3,
					Target: 4,
					Properties: EdgeProperties{
						Attributes: map[string]string{
							"label": "another-edge",
						},
						Weight: 50,
						Data:   "edge data #2",
					},
				},
			},
			expectedAdjacencyMap: map[int]map[int]Edge[int]{
				1: {
					2: {
						Source: 1,
						Target: 2,
						Properties: EdgeProperties{
							Attributes: map[string]string{
								"label": "my-edge",
							},
							Weight: 42,
							Data:   "edge data #1",
						},
					},
				},
				2: {},
				3: {
					4: {
						Source: 3,
						Target: 4,
						Properties: EdgeProperties{
							Attributes: map[string]string{
								"label": "another-edge",
							},
							Weight: 50,
							Data:   "edge data #2",
						},
					},
				},
				4: {},
			},
		},
	}

	for name, test := range tests {
		g := New(IntHash, Directed())

		for _, vertex := range test.gVertices {
			_ = g.AddVertex(vertex, copyVertexProperties(test.gVertexProperties[vertex]))
		}

		for _, edge := range test.gEdges {
			_ = g.AddEdge(edge.Source, edge.Target, copyEdgeProperties(edge.Properties))
		}

		h := New(IntHash, Directed())

		for _, vertex := range test.hVertices {
			_ = h.AddVertex(vertex, copyVertexProperties(test.gVertexProperties[vertex]))
		}

		for _, edge := range test.hEdges {
			_ = h.AddEdge(edge.Source, edge.Target, copyEdgeProperties(edge.Properties))
		}

		union, err := Union(g, h)
		if err != nil {
			t.Fatalf("%s: unexpected union error: %s", name, err.Error())
		}

		unionAdjacencyMap, err := union.AdjacencyMap()
		if err != nil {
			t.Fatalf("%s: unexpected adjaceny map error: %s", name, err.Error())
		}

		edgesAreEqual := g.(*directed[int, int]).edgesAreEqual

		if !adjacencyMapsAreEqual(test.expectedAdjacencyMap, unionAdjacencyMap, edgesAreEqual) {
			t.Fatalf("expected adjacency map %v, got %v", test.expectedAdjacencyMap, unionAdjacencyMap)
		}
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

func adjacencyMapsAreEqual[K comparable](a, b map[K]map[K]Edge[K], edgesAreEqual func(a, b Edge[K]) bool) bool {
	for aHash, aAdjacencies := range a {
		bAdjacencies, ok := b[aHash]
		if !ok {
			return false
		}

		for aAdjacency, aEdge := range aAdjacencies {
			bEdge, ok := bAdjacencies[aAdjacency]
			if !ok {
				return false
			}

			if !edgesAreEqual(aEdge, bEdge) {
				return false
			}

			for aKey, aValue := range aEdge.Properties.Attributes {
				bValue, ok := bEdge.Properties.Attributes[aKey]
				if !ok {
					return false
				}
				if bValue != aValue {
					return false
				}
			}

			if bEdge.Properties.Weight != aEdge.Properties.Weight {
				return false
			}
		}
	}

	for aHash := range a {
		if _, ok := b[aHash]; !ok {
			return false
		}
	}

	return true
}

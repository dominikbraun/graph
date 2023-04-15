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

		for expectedHash, expectedAdjacencies := range test.expectedAdjacencyMap {
			actualAdjacencies, ok := unionAdjacencyMap[expectedHash]
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

				if !union.(*directed[int, int]).edgesAreEqual(expectedEdge, actualEdge) {
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

		for actualHash, _ := range unionAdjacencyMap {
			if _, ok := test.expectedAdjacencyMap[actualHash]; !ok {
				t.Errorf("%s: unexpected key %v in union adjacency map", name, actualHash)
			}
		}
	}
}

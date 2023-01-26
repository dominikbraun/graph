package graph

import (
	"errors"
	"testing"
)

func TestDirected_Traits(t *testing.T) {
	tests := map[string]struct {
		traits   *Traits
		expected *Traits
	}{
		"default traits": {
			traits:   &Traits{},
			expected: &Traits{},
		},
		"directed": {
			traits:   &Traits{IsDirected: true},
			expected: &Traits{IsDirected: true},
		},
	}

	for name, test := range tests {
		g := newDirected(IntHash, test.traits, newMemoryStore[int, int]())
		traits := g.Traits()

		if !traitsAreEqual(traits, test.expected) {
			t.Errorf("%s: traits expectancy doesn't match: expected %v, got %v", name, test.expected, traits)
		}
	}
}

func TestDirected_AddVertex(t *testing.T) {
	tests := map[string]struct {
		vertices           []int
		properties         *VertexProperties
		expectedVertices   []int
		expectedProperties *VertexProperties
		// Even though some AddVertex calls might work, at least one of them could fail, for example
		// if the last call would add an existing vertex.
		finallyExpectedError error
	}{
		"graph with 3 vertices": {
			vertices: []int{1, 2, 3},
			properties: &VertexProperties{
				Attributes: map[string]string{"color": "red"},
				Weight:     10,
			},
			expectedVertices: []int{1, 2, 3},
			expectedProperties: &VertexProperties{
				Attributes: map[string]string{"color": "red"},
				Weight:     10,
			},
		},
		"graph with duplicated vertex": {
			vertices:             []int{1, 2, 2},
			expectedVertices:     []int{1, 2},
			finallyExpectedError: ErrVertexAlreadyExists,
		},
	}

	for name, test := range tests {
		graph := newDirected(IntHash, &Traits{}, newMemoryStore[int, int]())

		var err error

		for _, vertex := range test.vertices {
			if test.properties == nil {
				err = graph.AddVertex(vertex)
				continue
			}
			// If there are vertex attributes, iterate over them and call VertexAttribute for each
			// entry. A vertex should only have one attribute so that AddVertex is invoked once.
			for key, value := range test.properties.Attributes {
				err = graph.AddVertex(vertex, VertexWeight(test.properties.Weight), VertexAttribute(key, value))
			}
		}

		if !errors.Is(err, test.finallyExpectedError) {
			t.Errorf("%s: error expectancy doesn't match: expected %v, got %v", name, test.finallyExpectedError, err)
		}

		graphStore := graph.store.(*memoryStore[int, int])
		for _, vertex := range test.vertices {
			if len(graphStore.vertices) != len(test.expectedVertices) {
				t.Errorf("%s: vertex count doesn't match: expected %v, got %v", name, len(test.expectedVertices), len(graphStore.vertices))
			}

			hash := graph.hash(vertex)
			if _, _, err := graph.store.Vertex(hash); err != nil {
				vertices := graphStore.vertices
				t.Errorf("%s: vertex %v not found in graph: %v", name, vertex, vertices)
			}

			if test.properties == nil {
				continue
			}

			if graphStore.vertexProperties[hash].Weight != test.expectedProperties.Weight {
				t.Errorf("%s: edge weights don't match: expected weight %v, got %v", name, test.expectedProperties.Weight, graphStore.vertexProperties[hash].Weight)
			}

			if len(graphStore.vertexProperties[hash].Attributes) != len(test.expectedProperties.Attributes) {
				t.Fatalf("%s: attributes lengths don't match: expcted %v, got %v", name, len(test.expectedProperties.Attributes), len(graphStore.vertexProperties[hash].Attributes))
			}

			for expectedKey, expectedValue := range test.expectedProperties.Attributes {
				value, ok := graphStore.vertexProperties[hash].Attributes[expectedKey]
				if !ok {
					t.Errorf("%s: attribute keys don't match: expected key %v not found", name, expectedKey)
				}
				if value != expectedValue {
					t.Errorf("%s: attribute values don't match: expected value %v for key %v, got %v", name, expectedValue, expectedKey, value)
				}
			}
		}
	}
}

func TestDirected_Vertex(t *testing.T) {
	tests := map[string]struct {
		vertices      []int
		vertex        int
		expectedError error
	}{
		"existing vertex": {
			vertices: []int{1, 2, 3},
			vertex:   2,
		},
		"non-existent vertex": {
			vertices:      []int{1, 2, 3},
			vertex:        4,
			expectedError: ErrVertexNotFound,
		},
	}

	for name, test := range tests {
		graph := newDirected(IntHash, &Traits{}, newMemoryStore[int, int]())

		for _, vertex := range test.vertices {
			_ = graph.AddVertex(vertex)
		}

		vertex, err := graph.Vertex(test.vertex)

		if !errors.Is(err, test.expectedError) {
			t.Errorf("%s: error expectancy doesn't match: expected %v, got %v", name, test.expectedError, err)
		}

		if test.expectedError != nil {
			continue
		}

		if vertex != test.vertex {
			t.Errorf("%s: vertex expectancy doesn't match: expected %v, got %v", name, test.vertex, vertex)
		}
	}
}

func TestDirected_AddEdge(t *testing.T) {
	tests := map[string]struct {
		vertices      []int
		edges         []Edge[int]
		traits        *Traits
		expectedEdges []Edge[int]
		// Even though some AddEdge calls might work, at least one of them could fail, for example
		// if the last call would introduce a cycle.
		finallyExpectedError error
	}{
		"graph with 2 edges": {
			vertices: []int{1, 2, 3},
			edges: []Edge[int]{
				{Source: 1, Target: 2, Properties: EdgeProperties{Weight: 10}},
				{Source: 1, Target: 3, Properties: EdgeProperties{Weight: 20}},
			},
			traits: &Traits{},
			expectedEdges: []Edge[int]{
				{Source: 1, Target: 2, Properties: EdgeProperties{Weight: 10}},
				{Source: 1, Target: 3, Properties: EdgeProperties{Weight: 20}},
			},
		},
		"hashes for non-existent vertices": {
			vertices: []int{1, 2},
			edges: []Edge[int]{
				{Source: 1, Target: 3, Properties: EdgeProperties{Weight: 20}},
			},
			traits:               &Traits{},
			finallyExpectedError: ErrVertexNotFound,
		},
		"edge introducing a cycle in an acyclic graph": {
			vertices: []int{1, 2, 3},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 2, Target: 3},
				{Source: 3, Target: 1},
			},
			traits: &Traits{
				PreventCycles: true,
			},
			finallyExpectedError: ErrEdgeCreatesCycle,
		},
		"edge already exists": {
			vertices: []int{1, 2, 3},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 2, Target: 3},
				{Source: 3, Target: 1},
				{Source: 3, Target: 1},
			},
			traits:               &Traits{},
			finallyExpectedError: ErrEdgeAlreadyExists,
		},
		"edge with attributes": {
			vertices: []int{1, 2},
			edges: []Edge[int]{
				{
					Source: 1,
					Target: 2,
					Properties: EdgeProperties{
						Attributes: map[string]string{
							"color": "red",
						},
					},
				},
			},
			expectedEdges: []Edge[int]{
				{
					Source: 1,
					Target: 2,
					Properties: EdgeProperties{
						Attributes: map[string]string{
							"color": "red",
						},
					},
				},
			},
			traits: &Traits{},
		},
		"edge with data": {
			vertices: []int{1, 2},
			edges: []Edge[int]{
				{
					Source: 1,
					Target: 2,
					Properties: EdgeProperties{
						Data: "foo",
					},
				},
			},
			expectedEdges: []Edge[int]{
				{
					Source: 1,
					Target: 2,
					Properties: EdgeProperties{
						Data: "foo",
					},
				},
			},
			traits: &Traits{},
		},
	}

	for name, test := range tests {
		graph := newDirected(IntHash, test.traits, newMemoryStore[int, int]())

		for _, vertex := range test.vertices {
			_ = graph.AddVertex(vertex)
		}

		var err error

		for _, edge := range test.edges {
			if len(edge.Properties.Attributes) == 0 {
				err = graph.AddEdge(edge.Source, edge.Target, EdgeWeight(edge.Properties.Weight), EdgeData(edge.Properties.Data))
			}
			// If there are edge attributes, iterate over them and call EdgeAttribute for each
			// entry. An edge should only have one attribute so that AddEdge is invoked once.
			for key, value := range edge.Properties.Attributes {
				err = graph.AddEdge(
					edge.Source,
					edge.Target,
					EdgeWeight(edge.Properties.Weight),
					EdgeData(edge.Properties.Data),
					EdgeAttribute(key, value),
				)
			}
			if err != nil {
				break
			}
		}

		if !errors.Is(err, test.finallyExpectedError) {
			t.Fatalf("%s: error expectancy doesn't match: expected %v, got %v", name, test.finallyExpectedError, err)
		}

		for _, expectedEdge := range test.expectedEdges {
			sourceHash := graph.hash(expectedEdge.Source)
			targetHash := graph.hash(expectedEdge.Target)

			edge, err := graph.Edge(sourceHash, targetHash)
			if err != nil {
				t.Fatalf("%s: edge with source %v and target %v not found", name, expectedEdge.Source, expectedEdge.Target)
			}

			if edge.Source != expectedEdge.Source {
				t.Errorf("%s: edge sources don't match: expected source %v, got %v", name, expectedEdge.Source, edge.Source)
			}

			if edge.Target != expectedEdge.Target {
				t.Errorf("%s: edge targets don't match: expected target %v, got %v", name, expectedEdge.Target, edge.Target)
			}

			if edge.Properties.Weight != expectedEdge.Properties.Weight {
				t.Errorf("%s: edge weights don't match: expected weight %v, got %v", name, expectedEdge.Properties.Weight, edge.Properties.Weight)
			}

			if len(edge.Properties.Attributes) != len(expectedEdge.Properties.Attributes) {
				t.Fatalf("%s: attributes length don't match: expcted %v, got %v", name, len(expectedEdge.Properties.Attributes), len(edge.Properties.Attributes))
			}

			for expectedKey, expectedValue := range expectedEdge.Properties.Attributes {
				value, ok := edge.Properties.Attributes[expectedKey]
				if !ok {
					t.Errorf("%s: attribute keys don't match: expected key %v not found", name, expectedKey)
				}
				if value != expectedValue {
					t.Errorf("%s: attribute values don't match: expected value %v for key %v, got %v", name, expectedValue, expectedKey, value)
				}
			}

			if edge.Properties.Data != expectedEdge.Properties.Data {
				t.Errorf("%s: edge data doesn't match: expected data %v, got %v", name, expectedEdge.Properties.Data, edge.Properties.Data)
			}
		}
	}
}

func TestDirected_Edge(t *testing.T) {
	tests := map[string]struct {
		vertices      []int
		getEdgeHashes [2]int
		exists        bool
	}{
		"get edge of directed graph": {
			vertices:      []int{1, 2, 3},
			getEdgeHashes: [2]int{1, 2},
			exists:        true,
		},
		"get non-existent edge of directed graph": {
			vertices:      []int{1, 2, 3},
			getEdgeHashes: [2]int{1, 3},
		},
	}

	for name, test := range tests {
		graph := newDirected(IntHash, &Traits{}, newMemoryStore[int, int]())
		for _, vertex := range test.vertices {
			_ = graph.AddVertex(vertex)
		}

		sourceHash := graph.hash(test.vertices[0])
		targetHash := graph.hash(test.vertices[1])

		err := graph.AddEdge(sourceHash, targetHash)
		if err != nil {
			t.Fatalf("%s: error adding edge: %v", name, err)
		}

		_, err = graph.Edge(test.getEdgeHashes[0], test.getEdgeHashes[1])

		if test.exists != (err == nil) {
			t.Fatalf("%s: result expectancy doesn't match: expected %v, got %v", name, test.exists, err)
		}
	}
}

func TestDirected_RemoveEdge(t *testing.T) {
	tests := map[string]struct {
		vertices      []int
		edges         []Edge[int]
		removeEdges   []Edge[int]
		expectedError error
	}{
		"two-vertices graph": {
			vertices: []int{1, 2},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
			},
			removeEdges: []Edge[int]{
				{Source: 1, Target: 2},
			},
		},
		"remove 2 edges from triangle": {
			vertices: []int{1, 2, 3},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
				{Source: 2, Target: 3},
			},
			removeEdges: []Edge[int]{
				{Source: 1, Target: 3},
				{Source: 2, Target: 3},
			},
		},
		"remove non-existent edge": {
			vertices: []int{1, 2, 3},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
			},
			removeEdges: []Edge[int]{
				{Source: 2, Target: 3},
			},
			expectedError: ErrEdgeNotFound,
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

		for _, removeEdge := range test.removeEdges {
			if err := graph.RemoveEdge(removeEdge.Source, removeEdge.Target); !errors.Is(err, test.expectedError) {
				t.Errorf("%s: error expectancy doesn't match: expected %v, got %v", name, test.expectedError, err)
			}
			// After removing the edge, verify that it can't be retrieved using Edge anymore.
			if _, err := graph.Edge(removeEdge.Source, removeEdge.Target); err != ErrEdgeNotFound {
				t.Fatalf("%s: error expectancy doesn't match: expected %v, got %v", name, ErrEdgeNotFound, err)
			}
		}
	}
}

func TestDirected_AdjacencyList(t *testing.T) {
	tests := map[string]struct {
		vertices []int
		edges    []Edge[int]
		expected map[int]map[int]Edge[int]
	}{
		"Y-shaped graph": {
			vertices: []int{1, 2, 3, 4},
			edges: []Edge[int]{
				{Source: 1, Target: 3},
				{Source: 2, Target: 3},
				{Source: 3, Target: 4},
			},
			expected: map[int]map[int]Edge[int]{
				1: {
					3: {Source: 1, Target: 3},
				},
				2: {
					3: {Source: 2, Target: 3},
				},
				3: {
					4: {Source: 3, Target: 4},
				},
				4: {},
			},
		},
		"diamond-shaped graph": {
			vertices: []int{1, 2, 3, 4},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
				{Source: 2, Target: 4},
				{Source: 3, Target: 4},
			},
			expected: map[int]map[int]Edge[int]{
				1: {
					2: {Source: 1, Target: 2},
					3: {Source: 1, Target: 3},
				},
				2: {
					4: {Source: 2, Target: 4},
				},
				3: {
					4: {Source: 3, Target: 4},
				},
				4: {},
			},
		},
	}

	for name, test := range tests {
		graph := newDirected(IntHash, &Traits{}, newMemoryStore[int, int]())

		for _, vertex := range test.vertices {
			_ = graph.AddVertex(vertex)
		}

		for _, edge := range test.edges {
			if err := graph.AddEdge(edge.Source, edge.Target, EdgeWeight(edge.Properties.Weight)); err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		adjacencyMap, _ := graph.AdjacencyMap()

		for expectedVertex, expectedAdjacencies := range test.expected {
			adjacencies, ok := adjacencyMap[expectedVertex]
			if !ok {
				t.Errorf("%s: expected vertex %v does not exist in adjacency map", name, expectedVertex)
			}

			for expectedAdjacency, expectedEdge := range expectedAdjacencies {
				edge, ok := adjacencies[expectedAdjacency]
				if !ok {
					t.Errorf("%s: expected adjacency %v does not exist in map of %v", name, expectedAdjacency, expectedVertex)
				}
				if edge.Source != expectedEdge.Source || edge.Target != expectedEdge.Target {
					t.Errorf("%s: edge expectancy doesn't match: expected %v, got %v", name, expectedEdge, edge)
				}
			}
		}
	}
}

func TestDirected_PredecessorMap(t *testing.T) {
	tests := map[string]struct {
		vertices []int
		edges    []Edge[int]
		expected map[int]map[int]Edge[int]
	}{
		"Y-shaped graph": {
			vertices: []int{1, 2, 3, 4},
			edges: []Edge[int]{
				{Source: 1, Target: 3},
				{Source: 2, Target: 3},
				{Source: 3, Target: 4},
			},
			expected: map[int]map[int]Edge[int]{
				1: {},
				2: {},
				3: {
					1: {Source: 1, Target: 3},
					2: {Source: 2, Target: 3},
				},
				4: {
					3: {Source: 3, Target: 4},
				},
			},
		},
		"diamond-shaped graph": {
			vertices: []int{1, 2, 3, 4},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
				{Source: 2, Target: 4},
				{Source: 3, Target: 4},
			},
			expected: map[int]map[int]Edge[int]{
				1: {},
				2: {
					1: {Source: 1, Target: 2},
				},
				3: {
					1: {Source: 1, Target: 3},
				},
				4: {
					2: {Source: 2, Target: 4},
					3: {Source: 3, Target: 4},
				},
			},
		},
	}

	for name, test := range tests {
		graph := newDirected(IntHash, &Traits{}, newMemoryStore[int, int]())

		for _, vertex := range test.vertices {
			_ = graph.AddVertex(vertex)
		}

		for _, edge := range test.edges {
			if err := graph.AddEdge(edge.Source, edge.Target, EdgeWeight(edge.Properties.Weight)); err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		predecessors, _ := graph.PredecessorMap()

		for expectedVertex, expectedPredecessors := range test.expected {
			predecessors, ok := predecessors[expectedVertex]
			if !ok {
				t.Errorf("%s: expected vertex %v does not exist in adjacency map", name, expectedVertex)
			}

			for expectedPredecessor, expectedEdge := range expectedPredecessors {
				edge, ok := predecessors[expectedPredecessor]
				if !ok {
					t.Errorf("%s: expected adjacency %v does not exist in map of %v", name, expectedPredecessor, expectedVertex)
				}
				if edge.Source != expectedEdge.Source || edge.Target != expectedEdge.Target {
					t.Errorf("%s: edge expectancy doesn't match: expected %v, got %v", name, expectedEdge, edge)
				}
			}
		}
	}
}

func TestDirected_Clone(t *testing.T) {
	tests := map[string]struct {
		vertices []int
		edges    []Edge[int]
	}{
		"Y-shaped graph": {
			vertices: []int{1, 2, 3, 4},
			edges: []Edge[int]{
				{Source: 1, Target: 3},
				{Source: 2, Target: 3},
				{Source: 3, Target: 4},
			},
		},
		"diamond-shaped graph": {
			vertices: []int{1, 2, 3, 4},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
				{Source: 2, Target: 4},
				{Source: 3, Target: 4},
			},
		},
	}

	for name, test := range tests {
		graph := New(IntHash, Directed())

		for _, vertex := range test.vertices {
			_ = graph.AddVertex(vertex, VertexWeight(vertex), VertexAttribute("color", "red"))
		}

		for _, edge := range test.edges {
			if err := graph.AddEdge(edge.Source, edge.Target, EdgeWeight(edge.Properties.Weight)); err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		clonedGraph, err := graph.Clone()
		if err != nil {
			t.Fatalf("%s: failed to clone graph: %s", name, err.Error())
		}

		expected := graph.(*directed[int, int])
		actual := clonedGraph.(*directed[int, int])

		if actual.hash(5) != expected.hash(5) {
			t.Errorf("%s: hash expectancy doesn't match: expected %v, got %v", name, expected.hash, actual.hash)
		}

		if !traitsAreEqual(actual.traits, expected.traits) {
			t.Errorf("%s: traits expectancy doesn't match: expected %v, got %v", name, expected.traits, actual.traits)
		}

		if actual.store != expected.store {
			t.Fatalf("%s: stores don't match", name)
		}
	}
}

func TestDirected_OrderAndSize(t *testing.T) {
	tests := map[string]struct {
		vertices      []int
		edges         []Edge[int]
		expectedOrder int
		expectedSize  int
	}{
		"Y-shaped graph": {
			vertices: []int{1, 2, 3, 4},
			edges: []Edge[int]{
				{Source: 1, Target: 3},
				{Source: 2, Target: 3},
				{Source: 3, Target: 4},
			},
			expectedOrder: 4,
			expectedSize:  3,
		},
		"diamond-shaped graph": {
			vertices: []int{1, 2, 3, 4},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
				{Source: 2, Target: 4},
				{Source: 3, Target: 4},
			},
			expectedOrder: 4,
			expectedSize:  4,
		},
		"two-vertices two-edges graph": {
			vertices: []int{1, 2},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 2, Target: 1},
			},
			expectedOrder: 2,
			expectedSize:  2,
		},
		"edgeless graph": {
			vertices:      []int{1, 2},
			edges:         []Edge[int]{},
			expectedOrder: 2,
			expectedSize:  0,
		},
	}

	for name, test := range tests {
		graph := newDirected(IntHash, &Traits{}, newMemoryStore[int, int]())

		for _, vertex := range test.vertices {
			_ = graph.AddVertex(vertex)
		}

		for _, edge := range test.edges {
			if err := graph.AddEdge(edge.Source, edge.Target, EdgeWeight(edge.Properties.Weight)); err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		order, _ := graph.Order()
		size, _ := graph.Size()

		if order != test.expectedOrder {
			t.Errorf("%s: order expectancy doesn't match: expected %d, got %d", name, test.expectedOrder, order)
		}

		if size != test.expectedSize {
			t.Errorf("%s: size expectancy doesn't match: expected %d, got %d", name, test.expectedSize, size)
		}

	}
}

func TestDirected_edgesAreEqual(t *testing.T) {
	tests := map[string]struct {
		a             Edge[int]
		b             Edge[int]
		edgesAreEqual bool
	}{
		"equal edges in directed graph": {
			a:             Edge[int]{Source: 1, Target: 2},
			b:             Edge[int]{Source: 1, Target: 2},
			edgesAreEqual: true,
		},
		"swapped equal edges in directed graph": {
			a: Edge[int]{Source: 1, Target: 2},
			b: Edge[int]{Source: 2, Target: 1},
		},
	}

	for name, test := range tests {
		graph := newDirected(IntHash, &Traits{}, newMemoryStore[int, int]())
		actual := graph.edgesAreEqual(test.a, test.b)

		if actual != test.edgesAreEqual {
			t.Errorf("%s: equality expectations don't match: expected %v, got %v", name, test.edgesAreEqual, actual)
		}
	}
}

func TestDirected_addEdge(t *testing.T) {
	tests := map[string]struct {
		edges []Edge[int]
	}{
		"add 3 edges": {
			edges: []Edge[int]{
				{Source: 1, Target: 2, Properties: EdgeProperties{Weight: 1}},
				{Source: 2, Target: 3, Properties: EdgeProperties{Weight: 2}},
				{Source: 3, Target: 1, Properties: EdgeProperties{Weight: 3}},
			},
		},
	}

	for name, test := range tests {
		graph := newDirected(IntHash, &Traits{}, newMemoryStore[int, int]())

		for _, edge := range test.edges {
			sourceHash := graph.hash(edge.Source)
			TargetHash := graph.hash(edge.Target)
			err := graph.addEdge(sourceHash, TargetHash, edge)
			if err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		outEdges := graph.store.(*memoryStore[int, int]).outEdges
		if len(outEdges) != len(test.edges) {
			t.Errorf("%s: number of outgoing edges doesn't match: expected %v, got %v", name, len(test.edges), len(outEdges))
		}
		inEdges := graph.store.(*memoryStore[int, int]).inEdges
		if len(inEdges) != len(test.edges) {
			t.Errorf("%s: number of ingoing edges doesn't match: expected %v, got %v", name, len(test.edges), len(inEdges))
		}
	}
}

func TestDirected_predecessors(t *testing.T) {
	tests := map[string]struct {
		vertices             []int
		edges                []Edge[int]
		vertex               int
		expectedPredecessors []int
	}{
		"graph with 3 vertices": {
			vertices: []int{1, 2, 3},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
			},
			vertex:               2,
			expectedPredecessors: []int{1},
		},
		"graph with 6 vertices": {
			vertices: []int{1, 2, 3, 4, 5, 6},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
				{Source: 2, Target: 4},
				{Source: 2, Target: 5},
				{Source: 3, Target: 6},
			},
			vertex:               5,
			expectedPredecessors: []int{2},
		},
		"graph with 4 vertices and 3 predecessors": {
			vertices: []int{1, 2, 3, 4},
			edges: []Edge[int]{
				{Source: 1, Target: 4},
				{Source: 2, Target: 4},
				{Source: 3, Target: 4},
			},
			vertex:               4,
			expectedPredecessors: []int{1, 2, 3},
		},
	}

	for name, test := range tests {
		graph := newDirected(IntHash, &Traits{}, newMemoryStore[int, int]())

		for _, vertex := range test.vertices {
			_ = graph.AddVertex(vertex)
		}

		for _, edge := range test.edges {
			if err := graph.AddEdge(edge.Source, edge.Target); err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		predecessors, _ := predecessors(graph, graph.hash(test.vertex))

		if !slicesAreEqual(predecessors, test.expectedPredecessors) {
			t.Errorf("%s: predecessors don't match: expected %v, got %v", name, test.expectedPredecessors, predecessors)
		}
	}
}

func slicesAreEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	for _, aValue := range a {
		found := false
		for _, bValue := range b {
			if aValue == bValue {
				found = true
			}
		}
		if !found {
			return false
		}
	}

	return true
}

func predecessors[K comparable, T any](g *directed[K, T], vertexHash K) ([]K, error) {
	var predecessorHashes []K

	predecessorMap, err := g.PredecessorMap()
	if err != nil {
		return nil, err
	}

	for _, edge := range predecessorMap[vertexHash] {
		predecessorHashes = append(predecessorHashes, edge.Source)
	}

	return predecessorHashes, nil
}

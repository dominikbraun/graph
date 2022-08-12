package graph

import (
	"testing"
)

func TestUndirected_Traits(t *testing.T) {
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
		g := newUndirected(IntHash, test.traits)
		traits := g.Traits()

		if !traitsAreEqual(traits, test.expected) {
			t.Errorf("%s: traits expectancy doesn't match: expected %v, got %v", name, test.expected, traits)
		}
	}
}

func TestUndirected_AddVertex(t *testing.T) {
	tests := map[string]struct {
		vertices         []int
		expectedVertices []int
	}{
		"graph with 3 vertices": {
			vertices:         []int{1, 2, 3},
			expectedVertices: []int{1, 2, 3},
		},
		"graph with duplicated vertex": {
			vertices:         []int{1, 2, 2},
			expectedVertices: []int{1, 2},
		},
	}

	for name, test := range tests {
		graph := newUndirected(IntHash, &Traits{})

		for _, vertex := range test.vertices {
			graph.AddVertex(vertex)
		}

		for _, vertex := range test.vertices {
			hash := graph.hash(vertex)
			if _, ok := graph.vertices[hash]; !ok {
				t.Errorf("%s: vertex %v not found in graph: %v", name, vertex, graph.vertices)
			}
		}
	}
}

func TestUndirected_AddEdge(t *testing.T) {
	tests := map[string]struct {
		vertices      []int
		edgeHashes    [][3]int
		traits        *Traits
		expectedEdges []Edge[int]
		// Even though some AddEdge calls might work, at least one of them could fail, for example
		// if the last call would introduce a cycle.
		shouldFinallyFail bool
	}{
		"graph with 2 edges": {
			vertices:   []int{1, 2, 3},
			edgeHashes: [][3]int{{1, 2, 10}, {1, 3, 20}},
			traits:     &Traits{},
			expectedEdges: []Edge[int]{
				{Source: 1, Target: 2, Properties: EdgeProperties{Weight: 10}},
				{Source: 1, Target: 3, Properties: EdgeProperties{Weight: 20}},
			},
		},
		"hashes for non-existent vertices": {
			vertices:          []int{1, 2},
			edgeHashes:        [][3]int{{1, 3, 20}},
			traits:            &Traits{},
			shouldFinallyFail: true,
		},
		"edge introducing a cycle in an acyclic graph": {
			vertices:   []int{1, 2, 3},
			edgeHashes: [][3]int{{1, 2, 0}, {2, 3, 0}, {3, 1, 0}},
			traits: &Traits{
				IsAcyclic: true,
			},
			shouldFinallyFail: true,
		},
	}

	for name, test := range tests {
		graph := newUndirected(IntHash, test.traits)

		for _, vertex := range test.vertices {
			graph.AddVertex(vertex)
		}

		var err error

		for _, edge := range test.edgeHashes {
			if err = graph.AddEdge(edge[0], edge[1], EdgeWeight(edge[2])); err != nil {
				break
			}
		}

		if test.shouldFinallyFail != (err != nil) {
			t.Fatalf("%s: error expectancy doesn't match: expected %v, got %v (error: %v)", name, test.shouldFinallyFail, (err != nil), err)
		}

		for _, expectedEdge := range test.expectedEdges {
			sourceHash := graph.hash(expectedEdge.Source)
			targetHash := graph.hash(expectedEdge.Target)

			edge, ok := graph.outEdges[sourceHash][targetHash]
			if !ok {
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
		}
	}
}

func TestUndirected_Edge(t *testing.T) {
	tests := map[string]struct {
		vertices      []int
		getEdgeHashes [2]int
		exists        bool
	}{
		"get edge of undirected graph": {
			vertices:      []int{1, 2, 3},
			getEdgeHashes: [2]int{2, 1},
			exists:        true,
		},
		"get non-existent edge of undirected graph": {
			vertices:      []int{1, 2, 3},
			getEdgeHashes: [2]int{1, 3},
		},
	}

	for name, test := range tests {
		graph := newUndirected(IntHash, &Traits{})

		for _, vertex := range test.vertices {
			graph.AddVertex(vertex)
		}

		sourceHash := graph.hash(test.vertices[0])
		targetHash := graph.hash(test.vertices[1])

		_ = graph.AddEdge(sourceHash, targetHash)

		_, ok := graph.Edge(test.getEdgeHashes[0], test.getEdgeHashes[1])

		if test.exists != ok {
			t.Fatalf("%s: result expectancy doesn't match: expected %v, got %v", name, test.exists, ok)
		}
	}
}

func TestUndirected_CreatesCycle(t *testing.T) {
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
		graph := newUndirected(IntHash, &Traits{})

		for _, vertex := range test.vertices {
			graph.AddVertex(vertex)
		}

		for _, edge := range test.edges {
			if err := graph.AddEdge(edge.Source, edge.Target); err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		createsCycle, err := graph.CreatesCycle(test.sourceHash, test.targetHash)
		if err != nil {
			t.Fatalf("%s: failed to add edge: %s", name, err.Error())
		}

		if createsCycle != test.createsCycle {
			t.Errorf("%s: cycle expectancy doesn't match: expected %v, got %v", name, test.createsCycle, createsCycle)
		}
	}
}

func TestUndirected_Degree(t *testing.T) {
	tests := map[string]struct {
		vertices       []int
		edges          []Edge[int]
		vertex         int
		expectedDegree int
		shouldFail     bool
	}{
		"graph with 3 vertices and 2 edges": {
			vertices: []int{1, 2, 3},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
			},
			vertex:         1,
			expectedDegree: 2,
		},
		"graph with 3 vertices and 3 edges": {
			vertices: []int{1, 2, 3},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 2, Target: 3},
				{Source: 3, Target: 1},
			},
			vertex:         2,
			expectedDegree: 2,
		},
		"disconnected graph": {
			vertices: []int{1, 2, 3},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
			},
			vertex:         3,
			expectedDegree: 0,
		},
		"non-existent vertex": {
			vertices:   []int{1, 2, 3},
			vertex:     4,
			shouldFail: true,
		},
	}

	for name, test := range tests {
		graph := newUndirected(IntHash, &Traits{})

		for _, vertex := range test.vertices {
			graph.AddVertex(vertex)
		}

		for _, edge := range test.edges {
			if err := graph.AddEdge(edge.Source, edge.Target); err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		degree, err := graph.Degree(test.vertex)

		if test.shouldFail != (err != nil) {
			t.Fatalf("%s: error expectancy doesn't match: expected %v, got %v (error: %v)", name, test.shouldFail, (err != nil), err)
		}

		if degree != test.expectedDegree {
			t.Errorf("%s: degree expectancy doesn't match: expcted %v, got %v", name, test.expectedDegree, degree)
		}
	}
}

func TestUndirected_StronglyConnectedComponents(t *testing.T) {
	tests := map[string]struct {
		expectedSCCs [][]int
		shouldFail   bool
	}{
		"return error": {
			expectedSCCs: nil,
			shouldFail:   true,
		},
	}

	for name, test := range tests {
		graph := newUndirected(IntHash, &Traits{})

		sccs, err := graph.StronglyConnectedComponents()

		if test.shouldFail != (err != nil) {
			t.Errorf("%s: error expectancy doesn't match: expected %v, got %v (error: %v)", name, test.shouldFail, (err != nil), err)
		}

		if test.expectedSCCs == nil && sccs != nil {
			t.Errorf("%s: SCC expectancy doesn't match: expcted %v, got %v", name, test.expectedSCCs, sccs)
		}
	}
}

func TestUndirected_ShortestPath(t *testing.T) {
	tests := map[string]struct {
		vertices             []string
		edges                []Edge[string]
		sourceHash           string
		targetHash           string
		expectedShortestPath []string
		shouldFail           bool
	}{
		"graph as on img/dijkstra.svg": {
			vertices: []string{"A", "B", "C", "D", "E", "F", "G"},
			edges: []Edge[string]{
				{Source: "A", Target: "C", Properties: EdgeProperties{Weight: 3}},
				{Source: "A", Target: "F", Properties: EdgeProperties{Weight: 2}},
				{Source: "C", Target: "D", Properties: EdgeProperties{Weight: 4}},
				{Source: "C", Target: "E", Properties: EdgeProperties{Weight: 1}},
				{Source: "C", Target: "F", Properties: EdgeProperties{Weight: 2}},
				{Source: "D", Target: "B", Properties: EdgeProperties{Weight: 1}},
				{Source: "E", Target: "B", Properties: EdgeProperties{Weight: 2}},
				{Source: "E", Target: "F", Properties: EdgeProperties{Weight: 3}},
				{Source: "F", Target: "G", Properties: EdgeProperties{Weight: 5}},
				{Source: "G", Target: "B", Properties: EdgeProperties{Weight: 2}},
			},
			sourceHash:           "A",
			targetHash:           "B",
			expectedShortestPath: []string{"A", "C", "E", "B"},
		},
		"diamond-shaped graph": {
			vertices: []string{"A", "B", "C", "D"},
			edges: []Edge[string]{
				{Source: "A", Target: "B", Properties: EdgeProperties{Weight: 2}},
				{Source: "A", Target: "C", Properties: EdgeProperties{Weight: 4}},
				{Source: "B", Target: "D", Properties: EdgeProperties{Weight: 2}},
				{Source: "C", Target: "D", Properties: EdgeProperties{Weight: 2}},
			},
			sourceHash:           "A",
			targetHash:           "D",
			expectedShortestPath: []string{"A", "B", "D"},
		},
		"source equal to target": {
			vertices: []string{"A", "B", "C", "D"},
			edges: []Edge[string]{
				{Source: "A", Target: "B", Properties: EdgeProperties{Weight: 2}},
				{Source: "A", Target: "C", Properties: EdgeProperties{Weight: 4}},
				{Source: "B", Target: "D", Properties: EdgeProperties{Weight: 2}},
				{Source: "C", Target: "D", Properties: EdgeProperties{Weight: 2}},
			},
			sourceHash:           "B",
			targetHash:           "B",
			expectedShortestPath: []string{"B"},
		},
		"target not reachable": {
			vertices: []string{"A", "B", "C", "D"},
			edges: []Edge[string]{
				{Source: "A", Target: "B", Properties: EdgeProperties{Weight: 2}},
				{Source: "A", Target: "C", Properties: EdgeProperties{Weight: 4}},
			},
			sourceHash:           "A",
			targetHash:           "D",
			expectedShortestPath: []string{},
			shouldFail:           true,
		},
	}

	for name, test := range tests {
		graph := newUndirected(StringHash, &Traits{})

		for _, vertex := range test.vertices {
			graph.AddVertex(vertex)
		}

		for _, edge := range test.edges {
			if err := graph.AddEdge(edge.Source, edge.Target, EdgeWeight(edge.Properties.Weight)); err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		shortestPath, err := graph.ShortestPath(test.sourceHash, test.targetHash)

		if test.shouldFail != (err != nil) {
			t.Fatalf("%s: error expectancy doesn't match: expected %v, got %v (error: %v)", name, test.shouldFail, (err != nil), err)
		}

		if len(shortestPath) != len(test.expectedShortestPath) {
			t.Fatalf("%s: path length expectancy doesn't match: expected %v, got %v", name, len(test.expectedShortestPath), len(shortestPath))
		}

		for i, expectedVertex := range test.expectedShortestPath {
			if shortestPath[i] != expectedVertex {
				t.Errorf("%s: path vertex expectancy doesn't match: expected %v at index %d, got %v", name, expectedVertex, i, shortestPath[i])
			}
		}
	}
}

func TestUndirected_AdjacencyList(t *testing.T) {
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
					1: {Source: 3, Target: 1},
					2: {Source: 3, Target: 2},
					4: {Source: 3, Target: 4},
				},
				4: {
					3: {Source: 4, Target: 3},
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
				1: {
					2: {Source: 1, Target: 2},
					3: {Source: 1, Target: 3},
				},
				2: {
					1: {Source: 2, Target: 1},
					4: {Source: 2, Target: 4},
				},
				3: {
					1: {Source: 3, Target: 1},
					4: {Source: 3, Target: 4},
				},
				4: {
					2: {Source: 4, Target: 2},
					3: {Source: 4, Target: 3},
				},
			},
		},
	}

	for name, test := range tests {
		graph := newUndirected(IntHash, &Traits{})

		for _, vertex := range test.vertices {
			graph.AddVertex(vertex)
		}

		for _, edge := range test.edges {
			if err := graph.AddEdge(edge.Source, edge.Target, EdgeWeight(edge.Properties.Weight)); err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		adjacencyMap := graph.AdjacencyMap()

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

func TestUndirected_edgesAreEqual(t *testing.T) {
	tests := map[string]struct {
		a             Edge[int]
		b             Edge[int]
		edgesAreEqual bool
	}{
		"equal edges in undirected graph": {
			a:             Edge[int]{Source: 1, Target: 2},
			b:             Edge[int]{Source: 1, Target: 2},
			edgesAreEqual: true,
		},
		"swapped equal edges in undirected graph": {
			a:             Edge[int]{Source: 1, Target: 2},
			b:             Edge[int]{Source: 2, Target: 1},
			edgesAreEqual: true,
		},
		"unequal edges in undirected graph": {
			a: Edge[int]{Source: 1, Target: 2},
			b: Edge[int]{Source: 1, Target: 3},
		},
	}

	for name, test := range tests {
		graph := newUndirected(IntHash, &Traits{})
		actual := graph.edgesAreEqual(test.a, test.b)

		if actual != test.edgesAreEqual {
			t.Errorf("%s: equality expectations don't match: expected %v, got %v", name, test.edgesAreEqual, actual)
		}
	}
}

func TestUndirected_addEdge(t *testing.T) {
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
		graph := newUndirected(IntHash, &Traits{})

		for _, edge := range test.edges {
			sourceHash := graph.hash(edge.Source)
			TargetHash := graph.hash(edge.Target)
			graph.addEdge(sourceHash, TargetHash, edge)
		}

		if len(graph.outEdges) != len(test.edges) {
			t.Errorf("%s: number of outgoing edges doesn't match: expected %v, got %v", name, len(test.edges), len(graph.outEdges))
		}
		if len(graph.inEdges) != len(test.edges) {
			t.Errorf("%s: number of ingoing edges doesn't match: expected %v, got %v", name, len(test.edges), len(graph.inEdges))
		}
	}
}

func TestUndirected_adjacencies(t *testing.T) {
	tests := map[string]struct {
		vertices             []int
		edges                []Edge[int]
		vertex               int
		expectedAdjancencies []int
	}{
		"graph with 3 vertices": {
			vertices: []int{1, 2, 3},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
			},
			vertex:               2,
			expectedAdjancencies: []int{1},
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
			vertex:               2,
			expectedAdjancencies: []int{1, 4, 5},
		},
		"graph with 7 vertices and a diamond cycle (#1)": {
			vertices: []int{1, 2, 3, 4, 5, 6, 7},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
				{Source: 2, Target: 4},
				{Source: 2, Target: 5},
				{Source: 3, Target: 6},
				{Source: 4, Target: 7},
				{Source: 5, Target: 7},
			},
			vertex:               5,
			expectedAdjancencies: []int{2, 7},
		},
		"graph with 7 vertices and a diamond cycle (#2)": {
			vertices: []int{1, 2, 3, 4, 5, 6, 7},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
				{Source: 2, Target: 4},
				{Source: 2, Target: 5},
				{Source: 3, Target: 6},
				{Source: 4, Target: 7},
				{Source: 5, Target: 7},
			},
			vertex:               7,
			expectedAdjancencies: []int{4, 5},
		},
		"graph with 7 vertices and a diamond cycle (#3)": {
			vertices: []int{1, 2, 3, 4, 5, 6, 7},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
				{Source: 2, Target: 4},
				{Source: 2, Target: 5},
				{Source: 3, Target: 6},
				{Source: 4, Target: 7},
				{Source: 5, Target: 7},
			},
			vertex:               2,
			expectedAdjancencies: []int{1, 4, 5},
		},
	}

	for name, test := range tests {
		graph := newUndirected(IntHash, &Traits{})

		for _, vertex := range test.vertices {
			graph.AddVertex(vertex)
		}

		for _, edge := range test.edges {
			if err := graph.AddEdge(edge.Source, edge.Target); err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		adjacencies := graph.adjacencies(graph.hash(test.vertex))

		if !slicesAreEqual(adjacencies, test.expectedAdjancencies) {
			t.Errorf("%s: adjacencies don't match: expected %v, got %v", name, test.expectedAdjancencies, adjacencies)
		}
	}
}

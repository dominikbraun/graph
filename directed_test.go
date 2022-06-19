package graph

import (
	"testing"
)

func TestDirected_Vertex(t *testing.T) {
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
		graph := newDirected(IntHash, &properties{})

		for _, vertex := range test.vertices {
			graph.Vertex(vertex)
		}

		for _, vertex := range test.vertices {
			hash := graph.hash(vertex)
			if _, ok := graph.vertices[hash]; !ok {
				t.Errorf("%s: vertex %v not found in graph: %v", name, vertex, graph.vertices)
			}
		}
	}
}

func TestDirected_Edge(t *testing.T) {
	TestDirected_WeightedEdge(t)
}

func TestDirected_WeightedEdge(t *testing.T) {
	TestDirected_WeightedEdgeByHashes(t)
}

func TestDirected_EdgeByHashes(t *testing.T) {
	TestDirected_WeightedEdgeByHashes(t)
}

func TestDirected_WeightedEdgeByHashes(t *testing.T) {
	tests := map[string]struct {
		vertices      []int
		edgeHashes    [][3]int
		expectedEdges []Edge[int]
		shouldFail    bool
	}{
		"graph with 2 edges": {
			vertices:   []int{1, 2, 3},
			edgeHashes: [][3]int{{1, 2, 10}, {1, 3, 20}},
			expectedEdges: []Edge[int]{
				{Source: 1, Target: 2, Weight: 10},
				{Source: 1, Target: 3, Weight: 20},
			},
		},
		"hashes for non-existent vertices": {
			vertices:   []int{1, 2},
			edgeHashes: [][3]int{{1, 3, 20}},
			shouldFail: true,
		},
	}

	for name, test := range tests {
		graph := newDirected(IntHash, &properties{})

		for _, vertex := range test.vertices {
			graph.Vertex(vertex)
		}
		for _, edge := range test.edgeHashes {
			err := graph.WeightedEdgeByHashes(edge[0], edge[1], edge[2])

			if test.shouldFail != (err != nil) {
				t.Fatalf("%s: error expectancy doesn't match: expected %v, got %v (error: %v)", name, test.shouldFail, (err != nil), err)
			}
		}
		for _, expectedEdge := range test.expectedEdges {
			sourceHash := graph.hash(expectedEdge.Source)
			targetHash := graph.hash(expectedEdge.Target)

			edge, ok := graph.edges[sourceHash][targetHash]
			if !ok {
				t.Fatalf("%s: edge with source %v and target %v not found", name, expectedEdge.Source, expectedEdge.Target)
			}

			if edge.Source != expectedEdge.Source {
				t.Errorf("%s: edge sources don't match: expected source %v, got %v", name, expectedEdge.Source, edge.Source)
			}

			if edge.Target != expectedEdge.Target {
				t.Errorf("%s: edge targets don't match: expected target %v, got %v", name, expectedEdge.Target, edge.Target)
			}

			if edge.Weight != expectedEdge.Weight {
				t.Errorf("%s: edge weights don't match: expected weight %v, got %v", name, expectedEdge.Weight, edge.Weight)
			}
		}
	}
}

func TestDirected_GetEdge(t *testing.T) {
	TestDirected_GetEdgeByHashes(t)
}

func TestDirected_GetEdgeByHashes(t *testing.T) {
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
		graph := newDirected(IntHash, &properties{})
		for _, vertex := range test.vertices {
			graph.Vertex(vertex)
		}

		sourceHash := graph.hash(test.vertices[0])
		targetHash := graph.hash(test.vertices[1])

		graph.EdgeByHashes(sourceHash, targetHash)

		_, ok := graph.GetEdgeByHashes(test.getEdgeHashes[0], test.getEdgeHashes[1])

		if test.exists != ok {
			t.Fatalf("%s: result expectancy doesn't match: expected %v, got %v", name, test.exists, ok)
		}
	}
}

func TestDirected_DFS(t *testing.T) {
	TestDirected_DFSByHash(t)
}

func TestDirected_DFSByHash(t *testing.T) {
	tests := map[string]struct {
		vertices       []int
		edges          []Edge[int]
		startHash      int
		expectedVisits []int
		stopAtVertex   int
	}{
		"traverse entire graph with 3 vertices": {
			vertices: []int{1, 2, 3},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
			},
			startHash:      1,
			expectedVisits: []int{1, 2, 3},
			stopAtVertex:   -1,
		},
		"traverse entire triangle graph": {
			vertices: []int{1, 2, 3},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 2, Target: 3},
				{Source: 3, Target: 1},
			},
			startHash:      1,
			expectedVisits: []int{1, 2, 3},
			stopAtVertex:   -1,
		},
		"traverse graph with 3 vertices until vertex 2": {
			vertices: []int{1, 2, 3},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 2, Target: 3},
				{Source: 3, Target: 1},
			},
			startHash:      1,
			expectedVisits: []int{1, 2},
			stopAtVertex:   2,
		},
	}

	for name, test := range tests {
		graph := newDirected(IntHash, &properties{})

		for _, vertex := range test.vertices {
			graph.Vertex(vertex)
		}

		for _, edge := range test.edges {
			if err := graph.Edge(edge.Source, edge.Target); err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		visited := make(map[int]struct{})

		visit := func(value int) bool {
			visited[value] = struct{}{}

			if test.stopAtVertex != -1 {
				if value == test.stopAtVertex {
					return true
				}
			}
			return false
		}

		_ = graph.DFSByHash(test.startHash, visit)

		if len(visited) != len(test.expectedVisits) {
			t.Fatalf("%s: numbers of visited vertices don't match: expected %v, got %v", name, len(test.expectedVisits), len(visited))
		}

		for _, expectedVisit := range test.expectedVisits {
			if _, ok := visited[expectedVisit]; !ok {
				t.Errorf("%s: expected vertex %v to be visited, but it isn't", name, expectedVisit)
			}
		}
	}
}

func TestDirected_BFS(t *testing.T) {
	TestDirected_BFSByHash(t)
}

func TestDirected_BFSByHash(t *testing.T) {
	tests := map[string]struct {
		vertices       []int
		edges          []Edge[int]
		startHash      int
		expectedVisits []int
		stopAtVertex   int
	}{
		"traverse entire graph with 3 vertices": {
			vertices: []int{1, 2, 3},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
			},
			startHash:      1,
			expectedVisits: []int{1, 2, 3},
			stopAtVertex:   -1,
		},
		"traverse graph with 6 vertices until vertex 4": {
			vertices: []int{1, 2, 3, 4, 5, 6},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
				{Source: 2, Target: 4},
				{Source: 2, Target: 5},
				{Source: 3, Target: 6},
			},
			startHash:      1,
			expectedVisits: []int{1, 2, 3, 4},
			stopAtVertex:   4,
		},
	}

	for name, test := range tests {
		graph := newDirected(IntHash, &properties{})

		for _, vertex := range test.vertices {
			graph.Vertex(vertex)
		}

		for _, edge := range test.edges {
			if err := graph.Edge(edge.Source, edge.Target); err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		visited := make(map[int]struct{})

		visit := func(value int) bool {
			visited[value] = struct{}{}

			if test.stopAtVertex != -1 {
				if value == test.stopAtVertex {
					return true
				}
			}
			return false
		}

		_ = graph.BFSByHash(test.startHash, visit)

		for _, expectedVisit := range test.expectedVisits {
			if _, ok := visited[expectedVisit]; !ok {
				t.Errorf("%s: expected vertex %v to be visited, but it isn't", name, expectedVisit)
			}
		}
	}
}

func TestDirected_CreatesCycle(t *testing.T) {
	TestDirected_CreatesCycleByHashes(t)
}

func TestDirected_CreatesCycleByHashes(t *testing.T) {
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
		graph := newDirected(IntHash, &properties{})

		for _, vertex := range test.vertices {
			graph.Vertex(vertex)
		}

		for _, edge := range test.edges {
			if err := graph.Edge(edge.Source, edge.Target); err != nil {
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
		graph := newDirected(IntHash, &properties{})
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
				{Source: 1, Target: 2, Weight: 1},
				{Source: 2, Target: 3, Weight: 2},
				{Source: 3, Target: 1, Weight: 3},
			},
		},
	}

	for name, test := range tests {
		graph := newDirected(IntHash, &properties{})

		for _, edge := range test.edges {
			sourceHash := graph.hash(edge.Source)
			TargetHash := graph.hash(edge.Target)
			graph.addEdge(sourceHash, TargetHash, edge)
		}

		if len(graph.edges) != len(test.edges) {
			t.Errorf("%s: number of edges doesn't match: expected %v, got %v", name, len(test.edges), len(graph.edges))
		}
		if len(graph.outEdges) != len(test.edges) {
			t.Errorf("%s: number of outgoing edges doesn't match: expected %v, got %v", name, len(test.edges), len(graph.outEdges))
		}
		if len(graph.inEdges) != len(test.edges) {
			t.Errorf("%s: number of ingoing edges doesn't match: expected %v, got %v", name, len(test.edges), len(graph.inEdges))
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
		graph := newDirected(IntHash, &properties{})

		for _, vertex := range test.vertices {
			graph.Vertex(vertex)
		}

		for _, edge := range test.edges {
			if err := graph.Edge(edge.Source, edge.Target); err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		predecessors := graph.predecessors(graph.hash(test.vertex))

		if !slicesAreEqual(predecessors, test.expectedPredecessors) {
			t.Errorf("%s: predecessors don't match: expected %v, got %v", name, test.expectedPredecessors, predecessors)
		}
	}
}

func slicesAreEqual[T comparable](a []T, b []T) bool {
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

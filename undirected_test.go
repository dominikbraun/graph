package graph

import (
	"testing"
)

func TestUndirected_Vertex(t *testing.T) {
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
		graph := newUndirected(IntHash, newMemoryStore(IntHash), &traits{})

		for _, vertex := range test.vertices {
			graph.Vertex(vertex)
		}

		for _, vertex := range test.vertices {
			hash := graph.hash(vertex)
			if _, ok := graph.store.GetVertex(hash); !ok {
				vertices, _ := graph.store.GetAllVertexHashes()
				t.Errorf("%s: vertex %v not found in graph: %v", name, vertex, vertices)
			}
		}
	}
}

func TestUndirected_Edge(t *testing.T) {
	TestUndirected_WeightedEdge(t)
}

func TestUndirected_WeightedEdge(t *testing.T) {
	TestUndirected_WeightedEdgeByHashes(t)
}

func TestUndirected_EdgeByHashes(t *testing.T) {
	TestUndirected_WeightedEdgeByHashes(t)
}

func TestUndirected_WeightedEdgeByHashes(t *testing.T) {
	tests := map[string]struct {
		vertices      []int
		edgeHashes    [][3]int
		traits        *traits
		expectedEdges []Edge[int]
		// Even though some of the WeightedEdgeByHashes calls might work, at least one of them
		// could fail - for example if the last call would introduce a cycle.
		shouldFinallyFail bool
	}{
		"graph with 2 edges": {
			vertices:   []int{1, 2, 3},
			edgeHashes: [][3]int{{1, 2, 10}, {1, 3, 20}},
			traits:     &traits{},
			expectedEdges: []Edge[int]{
				{Source: 1, Target: 2, Weight: 10},
				{Source: 1, Target: 3, Weight: 20},
			},
		},
		"hashes for non-existent vertices": {
			vertices:          []int{1, 2},
			edgeHashes:        [][3]int{{1, 3, 20}},
			traits:            &traits{},
			shouldFinallyFail: true,
		},
		"edge introducing a cycle in an acyclic graph": {
			vertices:   []int{1, 2, 3},
			edgeHashes: [][3]int{{1, 2, 0}, {2, 3, 0}, {3, 1, 0}},
			traits: &traits{
				isAcyclic: true,
			},
			shouldFinallyFail: true,
		},
	}

	for name, test := range tests {
		graph := newUndirected(IntHash, newMemoryStore(IntHash), test.traits)

		for _, vertex := range test.vertices {
			graph.Vertex(vertex)
		}

		var err error

		for _, edge := range test.edgeHashes {
			if err = graph.WeightedEdgeByHashes(edge[0], edge[1], edge[2]); err != nil {
				break
			}
		}

		if test.shouldFinallyFail != (err != nil) {
			t.Fatalf("%s: error expectancy doesn't match: expected %v, got %v (error: %v)", name, test.shouldFinallyFail, (err != nil), err)
		}

		for _, expectedEdge := range test.expectedEdges {
			sourceHash := graph.hash(expectedEdge.Source)
			targetHash := graph.hash(expectedEdge.Target)

			edge, ok := graph.store.GetEdge(sourceHash, targetHash)
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

func TestUndirected_GetEdge(t *testing.T) {
	TestUndirected_GetEdgeByHashes(t)
}

func TestUndirected_GetEdgeByHashes(t *testing.T) {
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
		graph := newUndirected(IntHash, newMemoryStore(IntHash), &traits{})

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

func TestUndirected_DFS(t *testing.T) {
	TestUndirected_DFSByHash(t)
}

func TestUndirected_DFSByHash(t *testing.T) {
	tests := map[string]struct {
		vertices  []int
		edges     []Edge[int]
		startHash int
		// It is not possible to expect a strict list of vertices to be visited. If stopAtVertex is
		// a neighbor of another vertex, that other vertex might be visited before stopAtVertex.
		expectedMinimumVisits []int
		// In case stopAtVertex has downstream neighbors, those neighbors musn't be visited.
		forbiddenVisits []int
		stopAtVertex    int
	}{
		"traverse entire graph with 3 vertices": {
			vertices: []int{1, 2, 3},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
			},
			startHash:             1,
			expectedMinimumVisits: []int{1, 2, 3},
			stopAtVertex:          -1,
		},
		"traverse entire triangle graph": {
			vertices: []int{1, 2, 3},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 2, Target: 3},
				{Source: 3, Target: 1},
			},
			startHash:             1,
			expectedMinimumVisits: []int{1, 2, 3},
			stopAtVertex:          -1,
		},
		"traverse graph with 3 vertices until vertex 2": {
			vertices: []int{1, 2, 3},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 2, Target: 3},
				{Source: 3, Target: 1},
			},
			startHash:             1,
			expectedMinimumVisits: []int{1, 2},
			stopAtVertex:          2,
		},
		"traverse graph with 7 vertices until vertex 4": {
			vertices: []int{1, 2, 3, 4, 5, 6, 7},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
				{Source: 2, Target: 4},
				{Source: 2, Target: 5},
				{Source: 4, Target: 6},
				{Source: 5, Target: 7},
			},
			startHash:             1,
			expectedMinimumVisits: []int{1, 2, 4},
			forbiddenVisits:       []int{6},
			stopAtVertex:          4,
		},
		"traverse graph with 15 vertices until vertex 11": {
			vertices: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
				{Source: 3, Target: 4},
				{Source: 3, Target: 5},
				{Source: 3, Target: 6},
				{Source: 4, Target: 7},
				{Source: 5, Target: 13},
				{Source: 5, Target: 14},
				{Source: 6, Target: 7},
				{Source: 7, Target: 8},
				{Source: 7, Target: 9},
				{Source: 8, Target: 10},
				{Source: 9, Target: 11},
				{Source: 9, Target: 12},
				{Source: 10, Target: 14},
				{Source: 11, Target: 15},
			},
			startHash:             1,
			expectedMinimumVisits: []int{1, 3, 7, 9, 11},
			forbiddenVisits:       []int{15},
			stopAtVertex:          11,
		},
		"traverse a disconnected graph": {
			vertices: []int{1, 2, 3, 4},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 3, Target: 4},
			},
			startHash:             1,
			expectedMinimumVisits: []int{1, 2},
			stopAtVertex:          -1,
		},
	}

	for name, test := range tests {
		graph := newUndirected(IntHash, newMemoryStore(IntHash), &traits{})

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

		if len(visited) < len(test.expectedMinimumVisits) {
			t.Fatalf("%s: expected number of minimum visits doesn't match: expected %v, got %v", name, len(test.expectedMinimumVisits), len(visited))
		}

		if test.forbiddenVisits != nil {
			for _, forbiddenVisit := range test.forbiddenVisits {
				if _, ok := visited[forbiddenVisit]; ok {
					t.Errorf("%s: expected vertex %v to not be visited, but it is", name, forbiddenVisit)
				}
			}
		}

		for _, expectedVisit := range test.expectedMinimumVisits {
			if _, ok := visited[expectedVisit]; !ok {
				t.Errorf("%s: expected vertex %v to be visited, but it isn't", name, expectedVisit)
			}
		}
	}
}

func TestUndirected_BFS(t *testing.T) {
	TestUndirected_BFSByHash(t)
}

func TestUndirected_BFSByHash(t *testing.T) {
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
		"traverse a disconnected graph": {
			vertices: []int{1, 2, 3, 4},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 3, Target: 4},
			},
			startHash:      1,
			expectedVisits: []int{1, 2},
			stopAtVertex:   -1,
		},
	}

	for name, test := range tests {
		graph := newUndirected(IntHash, newMemoryStore(IntHash), &traits{})

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

func TestUndirected_CreatesCycle(t *testing.T) {
	TestUndirected_CreatesCycleByHashes(t)
}

func TestUndirected_CreatesCycleByHashes(t *testing.T) {
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
		graph := newUndirected(IntHash, newMemoryStore(IntHash), &traits{})

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

func TestUndirected_Degree(t *testing.T) {
	TestDirected_Degree(t)
}

func TestUndirected_DegreeByHash(t *testing.T) {
	tests := map[string]struct {
		vertices       []int
		edges          []Edge[int]
		vertexHash     int
		expectedDegree int
		shouldFail     bool
	}{}

	for name, test := range tests {
		graph := newUndirected(IntHash, newMemoryStore(IntHash), &traits{})

		for _, vertex := range test.vertices {
			graph.Vertex(vertex)
		}

		for _, edge := range test.edges {
			if err := graph.Edge(edge.Source, edge.Target); err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		degree, err := graph.DegreeByHash(test.vertexHash)

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
		graph := newUndirected(IntHash, newMemoryStore(IntHash), &traits{})

		sccs, err := graph.StronglyConnectedComponents()

		if test.shouldFail != (err != nil) {
			t.Errorf("%s: error expectancy doesn't match: expected %v, got %v (error: %v)", name, test.shouldFail, (err != nil), err)
		}

		if test.expectedSCCs == nil && sccs != nil {
			t.Errorf("%s: SCC expectancy doesn't match: expcted %v, got %v", name, test.expectedSCCs, sccs)
		}
	}
}

func TestUndirected_ShortesPath(t *testing.T) {
	TestUndirected_ShortestPathByHashes(t)
}

func TestUndirected_ShortestPathByHashes(t *testing.T) {
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
				{Source: "A", Target: "C", Weight: 3},
				{Source: "A", Target: "F", Weight: 2},
				{Source: "C", Target: "D", Weight: 4},
				{Source: "C", Target: "E", Weight: 1},
				{Source: "C", Target: "F", Weight: 2},
				{Source: "D", Target: "B", Weight: 1},
				{Source: "E", Target: "B", Weight: 2},
				{Source: "E", Target: "F", Weight: 3},
				{Source: "F", Target: "G", Weight: 5},
				{Source: "G", Target: "B", Weight: 2},
			},
			sourceHash:           "A",
			targetHash:           "B",
			expectedShortestPath: []string{"A", "C", "E", "B"},
		},
		"diamond-shaped graph": {
			vertices: []string{"A", "B", "C", "D"},
			edges: []Edge[string]{
				{Source: "A", Target: "B", Weight: 2},
				{Source: "A", Target: "C", Weight: 4},
				{Source: "B", Target: "D", Weight: 2},
				{Source: "C", Target: "D", Weight: 2},
			},
			sourceHash:           "A",
			targetHash:           "D",
			expectedShortestPath: []string{"A", "B", "D"},
		},
		"source equal to target": {
			vertices: []string{"A", "B", "C", "D"},
			edges: []Edge[string]{
				{Source: "A", Target: "B", Weight: 2},
				{Source: "A", Target: "C", Weight: 4},
				{Source: "B", Target: "D", Weight: 2},
				{Source: "C", Target: "D", Weight: 2},
			},
			sourceHash:           "B",
			targetHash:           "B",
			expectedShortestPath: []string{"B"},
		},
		"target not reachable": {
			vertices: []string{"A", "B", "C", "D"},
			edges: []Edge[string]{
				{Source: "A", Target: "B", Weight: 2},
				{Source: "A", Target: "C", Weight: 4},
			},
			sourceHash:           "A",
			targetHash:           "D",
			expectedShortestPath: []string{},
			shouldFail:           true,
		},
	}

	for name, test := range tests {
		graph := newUndirected(StringHash, newMemoryStore(StringHash), &traits{})

		for _, vertex := range test.vertices {
			graph.Vertex(vertex)
		}

		for _, edge := range test.edges {
			if err := graph.WeightedEdge(edge.Source, edge.Target, edge.Weight); err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		shortestPath, err := graph.ShortestPathByHashes(test.sourceHash, test.targetHash)

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
		graph := newUndirected(IntHash, newMemoryStore(IntHash), &traits{})
		actual := graph.edgesAreEqual(test.a, test.b)

		if actual != test.edgesAreEqual {
			t.Errorf("%s: equality expectations don't match: expected %v, got %v", name, test.edgesAreEqual, actual)
		}
	}
}

// TODO(geoah): update test
// func TestUndirected_addEdge(t *testing.T) {
// 	tests := map[string]struct {
// 		edges []Edge[int]
// 	}{
// 		"add 3 edges": {
// 			edges: []Edge[int]{
// 				{Source: 1, Target: 2, Weight: 1},
// 				{Source: 2, Target: 3, Weight: 2},
// 				{Source: 3, Target: 1, Weight: 3},
// 			},
// 		},
// 	}

// 	for name, test := range tests {
// 		graph := newUndirected(IntHash, newMemoryStore(IntHash), &traits{})

// 		for _, edge := range test.edges {
// 			sourceHash := graph.hash(edge.Source)
// 			TargetHash := graph.hash(edge.Target)
// 			graph.addEdge(sourceHash, TargetHash, edge)
// 		}

// 		if len(graph.edges) != len(test.edges) {
// 			t.Errorf("%s: number of edges doesn't match: expected %v, got %v", name, len(test.edges), len(graph.edges))
// 		}
// 		if len(graph.outEdges) != len(test.edges) {
// 			t.Errorf("%s: number of outgoing edges doesn't match: expected %v, got %v", name, len(test.edges), len(graph.outEdges))
// 		}
// 		if len(graph.inEdges) != len(test.edges) {
// 			t.Errorf("%s: number of ingoing edges doesn't match: expected %v, got %v", name, len(test.edges), len(graph.inEdges))
// 		}
// 	}
// }

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
		graph := newUndirected(IntHash, newMemoryStore(IntHash), &traits{})

		for _, vertex := range test.vertices {
			graph.Vertex(vertex)
		}

		for _, edge := range test.edges {
			if err := graph.Edge(edge.Source, edge.Target); err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		adjacencies := graph.adjacencies(graph.hash(test.vertex))

		if !slicesAreEqual(adjacencies, test.expectedAdjancencies) {
			t.Errorf("%s: adjacencies don't match: expected %v, got %v", name, test.expectedAdjancencies, adjacencies)
		}
	}
}

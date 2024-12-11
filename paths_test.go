package graph

import (
	"reflect"
	"sort"
	"testing"
)

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
			_ = graph.AddVertex(vertex)
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
			_ = graph.AddVertex(vertex)
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

func TestDirectedShortestPath(t *testing.T) {
	tests := map[string]struct {
		vertices             []string
		edges                []Edge[string]
		isWeighted           bool
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
			isWeighted:           true,
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
			isWeighted:           true,
			sourceHash:           "A",
			targetHash:           "D",
			expectedShortestPath: []string{"A", "B", "D"},
		},
		"unweighted graph": {
			vertices: []string{"A", "B", "C", "D", "E", "F", "G"},
			edges: []Edge[string]{
				{Source: "A", Target: "B", Properties: EdgeProperties{}},
				{Source: "A", Target: "C", Properties: EdgeProperties{}},
				{Source: "B", Target: "D", Properties: EdgeProperties{}},
				{Source: "C", Target: "F", Properties: EdgeProperties{}},
				{Source: "D", Target: "G", Properties: EdgeProperties{}},
				{Source: "E", Target: "G", Properties: EdgeProperties{}},
				{Source: "F", Target: "E", Properties: EdgeProperties{}},
			},
			sourceHash:           "A",
			targetHash:           "G",
			expectedShortestPath: []string{"A", "B", "D", "G"},
		},
		"source equal to target": {
			vertices: []string{"A", "B", "C", "D"},
			edges: []Edge[string]{
				{Source: "A", Target: "B", Properties: EdgeProperties{Weight: 2}},
				{Source: "A", Target: "C", Properties: EdgeProperties{Weight: 4}},
				{Source: "B", Target: "D", Properties: EdgeProperties{Weight: 2}},
				{Source: "C", Target: "D", Properties: EdgeProperties{Weight: 2}},
			},
			isWeighted:           true,
			sourceHash:           "B",
			targetHash:           "B",
			expectedShortestPath: []string{"B"},
		},
		"target not reachable in a disconnected graph": {
			vertices: []string{"A", "B", "C", "D"},
			edges: []Edge[string]{
				{Source: "A", Target: "B", Properties: EdgeProperties{Weight: 2}},
				{Source: "A", Target: "C", Properties: EdgeProperties{Weight: 4}},
			},
			isWeighted:           true,
			sourceHash:           "A",
			targetHash:           "D",
			expectedShortestPath: []string{},
			shouldFail:           true,
		},
		"target not reachable in a connected graph": {
			vertices: []string{"A", "B", "C"},
			edges: []Edge[string]{
				{Source: "A", Target: "B", Properties: EdgeProperties{Weight: 0}},
				{Source: "A", Target: "C", Properties: EdgeProperties{Weight: 0}},
			},
			sourceHash:           "B",
			targetHash:           "C",
			expectedShortestPath: []string{},
			shouldFail:           true,
		},
		"graph from issue 88": {
			vertices: []string{"A", "B", "C", "D"},
			edges: []Edge[string]{
				{Source: "A", Target: "B", Properties: EdgeProperties{Weight: 2}},
				{Source: "A", Target: "C", Properties: EdgeProperties{Weight: 6}},
				{Source: "B", Target: "C", Properties: EdgeProperties{Weight: 3}},
				{Source: "B", Target: "D", Properties: EdgeProperties{Weight: 5}},
				{Source: "C", Target: "D", Properties: EdgeProperties{Weight: 1}},
			},
			isWeighted:           true,
			sourceHash:           "A",
			targetHash:           "D",
			expectedShortestPath: []string{"A", "B", "C", "D"},
		},
	}

	for name, test := range tests {
		graph := New(StringHash, Directed())
		graph.(*directed[string, string]).traits.IsWeighted = test.isWeighted

		for _, vertex := range test.vertices {
			_ = graph.AddVertex(vertex)
		}

		for _, edge := range test.edges {
			if err := graph.AddEdge(edge.Source, edge.Target, EdgeWeight(edge.Properties.Weight)); err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		shortestPath, err := ShortestPath(graph, test.sourceHash, test.targetHash)

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

func TestUndirectedShortestPath(t *testing.T) {
	tests := map[string]struct {
		vertices             []string
		edges                []Edge[string]
		sourceHash           string
		targetHash           string
		isWeighted           bool
		isDirected           bool
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
			isWeighted:           true,
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
			isWeighted:           true,
			sourceHash:           "A",
			targetHash:           "D",
			expectedShortestPath: []string{"A", "B", "D"},
		},
		"unweighted graph": {
			vertices: []string{"A", "B", "C", "D", "E", "F", "G"},
			edges: []Edge[string]{
				{Source: "A", Target: "B", Properties: EdgeProperties{}},
				{Source: "A", Target: "C", Properties: EdgeProperties{}},
				{Source: "B", Target: "D", Properties: EdgeProperties{}},
				{Source: "C", Target: "F", Properties: EdgeProperties{}},
				{Source: "D", Target: "G", Properties: EdgeProperties{}},
				{Source: "E", Target: "G", Properties: EdgeProperties{}},
				{Source: "F", Target: "E", Properties: EdgeProperties{}},
			},
			sourceHash:           "A",
			targetHash:           "G",
			expectedShortestPath: []string{"A", "B", "D", "G"},
		},
		"source equal to target": {
			vertices: []string{"A", "B", "C", "D"},
			edges: []Edge[string]{
				{Source: "A", Target: "B", Properties: EdgeProperties{Weight: 2}},
				{Source: "A", Target: "C", Properties: EdgeProperties{Weight: 4}},
				{Source: "B", Target: "D", Properties: EdgeProperties{Weight: 2}},
				{Source: "C", Target: "D", Properties: EdgeProperties{Weight: 2}},
			},
			isWeighted:           true,
			sourceHash:           "B",
			targetHash:           "B",
			expectedShortestPath: []string{"B"},
		},
		"can process negative weights": {
			vertices: []string{"A", "B", "C", "D", "E"},
			edges: []Edge[string]{
				{Source: "A", Target: "B", Properties: EdgeProperties{Weight: 1}},
				{Source: "A", Target: "C", Properties: EdgeProperties{Weight: 2}},
				{Source: "B", Target: "C", Properties: EdgeProperties{Weight: 2}},
				{Source: "B", Target: "D", Properties: EdgeProperties{Weight: 2}},
				{Source: "C", Target: "E", Properties: EdgeProperties{Weight: 2}},
				{Source: "D", Target: "E", Properties: EdgeProperties{Weight: -1}},
			},
			isWeighted:           true,
			isDirected:           true,
			sourceHash:           "A",
			targetHash:           "E",
			expectedShortestPath: []string{"A", "B", "D", "E"},
		},
		"target not reachable in a disconnected graph": {
			vertices: []string{"A", "B", "C", "D"},
			edges: []Edge[string]{
				{Source: "A", Target: "B", Properties: EdgeProperties{Weight: 2}},
				{Source: "A", Target: "C", Properties: EdgeProperties{Weight: 4}},
			},
			isWeighted:           true,
			sourceHash:           "A",
			targetHash:           "D",
			expectedShortestPath: []string{},
			shouldFail:           true,
		},
	}

	for name, test := range tests {
		var graph Graph[string, string]
		if test.isDirected {
			graph = New(StringHash, Directed(), Weighted())
		} else {
			graph = New(StringHash, Weighted())
		}

		for _, vertex := range test.vertices {
			_ = graph.AddVertex(vertex)
		}

		for _, edge := range test.edges {
			if err := graph.AddEdge(edge.Source, edge.Target, EdgeWeight(edge.Properties.Weight)); err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		shortestPath, err := ShortestPath(graph, test.sourceHash, test.targetHash)

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

func Test_BellmanFord(t *testing.T) {
	tests := map[string]struct {
		vertices             []string
		edges                []Edge[string]
		sourceHash           string
		targetHash           string
		isWeighted           bool
		IsDirected           bool
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
			isWeighted:           true,
			IsDirected:           true,
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
			isWeighted:           true,
			IsDirected:           true,
			sourceHash:           "A",
			targetHash:           "D",
			expectedShortestPath: []string{"A", "B", "D"},
		},
		"unweighted graph": {
			vertices: []string{"A", "B", "C", "D", "E", "F", "G"},
			edges: []Edge[string]{
				{Source: "A", Target: "B", Properties: EdgeProperties{}},
				{Source: "A", Target: "C", Properties: EdgeProperties{}},
				{Source: "B", Target: "D", Properties: EdgeProperties{}},
				{Source: "C", Target: "F", Properties: EdgeProperties{}},
				{Source: "D", Target: "G", Properties: EdgeProperties{}},
				{Source: "E", Target: "G", Properties: EdgeProperties{}},
				{Source: "F", Target: "E", Properties: EdgeProperties{}},
			},
			IsDirected:           true,
			sourceHash:           "A",
			targetHash:           "G",
			expectedShortestPath: []string{"A", "B", "D", "G"},
		},
		"source equal to target": {
			vertices: []string{"A", "B", "C", "D"},
			edges: []Edge[string]{
				{Source: "A", Target: "B", Properties: EdgeProperties{Weight: 2}},
				{Source: "A", Target: "C", Properties: EdgeProperties{Weight: 4}},
				{Source: "B", Target: "D", Properties: EdgeProperties{Weight: 2}},
				{Source: "C", Target: "D", Properties: EdgeProperties{Weight: 2}},
			},
			isWeighted:           true,
			IsDirected:           true,
			sourceHash:           "B",
			targetHash:           "B",
			expectedShortestPath: []string{"B"},
		},
		"target not reachable in a disconnected graph": {
			vertices: []string{"A", "B", "C", "D"},
			edges: []Edge[string]{
				{Source: "A", Target: "B", Properties: EdgeProperties{Weight: 2}},
				{Source: "A", Target: "C", Properties: EdgeProperties{Weight: 4}},
			},
			isWeighted:           true,
			IsDirected:           true,
			sourceHash:           "A",
			targetHash:           "D",
			expectedShortestPath: []string{},
			shouldFail:           true,
		},
		"negative weights graph": {
			vertices: []string{"A", "B", "C", "D", "E"},
			edges: []Edge[string]{
				{Source: "A", Target: "B", Properties: EdgeProperties{Weight: 1}},
				{Source: "A", Target: "C", Properties: EdgeProperties{Weight: 2}},
				{Source: "B", Target: "C", Properties: EdgeProperties{Weight: 2}},
				{Source: "B", Target: "D", Properties: EdgeProperties{Weight: 2}},
				{Source: "C", Target: "E", Properties: EdgeProperties{Weight: 2}},
				{Source: "D", Target: "E", Properties: EdgeProperties{Weight: -1}},
			},
			isWeighted:           true,
			IsDirected:           true,
			sourceHash:           "A",
			targetHash:           "E",
			expectedShortestPath: []string{"A", "B", "D", "E"},
		},
		"fails on negative cycles": {
			vertices: []string{"A", "B", "C", "D", "E"},
			edges: []Edge[string]{
				{Source: "A", Target: "B", Properties: EdgeProperties{Weight: 1}},
				{Source: "A", Target: "C", Properties: EdgeProperties{Weight: 4}},
				{Source: "B", Target: "C", Properties: EdgeProperties{Weight: 2}},
				{Source: "B", Target: "D", Properties: EdgeProperties{Weight: 6}},
				{Source: "C", Target: "D", Properties: EdgeProperties{Weight: 3}},
				{Source: "C", Target: "E", Properties: EdgeProperties{Weight: 2}},
				{Source: "D", Target: "E", Properties: EdgeProperties{Weight: -3}},
				{Source: "E", Target: "C", Properties: EdgeProperties{Weight: -3}},
			},
			isWeighted:           true,
			IsDirected:           true,
			sourceHash:           "A",
			targetHash:           "E",
			expectedShortestPath: []string{},
			shouldFail:           true,
		},
		"fails if not directed": {
			vertices: []string{"A", "B", "C", "D"},
			edges: []Edge[string]{
				{Source: "A", Target: "B", Properties: EdgeProperties{Weight: 2}},
				{Source: "A", Target: "C", Properties: EdgeProperties{Weight: 4}},
			},
			isWeighted:           true,
			sourceHash:           "A",
			targetHash:           "D",
			expectedShortestPath: []string{},
			shouldFail:           true,
		},
	}
	for name, test := range tests {
		graph := New(StringHash, Directed(), Weighted())

		for _, vertex := range test.vertices {
			_ = graph.AddVertex(vertex)
		}

		for _, edge := range test.edges {
			if err := graph.AddEdge(edge.Source, edge.Target, EdgeWeight(edge.Properties.Weight)); err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		shortestPath, err := bellmanFord(graph, test.sourceHash, test.targetHash)

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

func TestDirectedStronglyConnectedComponents(t *testing.T) {
	tests := map[string]struct {
		vertices     []int
		edges        []Edge[int]
		expectedSCCs [][]int
	}{
		"graph with SCCs as on img/scc.svg": {
			vertices: []int{1, 2, 3, 4, 5, 6, 7, 8},
			edges: []Edge[int]{
				{Source: 1, Target: 2},
				{Source: 2, Target: 3},
				{Source: 2, Target: 5},
				{Source: 2, Target: 6},
				{Source: 3, Target: 4},
				{Source: 3, Target: 7},
				{Source: 4, Target: 3},
				{Source: 4, Target: 8},
				{Source: 5, Target: 1},
				{Source: 5, Target: 6},
				{Source: 6, Target: 7},
				{Source: 7, Target: 6},
				{Source: 8, Target: 4},
				{Source: 8, Target: 7},
			},
			expectedSCCs: [][]int{{1, 2, 5}, {3, 4, 8}, {6, 7}},
		},
	}

	for name, test := range tests {
		graph := New(IntHash, Directed())

		for _, vertex := range test.vertices {
			_ = graph.AddVertex(vertex)
		}

		for _, edge := range test.edges {
			if err := graph.AddEdge(edge.Source, edge.Target); err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		sccs, _ := StronglyConnectedComponents(graph)
		matchedSCCs := 0

		for _, scc := range sccs {
			for _, expectedSCC := range test.expectedSCCs {
				if slicesAreEqual(scc, expectedSCC) {
					matchedSCCs++
				}
			}
		}

		if matchedSCCs != len(test.expectedSCCs) {
			t.Errorf("%s: expected SCCs don't match: expected %v, got %v", name, test.expectedSCCs, sccs)
		}
	}
}

func TestUndirectedStronglyConnectedComponents(t *testing.T) {
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
		graph := New(IntHash)

		sccs, err := StronglyConnectedComponents(graph)

		if test.shouldFail != (err != nil) {
			t.Errorf("%s: error expectancy doesn't match: expected %v, got %v (error: %v)", name, test.shouldFail, (err != nil), err)
		}

		if test.expectedSCCs == nil && sccs != nil {
			t.Errorf("%s: SCC expectancy doesn't match: expcted %v, got %v", name, test.expectedSCCs, sccs)
		}
	}
}

func TestAllPathsBetween(t *testing.T) {
	type args[K comparable, T any] struct {
		g     Graph[K, T]
		start K
		end   K
	}
	type testCase[K comparable, T any] struct {
		name    string
		args    args[K, T]
		want    [][]K
		wantErr bool
	}
	tests := []testCase[int, int]{
		{
			name: "directed",
			args: args[int, int]{
				g: func() Graph[int, int] {
					g := New(IntHash, Directed())
					for i := 0; i <= 8; i++ {
						_ = g.AddVertex(i)
					}
					_ = g.AddEdge(0, 2)
					_ = g.AddEdge(1, 0)
					_ = g.AddEdge(1, 4)
					_ = g.AddEdge(2, 6)
					_ = g.AddEdge(3, 1)
					_ = g.AddEdge(3, 7)
					_ = g.AddEdge(4, 5)
					_ = g.AddEdge(5, 2)
					_ = g.AddEdge(5, 6)
					_ = g.AddEdge(6, 8)
					_ = g.AddEdge(7, 4)
					return g
				}(),
				start: 3,
				end:   6,
			},
			want: [][]int{
				{3, 1, 0, 2, 6},
				{3, 1, 4, 5, 6},
				{3, 1, 4, 5, 2, 6},
				{3, 7, 4, 5, 2, 6},
				{3, 7, 4, 5, 6},
			},
			wantErr: false,
		},
		{
			name: "undirected",
			args: args[int, int]{
				g: func() Graph[int, int] {
					g := New(IntHash)
					for i := 0; i <= 8; i++ {
						_ = g.AddVertex(i)
					}
					_ = g.AddEdge(0, 1)
					_ = g.AddEdge(0, 2)
					_ = g.AddEdge(1, 3)
					_ = g.AddEdge(1, 4)
					_ = g.AddEdge(2, 5)
					_ = g.AddEdge(2, 6)
					_ = g.AddEdge(3, 7)
					_ = g.AddEdge(4, 5)
					_ = g.AddEdge(4, 7)
					_ = g.AddEdge(5, 6)
					_ = g.AddEdge(6, 8)
					return g
				}(),
				start: 3,
				end:   6,
			},
			want: [][]int{
				{3, 1, 0, 2, 6},
				{3, 1, 0, 2, 5, 6},
				{3, 1, 4, 5, 6},
				{3, 1, 4, 5, 2, 6},
				{3, 7, 4, 5, 2, 6},
				{3, 7, 4, 5, 6},
				{3, 7, 4, 1, 0, 2, 6},
				{3, 7, 4, 1, 0, 2, 5, 6},
			},
			wantErr: false,
		},
		{
			name: "undirected (complex)",
			args: args[int, int]{
				g: func() Graph[int, int] {
					g := New(IntHash)
					for i := 0; i <= 9; i++ {
						_ = g.AddVertex(i)
					}
					_ = g.AddEdge(0, 1)
					_ = g.AddEdge(0, 2)
					_ = g.AddEdge(0, 3)
					_ = g.AddEdge(2, 3)
					_ = g.AddEdge(2, 6)
					_ = g.AddEdge(3, 4)
					_ = g.AddEdge(4, 8)
					_ = g.AddEdge(5, 6)
					_ = g.AddEdge(5, 8)
					_ = g.AddEdge(5, 9)
					_ = g.AddEdge(6, 7)
					_ = g.AddEdge(7, 9)
					_ = g.AddEdge(8, 9)
					return g
				}(),
				start: 0,
				end:   9,
			},
			want: [][]int{
				{0, 2, 3, 4, 8, 5, 6, 7, 9},
				{0, 2, 3, 4, 8, 5, 9},
				{0, 2, 3, 4, 8, 9},
				{0, 2, 6, 5, 8, 9},
				{0, 2, 6, 5, 9},
				{0, 2, 6, 7, 9},
				{0, 3, 2, 6, 5, 8, 9},
				{0, 3, 2, 6, 5, 9},
				{0, 3, 2, 6, 7, 9},
				{0, 3, 4, 8, 5, 6, 7, 9},
				{0, 3, 4, 8, 5, 9},
				{0, 3, 4, 8, 9},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AllPathsBetween(tt.args.g, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("AllPathsBetween() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			toStr := func(s []int) string {
				var num string
				for _, n := range s {
					num = num + string(rune(n))
				}
				return num
			}

			sort.Slice(got, func(i, j int) bool {
				return toStr(got[i]) < toStr(got[j])
			})

			sort.Slice(tt.want, func(i, j int) bool {
				return toStr(tt.want[i]) < toStr(tt.want[j])
			})

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AllPathsBetween() got = %v, want %v", got, tt.want)
			}
		})
	}
}

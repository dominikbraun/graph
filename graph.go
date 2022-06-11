package graph

import (
	"fmt"
)

// Graph represents a generic graph data structure consisting of vertices and nodes. Its vertices
// are of type T and each vertex is identified by a hash of type K.
//
// At the moment, Graph is not suited for representing a multigraph.
type Graph[K comparable, T any] struct {
	hash       Hash[K, T]
	properties *Properties
	vertices   map[K]T
	edges      map[K]map[K]Edge[T]
	outEdges   map[K]map[K]Edge[T]
	inEdges    map[K]map[K]Edge[T]
}

// Edge represents a graph edge with a source and target vertex as well as a weight, which has the
// same value for all edges in an unweighted graph. Even though the vertices are referred to as
// source and target, whether the graph is directed or not is determined by its properties.
type Edge[T any] struct {
	Source T
	Target T
	Weight int
}

// New creates a new graph with vertices of type T, identified by hash values of type K. These hash
// values will be obtained using the provided hash function (see Hash).
//
// For primitive vertex values, you may use the predefined hashing functions. As an example, this
// graph stores integer vertices:
//
//	g := graph.New(graph.IntHash)
//	g.Vertex(1)
//	g.Vertex(2)
//	g.Vertex(3)
//
// The provided IntHash hashing function takes an integer and uses it as a hash value at the same
// time. In a more complex scenario with custom objects, you should define your own function:
//
//	type City struct {
//		Name string
//	}
//
//	cityHash := func(c City) string {
//		return c.Name
//	}
//
//	g := graph.New(cityHash)
//	g.Vertex(london)
//
// This graph will store vertices of type City, identified by hashes of type string. Both type
// parameters can be inferred from the hashing function.
//
// All properties of the graph can be set using the predefined functional options. They can be
// combined arbitrarily. This example creates a directed acyclic graph:
//
//	g := graph.New(graph.IntHash, graph.Directed(), graph.Acyclic())
//
// The behavior of all graph methods is controlled by these particular options.
func New[K comparable, T any](hash Hash[K, T], options ...func(*Properties)) *Graph[K, T] {
	g := Graph[K, T]{
		hash:       hash,
		properties: &Properties{},
		vertices:   make(map[K]T),
		edges:      make(map[K]map[K]Edge[T]),
		outEdges:   make(map[K]map[K]Edge[T]),
		inEdges:    make(map[K]map[K]Edge[T]),
	}

	for _, option := range options {
		option(g.properties)
	}

	return &g
}

// Vertex creates a new vertex in the graph, which won't be connected to another vertex yet. This
// function is idempotent, but overwrites an existing vertex if the hash already exists.
func (g *Graph[K, T]) Vertex(value T) {
	hash := g.hash(value)
	g.vertices[hash] = value
}

// Edge creates an edge between the source and the target vertex. If the Directed option has been
// called on the graph, this is a directed edge. Returns an error if either vertex doesn't exist or
// the edge already exists.
func (g *Graph[K, T]) Edge(source, target T) error {
	return g.WeightedEdge(source, target, 0)
}

// WeightedEdge does the same as Edge, but adds an additional weight to the created edge. In an
// unweighted graph, all edges have the same weight of 0.
func (g *Graph[K, T]) WeightedEdge(source, target T, weight int) error {
	sourceHash := g.hash(source)
	targetHash := g.hash(target)

	return g.WeightedEdgeByHashes(sourceHash, targetHash, weight)
}

// EdgeByHashes creates an edge between the source and the target vertex, but uses hash values to
// identify the vertices. This is convenient when you don't have the full vertex objects at hand.
// Returns an error if either vertex doesn't exist or the edge already exists.
//
// To obtain the hash value for a vertex, call the hashing function passed to New.
func (g *Graph[K, T]) EdgeByHashes(sourceHash, targetHash K) error {
	return g.WeightedEdgeByHashes(sourceHash, targetHash, 0)
}

// WeightedEdgeByHashes does the same as EdgeByHashes, but adds an additional weight to the created
// edge. In an unweighted graph, all edges have the same weight of 0.
func (g *Graph[K, T]) WeightedEdgeByHashes(sourceHash, targetHash K, weight int) error {
	source, ok := g.vertices[sourceHash]
	if !ok {
		return fmt.Errorf("could not find source vertex with hash %v", source)
	}

	target, ok := g.vertices[targetHash]
	if !ok {
		return fmt.Errorf("could not find target vertex with hash %v", source)
	}

	if _, ok := g.GetEdgeByHashes(sourceHash, targetHash); ok {
		return fmt.Errorf("an edge between vertices %v and %v already exists", sourceHash, targetHash)
	}

	edge := Edge[T]{
		Source: source,
		Target: target,
		Weight: weight,
	}

	g.addEdge(sourceHash, targetHash, edge)

	return nil
}

// GetEdgeByHashes returns the edge between two vertices. The second return value indicates whether
// the edge exists. If the graph  is undirected, an edge with swapped source and target vertices
// does match.
func (g *Graph[K, T]) GetEdge(source, target T) (Edge[T], bool) {
	sourceHash := g.hash(source)
	targetHash := g.hash(target)

	return g.GetEdgeByHashes(sourceHash, targetHash)
}

// GetEdgeByHashes returns the edge between two vertices with the given hash values. The second
// return value indicates whether the edge exists. If the graph  is undirected, an edge with
// swapped source and target vertices does match.
func (g *Graph[K, T]) GetEdgeByHashes(sourceHash, targetHash K) (Edge[T], bool) {
	sourceEdges, ok := g.edges[sourceHash]
	if !ok && g.properties.isDirected {
		return Edge[T]{}, false
	}

	if edge, ok := sourceEdges[targetHash]; ok {
		return edge, true
	}

	if !g.properties.isDirected {
		targetEdges, ok := g.edges[targetHash]
		if !ok {
			return Edge[T]{}, false
		}

		if edge, ok := targetEdges[sourceHash]; ok {
			return edge, true
		}
	}

	return Edge[T]{}, false
}

// DFS performs a Depth-First Search on the graph, starting from the given vertex. The visit
// function will be invoked for each visited vertex. If it returns false, DFS will continue
// traversing the path, and if it returns true, the traversal will be stopped.
//
// This example prints all vertices of the graph in DFS-order:
//
//	g := graph.New(graph.IntHash)
//
//	g.Vertex(1)
//	g.Vertex(2)
//	g.Vertex(3)
//
//	_ = g.Edge(1, 2)
//	_ = g.Edge(2, 3)
//	_ = g.Edge(3, 1)
//
//	_ = g.DFS(1, func(value int) bool {
//		fmt.Println(value)
//		return false
//	})
//
// Similarily, if you have a graph of City vertices and the traversal should stop at London,
// the visit function would look as follows:
//
//	func(city City) bool {
//		return city.Name == "London"
//	}
//
// DFS is non-recursive and maintains a stack instead.
func (g *Graph[K, T]) DFS(start T, visit func(value T) bool) error {
	startHash := g.hash(start)

	return g.DFSByHash(startHash, visit)
}

// DFSByHash does the same as DFS, but uses a hash value to identify the starting vertex.
func (g *Graph[K, T]) DFSByHash(startHash K, visit func(value T) bool) error {
	if _, ok := g.vertices[startHash]; !ok {
		return fmt.Errorf("could not find start vertex with hash %v", startHash)
	}

	stack := make([]K, 0)
	visited := make(map[K]bool)

	stack = append(stack, startHash)

	for len(stack) > 0 {
		currentHash := stack[len(stack)-1]
		currentVertex := g.vertices[currentHash]

		stack = stack[:len(stack)-1]

		if _, ok := visited[currentHash]; !ok {
			// Stop traversing the graph if the visit function returns true.
			if visit(currentVertex) {
				break
			}
			visited[currentHash] = true

			for adjacency := range g.outEdges[currentHash] {
				stack = append(stack, adjacency)
			}
		}
	}

	return nil
}

// BFS performs a Breadth-First Search on the graph, starting from the given vertex. The visit
// function will be invoked for each visited vertex. If it returns false, BFS will continue
// traversing the path, and if it returns true, the traversal will be stopped.
//
// This example prints all vertices of the graph in BFS-order:
//
//	g := graph.New(graph.IntHash)
//
//	g.Vertex(1)
//	g.Vertex(2)
//	g.Vertex(3)
//
//	_ = g.Edge(1, 2)
//	_ = g.Edge(2, 3)
//	_ = g.Edge(3, 1)
//
//	_ = g.BFS(1, func(value int) bool {
//		fmt.Println(value)
//		return false
//	})
//
// Similarily, if you have a graph of City vertices and the traversal should stop at London,
// the visit function would look as follows:
//
//	func(city City) bool {
//		return city.Name == "London"
//	}
//
// BFS is non-recursive and maintains a stack instead.
func (g *Graph[K, T]) BFS(start T, visit func(value T) bool) error {
	startHash := g.hash(start)

	return g.BFSByHash(startHash, visit)
}

// BFSByHash does the same as BFS, but uses a hash value to identify the starting vertex.
func (g *Graph[K, T]) BFSByHash(startHash K, visit func(value T) bool) error {
	if _, ok := g.vertices[startHash]; !ok {
		return fmt.Errorf("could not find start vertex with hash %v", startHash)
	}

	queue := make([]K, 0)
	visited := make(map[K]bool)

	visited[startHash] = true
	queue = append(queue, startHash)

	for len(queue) > 0 {
		currentHash := queue[0]
		currentVertex := g.vertices[currentHash]

		queue = queue[1:]

		// Stop traversing the graph if the visit function returns true.
		if visit(currentVertex) {
			break
		}

		for adjacency := range g.outEdges[currentHash] {
			if _, ok := visited[adjacency]; !ok {
				visited[adjacency] = true
				queue = append(queue, adjacency)
			}
		}

	}

	return nil
}

// edgesAreEqual checks two given edges for equality. Two edges are considered equal if their
// source and target vertices are the same or, if the graph is undirected, the same but swapped.
func (g *Graph[K, T]) edgesAreEqual(a, b Edge[T]) bool {
	aSourceHash := g.hash(a.Source)
	aTargetHash := g.hash(a.Target)
	bSourceHash := g.hash(b.Source)
	bTargetHash := g.hash(b.Target)

	if aSourceHash == bSourceHash && aTargetHash == bTargetHash {
		return true
	}

	if !g.properties.isDirected {
		return aSourceHash == bTargetHash && aTargetHash == bSourceHash
	}

	return false
}

func (g *Graph[K, T]) addEdge(sourceHash, targetHash K, edge Edge[T]) {
	if _, ok := g.edges[sourceHash]; !ok {
		g.edges[sourceHash] = make(map[K]Edge[T])
	}

	g.edges[sourceHash][targetHash] = edge

	if _, ok := g.outEdges[sourceHash]; !ok {
		g.outEdges[sourceHash] = make(map[K]Edge[T])
	}

	g.outEdges[sourceHash][targetHash] = edge

	if _, ok := g.inEdges[targetHash]; !ok {
		g.inEdges[targetHash] = make(map[K]Edge[T])
	}

	g.inEdges[targetHash][sourceHash] = edge
}

package graph

// Graph represents a generic graph data structure consisting of vertices and edges. Its vertices
// are of type T, and each vertex is identified by a hash of type K.
type Graph[K comparable, T any] interface {

	// Vertex creates a new vertex in the graph, which won't be connected to another vertex yet.
	// This function is idempotent, but overwrites an existing vertex if the hash already exists.
	Vertex(value T)

	// Edge creates an edge between the source and the target vertex. If the Directed option has
	// been called on the graph, this is a directed edge. Returns an error if either vertex doesn't
	// exist or the edge already exists.
	Edge(source, target T) error

	// EdgeByHashes creates an edge between the source and the target vertex, but uses hash values
	// to identify the vertices. This is convenient when you don't have the full vertex objects at
	// hand. Returns an error if either vertex doesn't exist or the edge already exists.
	//
	// To obtain the hash value for a vertex, call the hashing function passed to New.
	EdgeByHashes(sourceHash, targetHash K) error

	// WeightedEdge does the same as Edge, but adds an additional weight to the created edge. In an
	// unweighted graph, all edges have the same weight of 0.
	WeightedEdge(source, target T, weight int) error

	// WeightedEdgeByHashes does the same as EdgeByHashes, but adds an additional weight to the
	// created edge. In an unweighted graph, all edges have the same weight of 0.
	WeightedEdgeByHashes(sourceHash, targetHash K, weight int) error

	// GetEdgeByHashes returns the edge between two vertices. The second return value indicates
	// whether the edge exists. If the graph  is undirected, an edge with swapped source and target
	// vertices does match.
	GetEdge(source, target T) (Edge[T], bool)

	// GetEdgeByHashes returns the edge between two vertices with the given hash values. The second
	// return value indicates whether the edge exists. If the graph  is undirected, an edge with
	// swapped source and target vertices does match.
	GetEdgeByHashes(sourceHash, targetHash K) (Edge[T], bool)

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
	DFS(start T, visit func(value T) bool) error

	// DFSByHash does the same as DFS, but uses a hash value to identify the starting vertex.
	DFSByHash(startHash K, visit func(value T) bool) error

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
	BFS(start T, visit func(value T) bool) error

	// BFSByHash does the same as BFS, but uses a hash value to identify the starting vertex.
	BFSByHash(startHash K, visit func(value T) bool) error

	// CreatesCycle determines whether an edge between the given source and target vertices would
	// introduce a cycle. It won't create that edge in any case.
	//
	// A potential edge would create a cycle if the target vertex is also a parent of the source
	// vertex. Given a graoh A-B-C-D, adding an edge DA would introduce a cycle:
	//
	//	A -
	//	|  |
	//	B  |
	//	|  |
	//	C  |
	//	|  |
	//	D -
	//
	// CreatesCycle backtracks the ingoing edges of D, resulting in a reverse walk C-B-A.
	CreatesCycle(source, target T) (bool, error)

	// CreatesCycleByHashes does the same as CreatesCycle, but uses a hash value to identify the
	// starting vertex.
	CreatesCycleByHashes(sourceHash, targetHash K) (bool, error)

	// Degree determines and returns the degree of a given vertex.
	Degree(vertex T) (int, error)

	// DegreeByHash does the same as Degree, but uses a hash value to identify the vertex.
	DegreeByHash(vertexHash K) (int, error)
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
// The obtained Graph implementation is depends on these properties.
func New[K comparable, T any](hash Hash[K, T], options ...func(*properties)) Graph[K, T] {
	var p properties

	for _, option := range options {
		option(&p)
	}

	if p.isDirected {
		return newDirected(hash, &p)
	}

	return newUndirected(hash, &p)
}

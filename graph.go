package graph

// Graph represents a generic graph data structure consisting of vertices and edges. Its vertices
// are of type T, and each vertex is identified by a hash of type K.
type Graph[K comparable, T any] interface {

	// Traits returns the graph's traits. Those traits must be set when creating a graph using New.
	Traits() *Traits

	// Vertex creates a new vertex in the graph, which won't be connected to another vertex yet.
	// This function is idempotent, but overwrites an existing vertex if the hash already exists.
	Vertex(value T)

	// Edge creates an edge between the source and the target vertex. If the Directed option has
	// been called on the graph, this is a directed edge. Returns an error if either vertex doesn't
	// exist or the edge already exists.
	//
	// Edge accepts a variety of functional options to set further edge details such as the weight
	// or an attribute:
	//
	//	_ = graph.Edge("A", "B", graph.EdgeWeight(4), graph.EdgeAttribute("label", "mylabel"))
	//
	Edge(source, target T, options ...func(*edgeProperties)) error

	// EdgeByHashes creates an edge between the source and the target vertex, but uses hash values
	// to identify the vertices. This is convenient when you don't have the full vertex objects at
	// hand. Returns an error if either vertex doesn't exist or the edge already exists.
	//
	// To obtain the hash value for a vertex, call the hashing function passed to New.
	//
	// EdgeByHashes accepts the same functional options as Edge.
	EdgeByHashes(sourceHash, targetHash K, options ...func(*edgeProperties)) error

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
	// traversing the path, and if it returns true, the traversal will be stopped. In case the
	// graph is diconnected, only the vertices joined with the starting vertex will be traversed.
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
	// traversing the path, and if it returns true, the traversal will be stopped. In case the
	// graph is diconnected, only the vertices joined with the starting vertex will be traversed.
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

	// StronglyConnectedComponents detects all strongly connected components within the graph and
	// returns the hash values of the respective vertices for each component. This only works for
	// directed graphs.
	//
	// Note that the current implementation uses recursive function calls.
	StronglyConnectedComponents() ([][]K, error)

	// ShortestPath computes the shortest path between a source and a target vertex using the edge
	// weights and returns the hash values of the vertices forming that path. This search runs in
	// O(|V|+|E|log(|V|)) time.
	//
	// The returned path includes the source and target vertices. If the target cannot be reached
	// from the source vertex, ShortestPath returns an error. If there are multiple shortest paths,
	// an arbitrary one will be returned.
	ShortestPath(source, target T) ([]K, error)

	// ShortestPathByHashes does the same as ShortestPath, but uses hash values to identify the
	// vertices.
	ShortestPathByHashes(sourceHash, targetHash K) ([]K, error)

	// AdjacencyMap computes and returns an adjacency map containing all vertices in the graph.
	//
	// There is an entry for each vertex, and each of those entries is another map whose keys are
	// the hash values of the adjacent vertices. The value is an Edge instance that stores the
	// source and target hash values (these are the same as the map keys) as well as edge metadata.
	//
	// For a graph with edges AB and AC, the adjacency map would look as follows:
	//
	//	map[string]map[string]Edge[string]{
	//		"A": map[string]Edge[string]{
	//			"B": {Source: "A", Target: "B"}
	//			"C": {Source: "A", Target: "B"}
	//		}
	//	}
	//
	// This design makes AdjacencyMap suitable for a wide variety of scenarios and demands.
	AdjacencyMap() map[K]map[K]Edge[K]
}

// Edge represents a graph edge with a source and target vertex as well as a weight, which has the
// same value for all edges in an unweighted graph. Even though the vertices are referred to as
// source and target, whether the graph is directed or not is determined by its traits.
type Edge[T any] struct {
	Source     T
	Target     T
	properties edgeProperties
}

type edgeProperties struct {
	Weight     int
	Attributes map[string]string
}

// Hash is a hashing function that takes a vertex of type T and returns a hash value of type K.
//
// Every graph has a hashing function and uses that function to retrieve the hash values of its
// vertices. You can either use one of the predefined hashing functions, or, if you want to store a
// custom data type, provide your own function:
//
//	cityHash := func(c City) string {
//		return c.Name
//	}
//
// The cityHash function returns the city name as a hash value. The types of T and K, in this case
// City and string, also define the types T and K of the graph.
type Hash[K comparable, T any] func(T) K

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
// All traits of the graph can be set using the predefined functional options. They can be combined
// arbitrarily. This example creates a directed acyclic graph:
//
//	g := graph.New(graph.IntHash, graph.Directed(), graph.Acyclic())
//
// Which Graph implementation will be returned depends on these traits.
func New[K comparable, T any](hash Hash[K, T], options ...func(*Traits)) Graph[K, T] {
	var p Traits

	for _, option := range options {
		option(&p)
	}

	if p.IsDirected {
		return newDirected(hash, &p)
	}

	return newUndirected(hash, &p)
}

// StringHash is a hashing function that accepts a string and uses that exact string as a hash
// value. Using it as Hash will yield a Graph[string, string].
func StringHash(v string) string {
	return v
}

// IntHash is a hashing function that accepts an integer and uses that exact integer as a hash
// value. Using it as Hash will yield a Graph[int, int].
func IntHash(v int) int {
	return v
}

// EdgeWeight returns a function that sets the weight of an edge to the given weight. This is a
// functional option for the Edge and EdgeByHashes methods.
func EdgeWeight(weight int) func(*edgeProperties) {
	return func(e *edgeProperties) {
		e.Weight = weight
	}
}

// EdgeAttribute returns a function that adds the given key-value pair to the attributes of an
// edge. This is a functional option for the Edge and EdgeByHashes methods.
func EdgeAttribute(key, value string) func(*edgeProperties) {
	return func(e *edgeProperties) {
		e.Attributes[key] = value
	}
}

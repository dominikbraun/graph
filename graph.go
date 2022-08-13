package graph

// Graph represents a generic graph data structure consisting of vertices and edges. Its vertices
// are of type T, and each vertex is identified by a hash of type K.
type Graph[K comparable, T any] interface {

	// Traits returns the graph's traits. Those traits must be set when creating a graph using New.
	Traits() *Traits

	// AddVertex creates a new vertex in the graph, which won't be connected to another vertex yet.
	// This function is idempotent, but overwrites an existing vertex if the hash already exists.
	AddVertex(value T)

	// Vertex returns the vertex with the given hash or an error if the vertex doesn't exist.
	Vertex(hash K) (T, error)

	// AddEdge creates an edge between the source and the target vertex. If the Directed option has
	// been called on the graph, this is a directed edge. Returns an error if either vertex doesn't
	// exist or the edge already exists.
	//
	// AddEdge accepts a variety of functional options to set further edge details such as the
	// weight or an attribute:
	//
	//	_ = graph.Edge("A", "B", graph.EdgeWeight(4), graph.EdgeAttribute("label", "my-label"))
	//
	AddEdge(sourceHash, targetHash K, options ...func(*EdgeProperties)) error

	// Edge returns the edge between two vertices. The second return value indicates whether the
	// edge exists. If the graph is undirected, an edge with swapped source and target vertices
	// does match.
	Edge(sourceHash, targetHash K) (Edge[T], bool)

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
	AdjacencyMap() (map[K]map[K]Edge[K], error)

	// Predecessors determines and returns the predecessors of the vertex with the given hash.
	//
	// In a directed graph, these are the directed predecessors with an outgoing edge to the given
	// edge whereas in an undirected graph, all adjacent vertices are considered as predecessors.
	//
	// Consequently, Predecessors returns a subset of AdjacencyMap for the given vertex.
	Predecessors(vertex K) (map[K]Edge[K], error)
}

// Edge represents a graph edge with a source and target vertex as well as a weight, which has the
// same value for all edges in an unweighted graph. Even though the vertices are referred to as
// source and target, whether the graph is directed or not is determined by its traits.
type Edge[T any] struct {
	Source     T
	Target     T
	Properties EdgeProperties
}

// EdgeProperties represents a set of properties that each edge possesses. They can be set when
// adding a new edge using the functional options provided by this library:
//
//	g.Edge("A", "B", graph.EdgeWeight(2), graph.EdgeAttribute("color", "red"))
//
// The example above will create an edge with weight 2 and a "color" attribute with value "red".
type EdgeProperties struct {
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
//	g.AddVertex(1)
//	g.AddVertex(2)
//	g.AddVertex(3)
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
//	g.AddVertex(london)
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
// functional option for the Edge and AddEdge methods.
func EdgeWeight(weight int) func(*EdgeProperties) {
	return func(e *EdgeProperties) {
		e.Weight = weight
	}
}

// EdgeAttribute returns a function that adds the given key-value pair to the attributes of an
// edge. This is a functional option for the Edge and AddEdge methods.
func EdgeAttribute(key, value string) func(*EdgeProperties) {
	return func(e *EdgeProperties) {
		e.Attributes[key] = value
	}
}

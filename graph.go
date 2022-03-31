package graph

// Set is a collection of objects. Typically, these objects are vertices
// or pairs. It is a set in a mathematical sense, however, it actually is
// a map data structure to enable access in constant time.
type Set[T comparable] map[T]T

// Values returns all values in the set as a slice. Because a set merely
// is a map internally, elements returned by Values are in random order.
func (s Set[T]) Values() []T {
	values := make([]T, 0, len(s))

	for _, v := range s {
		values = append(values, v)
	}

	return values
}

// Pair represents an pair of joined vertices. Whether this is a ordered
// or unordered pair is up to the graph implementation using the pair.
type Pair[T comparable] struct {
	A, B T
}

// Graph represents a graph according to the definition G = (V, E), where
// V is a set of vertices and E a set of paired vertices (edges). This is
// a general graph that can be a directed or undirected, cyclic or acyclic
// graph, or a special graph like a rooted tree.
type Graph[T comparable] struct {
	vertices Set[T]
	edges    Set[Pair[T]]
}

// Vertices returns a list of all vertices in the graph.
func (g *Graph[T]) Vertices() []T {
	return g.vertices.Values()
}

// Edges returns a list of all edges, i.e. pairs of vertices, in the graph.
func (g *Graph[T]) Edges() []Pair[T] {
	return g.edges.Values()
}

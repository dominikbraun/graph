package graph

// Set is a collection of objects. Typically, these objects are vertices
// or pairs. It is implemented as a map to enable access in constant time.
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

// Pair represents a pair of joined vertices. Whether this is an ordered
// or unordered pair is up to the graph implementation using the pair.
type Pair[V comparable] struct {
	A, B V
}

// Equals checks wether the pair is equal to another unordered pair.
func (p Pair[V]) Equals(other Pair[V]) bool {
	return p.A == other.A && p.B == other.B ||
		p.A == other.B && p.B == other.A
}

// Graph represents a graph according to the definition G = (V, E), where
// V is a set of vertices and E a set of paired vertices (edges).
//
// Graph is a general graph that can be a directed or undirected, cyclic
// or acyclic graph, or a special graph like a rooted tree. This library
// provides pre-defined structures for some of those special graph types.
type Graph[V comparable] struct {
	vertices Set[V]
	edges    Set[Pair[V]]
}

// Vertices returns a list of all vertices in the graph.
func (g *Graph[V]) Vertices() []V {
	return g.vertices.Values()
}

// Edges returns a list of all edges, i.e. pairs of vertices, in the graph.
func (g *Graph[V]) Edges() []Pair[V] {
	return g.edges.Values()
}

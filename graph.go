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
// called on the graph, this is a directed edge. Returns an error if either vertex doesn't exist.
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
// Returns an error if either vertex doesn't exist.
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

	if _, ok := g.edges[sourceHash]; !ok {
		g.edges[sourceHash] = make(map[K]Edge[T])
	}

	edge := Edge[T]{
		Source: source,
		Target: target,
		Weight: weight,
	}

	g.edges[sourceHash][targetHash] = edge

	return nil
}

// GetEdge returns the edge between two given vertices or an error if an edge doesn't exist. The
// order of the vertices isn't relevant if the graph is undirected.
func (g *Graph[K, T]) GetEdge(source, target T) (Edge[T], error) {
	sourceHash := g.hash(source)
	targetHash := g.hash(target)

	return g.GetEdgeByHashes(sourceHash, targetHash)
}

func (g *Graph[K, T]) GetEdgeByHashes(sourceHash, targetHash K) (Edge[T], error) {
	edgesOfSource, ok := g.edges[sourceHash]
	if !ok && g.properties.isDirected {
		return Edge[T]{}, fmt.Errorf("could not find edge with source %v", sourceHash)
	}

	edge, ok := edgesOfSource[targetHash]
	if !ok && g.properties.isDirected {
		return Edge[T]{}, fmt.Errorf("could not find edge with source %v and target %v", sourceHash, targetHash)
	}

	if ok {
		return edge, nil
	}

	if !g.properties.isDirected {
		edgesOftarget, ok := g.edges[targetHash]
		if !ok {
			return Edge[T]{}, fmt.Errorf("could not find edge with source %v", targetHash)
		}

		edge, ok := edgesOftarget[sourceHash]
		if !ok {
			return Edge[T]{}, fmt.Errorf("could not find edge with source %v and target %v", targetHash, sourceHash)
		}

		return edge, nil
	}

	return Edge[T]{}, fmt.Errorf("could not find edge with source %v and target %v", sourceHash, targetHash)
}

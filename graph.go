package graph

import "fmt"

type Graph[K comparable, T any] struct {
	hash       Hash[K, T]
	properties *Properties
	vertices   map[K]T
	edges      map[K]map[K]Edge[T]
}

type Edge[T any] struct {
	Source T
	Target T
	Weight int
}

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

func (g *Graph[K, T]) Vertex(value T) {
	hash := g.hash(value)
	g.vertices[hash] = value
}

func (g *Graph[K, T]) Edge(source, target T) {
	g.WeightedEdge(source, target, 0)
}

func (g *Graph[K, T]) WeightedEdge(source, target T, weight int) {
	edge := Edge[T]{
		Source: source,
		Target: target,
		Weight: weight,
	}

	sourceHash := g.hash(source)
	targetHash := g.hash(target)

	if _, ok := g.edges[sourceHash]; !ok {
		g.edges[sourceHash] = make(map[K]Edge[T])
	}

	g.edges[sourceHash][targetHash] = edge
}

func (g *Graph[K, T]) EdgeByHashes(sourceHash, targetHash K) error {
	return g.WeightedEdgeByHashes(sourceHash, targetHash, 0)
}

func (g *Graph[K, T]) WeightedEdgeByHashes(sourceHash, targetHash K, weight int) error {
	source, ok := g.vertices[sourceHash]
	if !ok {
		return fmt.Errorf("could not find source vertex with hash %v", source)
	}

	target, ok := g.vertices[targetHash]
	if !ok {
		return fmt.Errorf("could not find target vertex with hash %v", source)
	}

	g.WeightedEdge(source, target, weight)

	return nil
}

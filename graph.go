package graph

type Graph[K comparable, T any] struct {
	hash       Hash[K, T]
	properties *Properties
	vertices   map[K]T
}

func New[K comparable, T any](hash Hash[K, T], options ...func(*Properties)) *Graph[K, T] {
	g := Graph[K, T]{
		hash:       hash,
		properties: &Properties{},
	}

	for _, option := range options {
		option(g.properties)
	}

	return &g
}

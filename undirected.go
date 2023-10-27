package graph

import (
	"fmt"
)

type undirected[K comparable, T any] struct {
	graph[K, T]
}

func newUndirected[K comparable, T any](hash Hash[K, T], traits *Traits, store Store[K, T]) *undirected[K, T] {
	return &undirected[K, T]{
		graph: graph[K, T]{
			hash:   hash,
			traits: traits,
			store:  store,
		},
	}
}

func (u *undirected[K, T]) AddEdgesFrom(g Graph[K, T]) error {
	other := tograph(g)
	return u.addEdgesFrom(&other)
}

func (u *undirected[K, T]) AddVerticesFrom(g Graph[K, T]) error {
	other := tograph(g)
	return u.addVerticesFrom(&other)
}

func (u *undirected[K, T]) Clone() (Graph[K, T], error) {
	gclone, err := u.clone()
	if err != nil {
		return nil, err
	}

	dir := &undirected[K, T]{
		graph: *gclone,
	}
	return dir, nil
}

func (u *undirected[K, T]) UpdateEdge(source, target K, options ...func(properties *EdgeProperties)) error {
	existingEdge, err := u.updateEdge(source, target, options...)
	if err != nil {
		return err
	}

	reversedEdge := existingEdge
	reversedEdge.Source = existingEdge.Target
	reversedEdge.Target = existingEdge.Source

	return u.store.UpdateEdge(target, source, reversedEdge)
}

func (u *undirected[K, T]) RemoveEdge(source, target K) error {
	err := u.graph.RemoveEdge(source, target)
	if err != nil {
		return err
	}

	// remove reciprocal edge
	if err = u.store.RemoveEdge(target, source); err != nil {
		return fmt.Errorf("failed to remove edge from %v to %v: %w", target, source, err)
	}

	return nil
}

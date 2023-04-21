package graph

import (
	"errors"
	"fmt"
	"sort"
)

func MinimumSpanningTree[K comparable, T any](g Graph[K, T]) (Graph[K, T], error) {
	if g.Traits().IsDirected {
		return nil, errors.New("spanning trees can only be determined for undirected graphs")
	}

	adjacencyMap, err := g.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("failed to get adjacency map: %w", err)
	}

	edges := make([]Edge[K], 0)
	subtrees := newUnionFind[K]()

	mst := New(g.(*undirected[K, T]).hash)

	for v, adjacencies := range adjacencyMap {
		subtrees.add(v)

		vertex, properties, err := g.VertexWithProperties(v)
		if err != nil {
			return nil, fmt.Errorf("failed to get vertex %v: %w", v, err)
		}

		err = mst.AddVertex(vertex, copyVertexProperties(properties))
		if err != nil {
			return nil, fmt.Errorf("failed to add vertex %v: %w", v, err)
		}

		for _, edge := range adjacencies {
			edges = append(edges, edge)
		}
	}

	sort.Slice(edges, func(i, j int) bool {
		return edges[i].Properties.Weight < edges[j].Properties.Weight
	})

	for _, edge := range edges {
		sourceRoot := subtrees.find(edge.Source)
		targetRoot := subtrees.find(edge.Target)

		if sourceRoot != targetRoot {
			subtrees.union(sourceRoot, targetRoot)

			err := mst.AddEdge(edge.Source, edge.Target, copyEdgeProperties(edge.Properties))
			if err != nil {
				return nil, fmt.Errorf("failed to add edge (%v, %v): %w", edge.Source, edge.Target, err)
			}
		}
	}

	return mst, nil
}

// unionFind implements a union-find or disjoint set data structure that works
// with vertex hashes as vertices. It's an internal helper type at the moment,
// but could perhaps be exposed publicly in the future.
type unionFind[K comparable] struct {
	parents map[K]K
}

func newUnionFind[K comparable](vertices ...K) *unionFind[K] {
	u := &unionFind[K]{
		parents: make(map[K]K),
	}

	for _, vertex := range vertices {
		u.parents[vertex] = vertex
	}

	return u
}

func (u *unionFind[K]) add(vertex K) {
	u.parents[vertex] = vertex
}

func (u *unionFind[K]) union(vertex1, vertex2 K) {
	root1 := u.find(vertex1)
	root2 := u.find(vertex2)

	if root1 == root2 {
		return
	}

	u.parents[root2] = root1
}

func (u *unionFind[K]) find(vertex K) K {
	root := vertex

	for u.parents[root] != root {
		root = u.parents[root]
	}

	// Perform a path compression in order to optimize of future find calls.
	current := vertex

	for u.parents[current] != root {
		parent := u.parents[vertex]
		u.parents[vertex] = root
		current = parent
	}

	return root
}

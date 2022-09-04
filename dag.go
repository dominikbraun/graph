package graph

import (
	"errors"
	"fmt"
)

// TopologicalSort performs a topological sort on a given graph and returns the vertex hashes in
// topological order. A topological order is a non-unique order of the vertices in a directed graph
// where an edge from vertex A to vertex B implies that vertex A appears before vertex B.
//
// TopologicalSort only works for directed acyclic graphs. The current implementation works non-
// recursively and uses Kahn's algorithm.
func TopologicalSort[K comparable, T any](g Graph[K, T]) ([]K, error) {
	if !g.Traits().IsDirected || !g.Traits().IsAcyclic {
		return nil, errors.New("topological sort can only be performed on DAGs")
	}

	predecessorMap, err := g.PredecessorMap()
	if err != nil {
		return nil, fmt.Errorf("failed to get predecessor map: %w", err)
	}

	queue := make([]K, 0)

	for vertex, predecessors := range predecessorMap {
		if len(predecessors) == 0 {
			queue = append(queue, vertex)
		}
	}

	order := make([]K, 0, len(predecessorMap))
	visited := make(map[K]struct{})

	for len(queue) > 0 {
		currentVertex := queue[0]
		queue = queue[1:]

		if _, ok := visited[currentVertex]; ok {
			continue
		}

		order = append(order, currentVertex)
		visited[currentVertex] = struct{}{}

		for vertex, predecessors := range predecessorMap {
			delete(predecessors, currentVertex)

			if len(predecessors) == 0 {
				queue = append(queue, vertex)
			}
		}
	}

	return order, nil
}

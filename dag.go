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
	if !isDAG(g) {
		return nil, errors.New("topological sort can only be performed on DAGs created with the PermitCycles option")
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

// TransitiveReduction transforms the graph into its own transitive reduction. The transitive
// reduction of the given graph is another graph with the same vertices and the same reachability,
// but with as few edges as possible. This greatly reduces the complexity of the graph.
//
// With a time complexity of O(V(V+E)), TransitiveReduction is a very costly operation.
func TransitiveReduction[K comparable, T any](g Graph[K, T]) error {
	if !isDAG(g) {
		return errors.New("topological sort can only be performed on DAGs created with the PermitCycles option")
	}

	adjacencyMap, err := g.AdjacencyMap()
	if err != nil {
		return fmt.Errorf("failed to get adajcency map: %w", err)
	}

	for vertex, successors := range adjacencyMap {
		// For each direct successor of the current vertex, run a DFS starting from that successor.
		// Then, for each vertex visited in the DFS, inspect all of its edges. Remove the edges that
		// also appear in the edges of the top-level iteration vertex.
		//
		// These edges are redundant because their targets obviously are reachable via DFS, i.e.
		// they aren't needed in the top-level vertex anymore and can be removed from there.
		for successor := range successors {
			err := DFS(g, successor, func(current K) bool {
				for _, edge := range adjacencyMap[current] {
					if _, ok := adjacencyMap[vertex][edge.Target]; ok {
						_ = g.RemoveEdge(vertex, edge.Target)
					}
				}

				return false
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func isDAG[K comparable, T any](g Graph[K, T]) bool {
	return g.Traits().IsDirected && g.Traits().PermitCycles
}

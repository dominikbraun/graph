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
	if !g.Traits().IsDirected {
		return nil, fmt.Errorf("topological sort cannot be computed on undirected graph")
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

	if len(order) != g.Order() {
		return nil, errors.New("topological sort cannot be computed on graph with cycles")
	}

	return order, nil
}

// TransitiveReduction returns another graph with the same vertices and the same reachability, but
// with as few edges as possible. This greatly reduces the complexity of the graph.
//
// With a time complexity of O(V(V+E)), TransitiveReduction is a very costly operation.
func TransitiveReduction[K comparable, T any](g Graph[K, T]) (Graph[K, T], error) {
	if !g.Traits().IsDirected {
		return nil, fmt.Errorf("transitive reduction cannot be performed on undirected graph")
	}

	transitiveReduction, err := g.Clone()
	if err != nil {
		return nil, fmt.Errorf("failed to clone the graph: %w", err)
	}

	adjacencyMap, err := transitiveReduction.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("failed to get adajcency map: %w", err)
	}

	for vertex, successors := range adjacencyMap {
		// For each direct successor of the current vertex, run a DFS starting from that successor.
		// Then, for each vertex visited in the DFS, inspect all of its edges. Remove the edges that
		// also appear in the edges of the top-level iteration vertex.
		//
		// These edges are redundant because their targets obviously are reachable via DFS, i.e.
		// they aren't needed in the top-level vertex anymore and can be removed from there.
		for successor := range successors {
			visited := make(map[K]struct{}, transitiveReduction.Order())
			onStack := make(map[K]struct{}, transitiveReduction.Order())
			stack := append(make([]K, 0, transitiveReduction.Order()), successor)

			for len(stack) > 0 {
				current := stack[len(stack)-1]
				stack = stack[:len(stack)-1]

				if _, ok := visited[current]; !ok {
					// If the vertex is not yet visited, mark it as visited and put it on the stack.
					visited[current] = struct{}{}
					onStack[current] = struct{}{}
				} else {
					// Otherwise, remove the vertex from the stack.
					delete(onStack, current)
					continue
				}

				// If the vertex is a leaf node, remove it from the stack.
				if len(adjacencyMap[current]) == 0 {
					delete(onStack, current)
				}

				for adjacency := range adjacencyMap[current] {
					if _, ok := visited[adjacency]; ok {
						if _, ok := onStack[adjacency]; ok {
							// If this vertex is visited as well as on the stack, we have a cycle.
							return nil, fmt.Errorf("transitive reduction cannot be performed on graph with cycle")
						}
						continue
					}

					if _, ok := adjacencyMap[vertex][adjacency]; ok {
						_ = transitiveReduction.RemoveEdge(vertex, adjacency)
					}
					stack = append(stack, adjacency)
				}
			}
		}
	}

	return transitiveReduction, nil
}

package graph

import "fmt"

// CreatesCycle determines whether an edge between the given source and target vertices would
// introduce a cycle. It won't create that edge in any case.
//
// A potential edge would create a cycle if the target vertex is also a parent of the source vertex.
// Given a graph A-B-C-D, adding an edge DA would introduce a cycle:
//
//	A -
//	|  |
//	B  |
//	|  |
//	C  |
//	|  |
//	D -
//
// CreatesCycle backtracks the ingoing edges of D, resulting in a reverse walk C-B-A.
func CreatesCycle[K comparable, T any](g Graph[K, T], source, target K) (bool, error) {
	if _, err := g.Vertex(source); err != nil {
		return false, fmt.Errorf("could not get vertex with hash %v: %w", source, err)
	}

	if _, err := g.Vertex(target); err != nil {
		return false, fmt.Errorf("could not get vertex with hash %v: %w", target, err)
	}

	if source == target {
		return true, nil
	}

	stack := make([]K, 0)
	visited := make(map[K]bool)

	stack = append(stack, source)

	for len(stack) > 0 {
		currentHash := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if _, ok := visited[currentHash]; !ok {
			// If the current vertex, e.g. an adjacency of the source vertex, also is the target
			// vertex, an edge between these two would create a cycle.
			if currentHash == target {
				return true, nil
			}
			visited[currentHash] = true

			predecessors, _ := g.Predecessors(currentHash)

			for adjacency := range predecessors {
				stack = append(stack, adjacency)
			}
		}
	}

	return false, nil
}

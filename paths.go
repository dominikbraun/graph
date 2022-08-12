package graph

import (
	"fmt"
	"math"
)

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

// ShortestPath computes the shortest path between a source and a target vertex using the edge
// weights and returns the hash values of the vertices forming that path. This search runs in
// O(|V|+|E|log(|V|)) time.
//
// The returned path includes the source and target vertices. If the target cannot be reached
// from the source vertex, ShortestPath returns an error. If there are multiple shortest paths,
// an arbitrary one will be returned.
func ShortestPath[K comparable, T any](g Graph[K, T], source, target K) ([]K, error) {
	weights := make(map[K]float64)
	visited := make(map[K]bool)
	predecessors := make(map[K]K)

	weights[source] = 0
	visited[target] = true

	queue := newPriorityQueue[K]()
	adjacencyMap := g.AdjacencyMap()

	for hash := range adjacencyMap {
		if hash != source {
			weights[hash] = math.Inf(1)
			visited[hash] = false
		}

		queue.Push(hash, weights[hash])
	}

	for queue.Len() > 0 {
		vertex, _ := queue.Pop()
		hasInfiniteWeight := math.IsInf(weights[vertex], 1)

		if vertex == target {
			targetPredecessors, err := g.Predecessors(target)
			if err != nil {
				return nil, fmt.Errorf("failed to get predecessors of %v: %w", target, err)
			}
			if len(targetPredecessors) == 0 {
				return nil, fmt.Errorf("vertex %v is not reachable from vertex %v", target, source)
			}
		}

		for adjacency, edge := range adjacencyMap[vertex] {
			weight := weights[vertex] + float64(edge.Properties.Weight)

			if weight < weights[adjacency] && !hasInfiniteWeight {
				weights[adjacency] = weight
				predecessors[adjacency] = vertex
				queue.DecreasePriority(adjacency, weight)
			}
		}
	}

	// Backtrack the predecessors from target to source. These are the least-weighted edges.
	path := []K{target}
	hashCursor := target

	for hashCursor != source {
		hashCursor = predecessors[hashCursor]
		path = append([]K{hashCursor}, path...)
	}

	return path, nil
}

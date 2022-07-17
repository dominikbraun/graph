package graph

import (
	"errors"
	"fmt"
	"math"
)

type undirected[K comparable, T any] struct {
	hash   Hash[K, T]
	traits *traits
	store  Store[K, T]
}

func newUndirected[K comparable, T any](
	hash Hash[K, T],
	store Store[K, T],
	traits *traits,
) *undirected[K, T] {
	return &undirected[K, T]{
		hash:   hash,
		traits: traits,
		store:  store,
	}
}

func (u *undirected[K, T]) Vertex(value T) {
	u.store.AddVertex(value) // TODO: error
}

func (u *undirected[K, T]) Edge(source, target T) error {
	return u.WeightedEdge(source, target, 0)
}

func (u *undirected[K, T]) WeightedEdge(source, target T, weight int) error {
	sourceHash := u.hash(source)
	targetHash := u.hash(target)

	return u.WeightedEdgeByHashes(sourceHash, targetHash, weight)
}

func (u *undirected[K, T]) EdgeByHashes(sourceHash, targetHash K) error {
	return u.WeightedEdgeByHashes(sourceHash, targetHash, 0)
}

func (u *undirected[K, T]) WeightedEdgeByHashes(sourceHash, targetHash K, weight int) error {
	_, ok := u.store.GetVertex(sourceHash)
	if !ok {
		return fmt.Errorf("could not find source vertex with hash %v", sourceHash)
	}

	_, ok = u.store.GetVertex(targetHash)
	if !ok {
		return fmt.Errorf("could not find target vertex with hash %v", targetHash)
	}

	if _, ok := u.GetEdgeByHashes(sourceHash, targetHash); ok {
		return fmt.Errorf("an edge between vertices %v and %v already exists", sourceHash, targetHash)
	}

	// If the graph was declared to be acyclic, permit the creation of a cycle.
	if u.traits.isAcyclic {
		createsCycle, err := u.CreatesCycleByHashes(sourceHash, targetHash)
		if err != nil {
			return fmt.Errorf("failed to check for cycles: %w", err)
		}
		if createsCycle {
			return fmt.Errorf("an edge between %v and %v would introduce a cycle", sourceHash, targetHash)
		}
	}

	edge := Edge[K]{
		Source: sourceHash,
		Target: targetHash,
		Weight: weight,
	}

	// Note(geoah): note sure about this
	// In an undirected graph, since multigraphs aren't supported, the edge AB is the same as BA.
	// To allow for this, while keeping the API friendly, we add the edge both directions.
	u.store.AddEdge(sourceHash, targetHash, edge)
	u.store.AddEdge(targetHash, sourceHash, edge)

	return nil
}

func (u *undirected[K, T]) GetEdge(source, target T) (Edge[K], bool) {
	sourceHash := u.hash(source)
	targetHash := u.hash(target)

	return u.GetEdgeByHashes(sourceHash, targetHash)
}

func (u *undirected[K, T]) GetEdgeByHashes(sourceHash, targetHash K) (Edge[K], bool) {
	edge, ok := u.store.GetEdge(sourceHash, targetHash)
	if !ok {
		return Edge[K]{}, false
	}

	return edge, ok
}

func (u *undirected[K, T]) DFS(start T, visit func(value T) bool) error {
	startHash := u.hash(start)

	return u.DFSByHash(startHash, visit)
}

func (u *undirected[K, T]) DFSByHash(startHash K, visit func(value T) bool) error {
	if _, ok := u.store.GetVertex(startHash); !ok {
		return fmt.Errorf("could not find start vertex with hash %v", startHash)
	}

	stack := make([]K, 0)
	visited := make(map[K]bool)

	stack = append(stack, startHash)

	for len(stack) > 0 {
		currentHash := stack[len(stack)-1]
		currentVertex, _ := u.store.GetVertex(currentHash)

		stack = stack[:len(stack)-1]

		if _, ok := visited[currentHash]; !ok {
			// Stop traversing the graph if the visit function returns true.
			if visit(*currentVertex) {
				break
			}
			visited[currentHash] = true

			targetHashes, _ := u.store.GetEdgeTargetHashes(currentHash) // TODO: error
			stack = append(stack, targetHashes...)
		}
	}

	return nil
}

func (u *undirected[K, T]) BFS(start T, visit func(value T) bool) error {
	startHash := u.hash(start)

	return u.BFSByHash(startHash, visit)
}

func (u *undirected[K, T]) BFSByHash(startHash K, visit func(value T) bool) error {
	if _, ok := u.store.GetVertex(startHash); !ok {
		return fmt.Errorf("could not find start vertex with hash %v", startHash)
	}

	queue := make([]K, 0)
	visited := make(map[K]bool)

	visited[startHash] = true
	queue = append(queue, startHash)

	for len(queue) > 0 {
		currentHash := queue[0]
		currentVertex, _ := u.store.GetVertex(currentHash)

		queue = queue[1:]

		// Stop traversing the graph if the visit function returns true.
		if visit(*currentVertex) {
			break
		}

		targetHashes, _ := u.store.GetEdgeTargetHashes(currentHash) // TODO: error
		for _, adjacency := range targetHashes {
			if _, ok := visited[adjacency]; !ok {
				visited[adjacency] = true
				queue = append(queue, adjacency)
			}
		}

	}

	return nil
}

func (u *undirected[K, T]) CreatesCycle(source, target T) (bool, error) {
	sourceHash := u.hash(source)
	targetHash := u.hash(target)

	return u.CreatesCycleByHashes(sourceHash, targetHash)
}

func (u *undirected[K, T]) CreatesCycleByHashes(sourceHash, targetHash K) (bool, error) {
	_, ok := u.store.GetVertex(sourceHash)
	if !ok {
		return false, fmt.Errorf("could not find source vertex with hash %v", sourceHash)
	}

	_, ok = u.store.GetVertex(targetHash)
	if !ok {
		return false, fmt.Errorf("could not find target vertex with hash %v", targetHash)
	}

	if sourceHash == targetHash {
		return true, nil
	}

	stack := make([]K, 0)
	visited := make(map[K]bool)

	stack = append(stack, sourceHash)

	for len(stack) > 0 {
		currentHash := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if _, ok := visited[currentHash]; !ok {
			// If the current vertex, e.g. a predecessor of the source vertex, is also the target
			// vertex, an edge between these two would create a cycle.
			if currentHash == targetHash {
				return true, nil
			}
			visited[currentHash] = true

			for _, adjacency := range u.adjacencies(currentHash) {
				stack = append(stack, adjacency)
			}
		}
	}

	return false, nil
}

func (u *undirected[K, T]) Degree(vertex T) (int, error) {
	sourceHash := u.hash(vertex)

	return u.DegreeByHash(sourceHash)
}

func (u *undirected[K, T]) DegreeByHash(vertexHash K) (int, error) {
	if _, ok := u.store.GetVertex(vertexHash); !ok {
		return 0, fmt.Errorf("could not find vertex with hash %v", vertexHash)
	}

	degree := 0

	// Adding the number of ingoing edges is sufficient for undirected graphs, because all edges
	// exist twice (as two directed edges in opposite directions). Either dividing the number of
	// ingoing + outgoing edges by 2 or just using the number of ingoing edges is appropriate.
	if inEdges, ok := u.store.GetEdgeSourceHashes(vertexHash); ok {
		degree += len(inEdges)
	}

	return degree, nil
}

func (u *undirected[K, T]) StronglyConnectedComponents() ([][]K, error) {
	return nil, errors.New("strongly connected components may only be detected in directed graphs")
}

// ShortestPath computes the shortest path between two vertices and returns the hashes of the
// vertices forming that path. The current implementation uses Dijkstra with a priority queue.
func (u *undirected[K, T]) ShortestPath(source, target T) ([]K, error) {
	sourceHash := u.hash(source)
	targetHash := u.hash(target)

	return u.ShortestPathByHashes(sourceHash, targetHash)
}

func (u *undirected[K, T]) ShortestPathByHashes(sourceHash, targetHash K) ([]K, error) {
	weights := make(map[K]float64)
	visited := make(map[K]bool)
	predecessors := make(map[K]K)

	weights[sourceHash] = 0
	visited[sourceHash] = true

	queue := newPriorityQueue[K]()

	hashes, _ := u.store.GetAllVertexHashes() // TODO: error
	for _, hash := range hashes {
		if hash != sourceHash {
			weights[hash] = math.Inf(1)
			visited[hash] = false
		}

		queue.Push(hash, weights[hash])
	}

	for queue.Len() > 0 {
		vertex, _ := queue.Pop()
		hasInfiniteWeight := math.IsInf(float64(weights[vertex]), 1)

		if vertex == targetHash {
			if _, ok := u.store.GetEdgeSourceHashes(vertex); !ok {
				return nil, fmt.Errorf("vertex %v is not reachable from vertex %v", targetHash, sourceHash)
			}
		}

		inEdges, ok := u.store.GetEdgeSources(vertex)
		if !ok {
			continue
		}

		for _, edge := range inEdges {
			successor := edge.Target
			weight := weights[vertex] + float64(edge.Weight)

			if weight < weights[successor] && !hasInfiniteWeight {
				weights[successor] = weight
				predecessors[successor] = vertex
				queue.DecreasePriority(successor, weight)
			}
		}
	}

	// Backtrack the predecessors from target to source. These are the least-weighted edges.
	path := []K{targetHash}
	hashCursor := targetHash

	for hashCursor != sourceHash {
		hashCursor = predecessors[hashCursor]
		path = append([]K{hashCursor}, path...)
	}

	return path, nil
}

func (u *undirected[K, T]) edgesAreEqual(a, b Edge[K]) bool {
	if a.Source == b.Source && a.Target == b.Target {
		return true
	}

	if !u.traits.isDirected {
		return a.Source == b.Target && a.Target == b.Source
	}

	return false
}

// func (u *undirected[K, T]) addEdge(sourceHash, targetHash K, edge Edge[K]) {
// 	if _, ok := u.edges[sourceHash]; !ok {
// 		u.edges[sourceHash] = make(map[K]Edge[K])
// 	}
// 	if _, ok := u.edges[targetHash]; !ok {
// 		u.edges[targetHash] = make(map[K]Edge[K])
// 	}

// 	u.edges[sourceHash][targetHash] = edge
// 	u.edges[targetHash][sourceHash] = edge

// 	if _, ok := u.outEdges[sourceHash]; !ok {
// 		u.outEdges[sourceHash] = make(map[K]Edge[K])
// 	}
// 	if _, ok := u.outEdges[targetHash]; !ok {
// 		u.outEdges[targetHash] = make(map[K]Edge[K])
// 	}

// 	u.outEdges[sourceHash][targetHash] = edge
// 	u.outEdges[targetHash][sourceHash] = edge

// 	if _, ok := u.inEdges[targetHash]; !ok {
// 		u.inEdges[targetHash] = make(map[K]Edge[K])
// 	}
// 	if _, ok := u.inEdges[sourceHash]; !ok {
// 		u.inEdges[sourceHash] = make(map[K]Edge[K])
// 	}

// 	u.inEdges[targetHash][sourceHash] = edge
// 	u.inEdges[sourceHash][targetHash] = edge
// }

func (u *undirected[K, T]) adjacencies(vertexHash K) []K {
	var adjacencyHashes []K

	// An undirected graph creates an undirected edge as two directed edges in the opposite
	// direction, so both the in-edges and the out-edges work here.
	targetHashes, ok := u.store.GetEdgeTargetHashes(vertexHash)
	if !ok {
		return adjacencyHashes
	}

	adjacencyHashes = append(adjacencyHashes, targetHashes...)

	return adjacencyHashes
}

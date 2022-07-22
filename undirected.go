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

func (u *undirected[K, T]) Vertex(value T) error {
	hash := u.hash(value)
	err := u.store.AddVertex(hash, value)
	if err != nil {
		return fmt.Errorf("could not get vertex with hash %v, %w", hash, err)
	}
	return nil
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
	if _, err := u.store.GetVertex(sourceHash); err != nil {
		return fmt.Errorf("could not find source vertex with hash %v, %w", sourceHash, err)
	}

	if _, err := u.store.GetVertex(targetHash); err != nil {
		return fmt.Errorf("could not find target vertex with hash %v, %w", targetHash, err)
	}

	if _, err := u.GetEdgeByHashes(sourceHash, targetHash); err == nil {
		return fmt.Errorf("an edge between vertices %v and %v already exists, %w", sourceHash, targetHash, err)
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

	edgeAB := Edge[K]{
		Source: sourceHash,
		Target: targetHash,
		Weight: weight,
	}

	edgeBA := Edge[K]{
		Source: targetHash,
		Target: sourceHash,
		Weight: weight,
	}

	// Note(geoah): note sure about this
	// In an undirected graph, since multigraphs aren't supported, the edge AB is the same as BA.
	// To allow for this, while keeping the API friendly, we add the edge both directions.
	u.store.AddEdge(sourceHash, targetHash, edgeAB)
	u.store.AddEdge(targetHash, sourceHash, edgeBA)

	return nil
}

func (u *undirected[K, T]) GetEdge(source, target T) (*Edge[K], error) {
	sourceHash := u.hash(source)
	targetHash := u.hash(target)

	return u.GetEdgeByHashes(sourceHash, targetHash)
}

func (u *undirected[K, T]) GetEdgeByHashes(sourceHash, targetHash K) (*Edge[K], error) {
	edge, err := u.store.GetEdge(sourceHash, targetHash)
	if err != nil {
		return nil, fmt.Errorf("could not find edge with hashes %v %v, %w", sourceHash, targetHash, err)
	}

	return edge, nil
}

func (u *undirected[K, T]) DFS(start T, visit func(value T) bool) error {
	startHash := u.hash(start)

	return u.DFSByHash(startHash, visit)
}

func (u *undirected[K, T]) DFSByHash(startHash K, visit func(value T) bool) error {
	if _, err := u.store.GetVertex(startHash); err != nil {
		return fmt.Errorf("could not find start vertex with hash %v, %w", startHash, err)
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

			edges, err := u.store.GetEdgesBySource(currentHash)
			if err != nil && !errors.Is(err, ErrNotFound) {
				return fmt.Errorf("could not get edges by source with hash %v, %w", currentHash, err)
			}
			for _, edge := range edges {
				stack = append(stack, edge.Target)
			}
		}
	}

	return nil
}

func (u *undirected[K, T]) BFS(start T, visit func(value T) bool) error {
	startHash := u.hash(start)

	return u.BFSByHash(startHash, visit)
}

func (u *undirected[K, T]) BFSByHash(startHash K, visit func(value T) bool) error {
	if _, err := u.store.GetVertex(startHash); err != nil {
		return fmt.Errorf("could not find start vertex with hash %v, %w", startHash, err)
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

		edges, _ := u.store.GetEdgesBySource(currentHash) // TODO: error
		for _, edge := range edges {
			if _, ok := visited[edge.Target]; !ok {
				visited[edge.Target] = true
				queue = append(queue, edge.Target)
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
	if _, err := u.store.GetVertex(sourceHash); err != nil {
		return false, fmt.Errorf("could not find source vertex with hash %v, %w", sourceHash, err)
	}

	if _, err := u.store.GetVertex(targetHash); err != nil {
		return false, fmt.Errorf("could not find target vertex with hash %v, %w", targetHash, err)
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

			adjacencies, err := u.adjacencies(currentHash)
			if err != nil {
				return false, fmt.Errorf("could not get adjacencies with hash %v, %w", currentHash, err)
			}
			for _, adjacency := range adjacencies {
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
	if _, err := u.store.GetVertex(vertexHash); err != nil {
		return 0, fmt.Errorf("could not find vertex with hash %v, %w", vertexHash, err)
	}

	degree := 0

	// Adding the number of ingoing edges is sufficient for undirected graphs, because all edges
	// exist twice (as two directed edges in opposite directions). Either dividing the number of
	// ingoing + outgoing edges by 2 or just using the number of ingoing edges is appropriate.
	inEdges, err := u.store.GetEdgesByTarget(vertexHash)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return 0, fmt.Errorf("could not find edges with hash %v, %w", vertexHash, err)
	}
	degree += len(inEdges)

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

	hashes, _ := u.store.ListVertices() // TODO: error
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
			if _, err := u.store.GetEdgesByTarget(vertex); err != nil {
				return nil, fmt.Errorf("vertex %v is not reachable from vertex %v, %w", targetHash, sourceHash, err)
			}
		}

		inEdges, err := u.store.GetEdgesByTarget(vertex)
		if err != nil {
			if !errors.Is(err, ErrNotFound) {
				return nil, fmt.Errorf("could not find edges by target with hash %v, %w", vertex, err)
			}
			continue
		}

		for _, edge := range inEdges {
			successor := edge.Source
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

func (u *undirected[K, T]) AdjacencyMap() map[K]map[K]Edge[K] {
	adjacencyMap := make(map[K]map[K]Edge[K])
	// ToDo(dominikbraun): Don't ignore this and the other error.
	vertices, _ := u.store.ListVertices()

	// Create an entry for each vertex to guarantee that all vertices are contained and its
	// adjacencies can be safely accessed without a preceding check.
	for _, vertexHash := range vertices {
		adjacencyMap[vertexHash] = make(map[K]Edge[K])
		edges, _ := u.store.GetEdgesBySource(vertexHash)

		for _, edge := range edges {
			adjacencyMap[vertexHash][edge.Target] = Edge[K]{
				Source: vertexHash,
				Target: edge.Target,
				Weight: edge.Weight,
			}
		}
	}

	return adjacencyMap
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

func (u *undirected[K, T]) adjacencies(vertexHash K) ([]K, error) {
	var adjacencyHashes []K

	// An undirected graph creates an undirected edge as two directed edges in the opposite
	// direction, so both the in-edges and the out-edges work here.
	inEdges, err := u.store.GetEdgesByTarget(vertexHash)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return adjacencyHashes, nil
		}
		return nil, fmt.Errorf("could not get edges by target with hash %v, %w", vertexHash, err)
	}

	for _, edge := range inEdges {
		adjacencyHashes = append(adjacencyHashes, edge.Source)
	}

	return adjacencyHashes, nil
}

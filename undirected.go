package graph

import (
	"errors"
	"fmt"
	"math"
)

type undirected[K comparable, T any] struct {
	hash     Hash[K, T]
	traits   *Traits
	vertices map[K]T
	outEdges map[K]map[K]Edge[T]
	inEdges  map[K]map[K]Edge[T]
}

func newUndirected[K comparable, T any](hash Hash[K, T], traits *Traits) *undirected[K, T] {
	return &undirected[K, T]{
		hash:     hash,
		traits:   traits,
		vertices: make(map[K]T),
		outEdges: make(map[K]map[K]Edge[T]),
		inEdges:  make(map[K]map[K]Edge[T]),
	}
}

func (u *undirected[K, T]) Traits() *Traits {
	return u.traits
}

func (u *undirected[K, T]) AddVertex(value T) {
	hash := u.hash(value)
	u.vertices[hash] = value
}

func (u *undirected[K, T]) Vertex(hash K) (T, error) {
	vertex, ok := u.vertices[hash]
	if !ok {
		return vertex, fmt.Errorf("vertex with hash %v doesn't exist", hash)
	}

	return vertex, nil
}

func (u *undirected[K, T]) AddEdge(sourceHash, targetHash K, options ...func(*EdgeProperties)) error {
	source, ok := u.vertices[sourceHash]
	if !ok {
		return fmt.Errorf("could not find source vertex with hash %v", sourceHash)
	}

	target, ok := u.vertices[targetHash]
	if !ok {
		return fmt.Errorf("could not find target vertex with hash %v", targetHash)
	}

	if _, ok := u.Edge(sourceHash, targetHash); ok {
		return fmt.Errorf("an edge between vertices %v and %v already exists", sourceHash, targetHash)
	}

	// If the graph was declared to be acyclic, permit the creation of a cycle.
	if u.traits.IsAcyclic {
		createsCycle, err := CreatesCycle[K, T](u, sourceHash, targetHash)
		if err != nil {
			return fmt.Errorf("failed to check for cycles: %w", err)
		}
		if createsCycle {
			return fmt.Errorf("an edge between %v and %v would introduce a cycle", sourceHash, targetHash)
		}
	}

	edge := Edge[T]{
		Source: source,
		Target: target,
		Properties: EdgeProperties{
			Attributes: make(map[string]string),
		},
	}

	for _, option := range options {
		option(&edge.Properties)
	}

	u.addEdge(sourceHash, targetHash, edge)

	return nil
}

func (u *undirected[K, T]) Edge(sourceHash, targetHash K) (Edge[T], bool) {
	// In an undirected graph, since multigraphs aren't supported, the edge AB is the same as BA.
	// Therefore, if source[target] cannot be found, this function also looks for target[source].
	sourceEdges, ok := u.outEdges[sourceHash]
	if ok {
		if edge, ok := sourceEdges[targetHash]; ok {
			return edge, true
		}
	}

	targetEdges, ok := u.outEdges[targetHash]
	if ok {
		if edge, ok := targetEdges[sourceHash]; ok {
			return edge, ok
		}
	}

	return Edge[T]{}, false
}

func (u *undirected[K, T]) Degree(vertexHash K) (int, error) {
	if _, ok := u.vertices[vertexHash]; !ok {
		return 0, fmt.Errorf("could not find vertex with hash %v", vertexHash)
	}

	degree := 0

	// Adding the number of ingoing edges is sufficient for undirected graphs, because all edges
	// exist twice (as two directed edges in opposite directions). Either dividing the number of
	// ingoing + outgoing edges by 2 or just using the number of ingoing edges is appropriate.
	if inEdges, ok := u.inEdges[vertexHash]; ok {
		degree += len(inEdges)
	}

	return degree, nil
}

func (u *undirected[K, T]) StronglyConnectedComponents() ([][]K, error) {
	return nil, errors.New("strongly connected components may only be detected in directed graphs")
}

// ShortestPath computes the shortest path between two vertices and returns the hashes of the
// vertices forming that path. The current implementation uses Dijkstra with a priority queue.
func (u *undirected[K, T]) ShortestPath(sourceHash, targetHash K) ([]K, error) {
	weights := make(map[K]float64)
	visited := make(map[K]bool)
	predecessors := make(map[K]K)

	weights[sourceHash] = 0
	visited[sourceHash] = true

	queue := newPriorityQueue[K]()

	for hash := range u.vertices {
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
			if _, ok := u.inEdges[vertex]; !ok {
				return nil, fmt.Errorf("vertex %v is not reachable from vertex %v", targetHash, sourceHash)
			}
		}

		inEdges, ok := u.inEdges[vertex]
		if !ok {
			continue
		}

		for successor, edge := range inEdges {
			weight := weights[vertex] + float64(edge.Properties.Weight)

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

	// Create an entry for each vertex to guarantee that all vertices are contained and its
	// adjacencies can be safely accessed without a preceding check.
	for vertexHash := range u.vertices {
		adjacencyMap[vertexHash] = make(map[K]Edge[K])
	}

	for vertexHash, outEdges := range u.outEdges {
		for adjacencyHash, edge := range outEdges {
			adjacencyMap[vertexHash][adjacencyHash] = Edge[K]{
				Source: vertexHash,
				Target: adjacencyHash,
				Properties: EdgeProperties{
					Weight:     edge.Properties.Weight,
					Attributes: edge.Properties.Attributes,
				},
			}
		}
	}

	return adjacencyMap
}

func (u *undirected[K, T]) Predecessors(vertex K) (map[K]Edge[K], error) {
	if _, ok := u.vertices[vertex]; !ok {
		return nil, fmt.Errorf("vertex with hash %v doesn't exist", vertex)
	}

	predecessors := make(map[K]Edge[K])

	for predecessor, edge := range u.inEdges[vertex] {
		predecessors[predecessor] = Edge[K]{
			Source:     vertex,
			Target:     predecessor,
			Properties: edge.Properties,
		}
	}

	return predecessors, nil
}

func (u *undirected[K, T]) edgesAreEqual(a, b Edge[T]) bool {
	aSourceHash := u.hash(a.Source)
	aTargetHash := u.hash(a.Target)
	bSourceHash := u.hash(b.Source)
	bTargetHash := u.hash(b.Target)

	if aSourceHash == bSourceHash && aTargetHash == bTargetHash {
		return true
	}

	if !u.traits.IsDirected {
		return aSourceHash == bTargetHash && aTargetHash == bSourceHash
	}

	return false
}

func (u *undirected[K, T]) addEdge(sourceHash, targetHash K, edge Edge[T]) {
	if _, ok := u.outEdges[sourceHash]; !ok {
		u.outEdges[sourceHash] = make(map[K]Edge[T])
	}
	if _, ok := u.outEdges[targetHash]; !ok {
		u.outEdges[targetHash] = make(map[K]Edge[T])
	}

	u.outEdges[sourceHash][targetHash] = edge
	u.outEdges[targetHash][sourceHash] = edge

	if _, ok := u.inEdges[targetHash]; !ok {
		u.inEdges[targetHash] = make(map[K]Edge[T])
	}
	if _, ok := u.inEdges[sourceHash]; !ok {
		u.inEdges[sourceHash] = make(map[K]Edge[T])
	}

	u.inEdges[targetHash][sourceHash] = edge
	u.inEdges[sourceHash][targetHash] = edge
}

func (u *undirected[K, T]) adjacencies(vertexHash K) []K {
	var adjacencyHashes []K

	// An undirected graph creates an undirected edge as two directed edges in the opposite
	// direction, so both the in-edges and the out-edges work here.
	inEdges, ok := u.inEdges[vertexHash]
	if !ok {
		return adjacencyHashes
	}

	for hash := range inEdges {
		adjacencyHashes = append(adjacencyHashes, hash)
	}

	return adjacencyHashes
}

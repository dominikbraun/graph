package graph

import (
	"errors"
	"fmt"
)

type undirected[K comparable, T any] struct {
	hash       Hash[K, T]
	properties *properties
	vertices   map[K]T
	edges      map[K]map[K]Edge[T]
	outEdges   map[K]map[K]Edge[T]
	inEdges    map[K]map[K]Edge[T]
}

func newUndirected[K comparable, T any](hash Hash[K, T], properties *properties) *undirected[K, T] {
	return &undirected[K, T]{
		hash:       hash,
		properties: properties,
		vertices:   make(map[K]T),
		edges:      make(map[K]map[K]Edge[T]),
		outEdges:   make(map[K]map[K]Edge[T]),
		inEdges:    make(map[K]map[K]Edge[T]),
	}
}

func (u *undirected[K, T]) Vertex(value T) {
	hash := u.hash(value)
	u.vertices[hash] = value
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
	source, ok := u.vertices[sourceHash]
	if !ok {
		return fmt.Errorf("could not find source vertex with hash %v", sourceHash)
	}

	target, ok := u.vertices[targetHash]
	if !ok {
		return fmt.Errorf("could not find target vertex with hash %v", targetHash)
	}

	if _, ok := u.GetEdgeByHashes(sourceHash, targetHash); ok {
		return fmt.Errorf("an edge between vertices %v and %v already exists", sourceHash, targetHash)
	}

	// If the graph was declared to be acyclic, permit the creation of a cycle.
	if u.properties.isAcyclic {
		createsCycle, err := u.CreatesCycleByHashes(sourceHash, targetHash)
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
		Weight: weight,
	}

	u.addEdge(sourceHash, targetHash, edge)

	return nil
}

func (u *undirected[K, T]) GetEdge(source, target T) (Edge[T], bool) {
	sourceHash := u.hash(source)
	targetHash := u.hash(target)

	return u.GetEdgeByHashes(sourceHash, targetHash)
}

func (u *undirected[K, T]) GetEdgeByHashes(sourceHash, targetHash K) (Edge[T], bool) {
	// In an undirected graph, since multigraphs aren't supported, the edge AB is the same as BA.
	// Therefore, if source[target] cannot be found, this function also looks for target[source].
	sourceEdges, ok := u.edges[sourceHash]
	if ok {
		if edge, ok := sourceEdges[targetHash]; ok {
			return edge, true
		}
	}

	targetEdges, ok := u.edges[targetHash]
	if ok {
		if edge, ok := targetEdges[sourceHash]; ok {
			return edge, ok
		}
	}

	return Edge[T]{}, false
}

func (u *undirected[K, T]) DFS(start T, visit func(value T) bool) error {
	startHash := u.hash(start)

	return u.DFSByHash(startHash, visit)
}

func (u *undirected[K, T]) DFSByHash(startHash K, visit func(value T) bool) error {
	if _, ok := u.vertices[startHash]; !ok {
		return fmt.Errorf("could not find start vertex with hash %v", startHash)
	}

	stack := make([]K, 0)
	visited := make(map[K]bool)

	stack = append(stack, startHash)

	for len(stack) > 0 {
		currentHash := stack[len(stack)-1]
		currentVertex := u.vertices[currentHash]

		stack = stack[:len(stack)-1]

		if _, ok := visited[currentHash]; !ok {
			// Stop traversing the graph if the visit function returns true.
			if visit(currentVertex) {
				break
			}
			visited[currentHash] = true

			for adjacency := range u.outEdges[currentHash] {
				stack = append(stack, adjacency)
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
	if _, ok := u.vertices[startHash]; !ok {
		return fmt.Errorf("could not find start vertex with hash %v", startHash)
	}

	queue := make([]K, 0)
	visited := make(map[K]bool)

	visited[startHash] = true
	queue = append(queue, startHash)

	for len(queue) > 0 {
		currentHash := queue[0]
		currentVertex := u.vertices[currentHash]

		queue = queue[1:]

		// Stop traversing the graph if the visit function returns true.
		if visit(currentVertex) {
			break
		}

		for adjacency := range u.outEdges[currentHash] {
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
	source, ok := u.vertices[sourceHash]
	if !ok {
		return false, fmt.Errorf("could not find source vertex with hash %v", source)
	}

	_, ok = u.vertices[targetHash]
	if !ok {
		return false, fmt.Errorf("could not find target vertex with hash %v", source)
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

func (u *undirected[K, T]) edgesAreEqual(a, b Edge[T]) bool {
	aSourceHash := u.hash(a.Source)
	aTargetHash := u.hash(a.Target)
	bSourceHash := u.hash(b.Source)
	bTargetHash := u.hash(b.Target)

	if aSourceHash == bSourceHash && aTargetHash == bTargetHash {
		return true
	}

	if !u.properties.isDirected {
		return aSourceHash == bTargetHash && aTargetHash == bSourceHash
	}

	return false
}

func (u *undirected[K, T]) addEdge(sourceHash, targetHash K, edge Edge[T]) {
	if _, ok := u.edges[sourceHash]; !ok {
		u.edges[sourceHash] = make(map[K]Edge[T])
	}
	if _, ok := u.edges[targetHash]; !ok {
		u.edges[targetHash] = make(map[K]Edge[T])
	}

	u.edges[sourceHash][targetHash] = edge
	u.edges[targetHash][sourceHash] = edge

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

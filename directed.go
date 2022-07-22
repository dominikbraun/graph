package graph

import (
	"fmt"
	"math"
)

type directed[K comparable, T any] struct {
	hash     Hash[K, T]
	traits   *traits
	vertices map[K]T
	edges    map[K]map[K]Edge[T]
	outEdges map[K]map[K]Edge[T]
	inEdges  map[K]map[K]Edge[T]
}

func newDirected[K comparable, T any](hash Hash[K, T], traits *traits) *directed[K, T] {
	return &directed[K, T]{
		hash:     hash,
		traits:   traits,
		vertices: make(map[K]T),
		edges:    make(map[K]map[K]Edge[T]),
		outEdges: make(map[K]map[K]Edge[T]),
		inEdges:  make(map[K]map[K]Edge[T]),
	}
}

func (d *directed[K, T]) Vertex(value T) {
	hash := d.hash(value)
	d.vertices[hash] = value
}

func (d *directed[K, T]) Edge(source, target T) error {
	return d.WeightedEdge(source, target, 0)
}

func (d *directed[K, T]) WeightedEdge(source, target T, weight int) error {
	sourceHash := d.hash(source)
	targetHash := d.hash(target)

	return d.WeightedEdgeByHashes(sourceHash, targetHash, weight)
}

func (d *directed[K, T]) EdgeByHashes(sourceHash, targetHash K) error {
	return d.WeightedEdgeByHashes(sourceHash, targetHash, 0)
}

func (d *directed[K, T]) WeightedEdgeByHashes(sourceHash, targetHash K, weight int) error {
	source, ok := d.vertices[sourceHash]
	if !ok {
		return fmt.Errorf("could not find source vertex with hash %v", sourceHash)
	}

	target, ok := d.vertices[targetHash]
	if !ok {
		return fmt.Errorf("could not find target vertex with hash %v", targetHash)
	}

	if _, ok := d.GetEdgeByHashes(sourceHash, targetHash); ok {
		return fmt.Errorf("an edge between vertices %v and %v already exists", sourceHash, targetHash)
	}

	// If the graph was declared to be acyclic, permit the creation of a cycle.
	if d.traits.isAcyclic {
		createsCycle, err := d.CreatesCycleByHashes(sourceHash, targetHash)
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

	d.addEdge(sourceHash, targetHash, edge)

	return nil
}

func (d *directed[K, T]) GetEdge(source, target T) (Edge[T], bool) {
	sourceHash := d.hash(source)
	targetHash := d.hash(target)

	return d.GetEdgeByHashes(sourceHash, targetHash)
}

func (d *directed[K, T]) GetEdgeByHashes(sourceHash, targetHash K) (Edge[T], bool) {
	sourceEdges, ok := d.edges[sourceHash]
	if !ok {
		return Edge[T]{}, false
	}

	if edge, ok := sourceEdges[targetHash]; ok {
		return edge, true
	}

	return Edge[T]{}, false
}

func (d *directed[K, T]) DFS(start T, visit func(value T) bool) error {
	startHash := d.hash(start)

	return d.DFSByHash(startHash, visit)
}

func (d *directed[K, T]) DFSByHash(startHash K, visit func(value T) bool) error {
	if _, ok := d.vertices[startHash]; !ok {
		return fmt.Errorf("could not find start vertex with hash %v", startHash)
	}

	stack := make([]K, 0)
	visited := make(map[K]bool)

	stack = append(stack, startHash)

	for len(stack) > 0 {
		currentHash := stack[len(stack)-1]
		currentVertex := d.vertices[currentHash]

		stack = stack[:len(stack)-1]

		if _, ok := visited[currentHash]; !ok {
			// Stop traversing the graph if the visit function returns true.
			if visit(currentVertex) {
				break
			}
			visited[currentHash] = true

			for adjacency := range d.outEdges[currentHash] {
				stack = append(stack, adjacency)
			}
		}
	}

	return nil
}

func (d *directed[K, T]) BFS(start T, visit func(value T) bool) error {
	startHash := d.hash(start)

	return d.BFSByHash(startHash, visit)
}

func (d *directed[K, T]) BFSByHash(startHash K, visit func(value T) bool) error {
	if _, ok := d.vertices[startHash]; !ok {
		return fmt.Errorf("could not find start vertex with hash %v", startHash)
	}

	queue := make([]K, 0)
	visited := make(map[K]bool)

	visited[startHash] = true
	queue = append(queue, startHash)

	for len(queue) > 0 {
		currentHash := queue[0]
		currentVertex := d.vertices[currentHash]

		queue = queue[1:]

		// Stop traversing the graph if the visit function returns true.
		if visit(currentVertex) {
			break
		}

		for adjacency := range d.outEdges[currentHash] {
			if _, ok := visited[adjacency]; !ok {
				visited[adjacency] = true
				queue = append(queue, adjacency)
			}
		}

	}

	return nil
}

func (d *directed[K, T]) CreatesCycle(source, target T) (bool, error) {
	sourceHash := d.hash(source)
	targetHash := d.hash(target)

	return d.CreatesCycleByHashes(sourceHash, targetHash)
}

func (d *directed[K, T]) CreatesCycleByHashes(sourceHash, targetHash K) (bool, error) {
	source, ok := d.vertices[sourceHash]
	if !ok {
		return false, fmt.Errorf("could not find source vertex with hash %v", source)
	}

	_, ok = d.vertices[targetHash]
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
			// If the current vertex, e.g. a predecessor of the source vertex, also is the target
			// vertex, an edge between these two would create a cycle.
			if currentHash == targetHash {
				return true, nil
			}
			visited[currentHash] = true

			for _, predecessor := range d.predecessors(currentHash) {
				stack = append(stack, predecessor)
			}
		}
	}

	return false, nil
}

func (d *directed[K, T]) Degree(vertex T) (int, error) {
	sourceHash := d.hash(vertex)

	return d.DegreeByHash(sourceHash)
}

func (d *directed[K, T]) DegreeByHash(vertexHash K) (int, error) {
	if _, ok := d.vertices[vertexHash]; !ok {
		return 0, fmt.Errorf("could not find vertex with hash %v", vertexHash)
	}

	degree := 0

	if inEdges, ok := d.inEdges[vertexHash]; ok {
		degree += len(inEdges)
	}
	if outEdges, ok := d.outEdges[vertexHash]; ok {
		degree += len(outEdges)
	}

	return degree, nil
}

type sccState[K comparable] struct {
	components [][]K
	stack      []K
	onStack    map[K]bool
	visited    map[K]bool
	lowlink    map[K]int
	index      map[K]int
	time       int
}

// StronglyConnectedComponents searches strongly connected components within the graph, and returns
// the hashes of the vertices shaping these components. The current implementation of this function
// uses Tarjan's algorithm and runs recursively.
func (d *directed[K, T]) StronglyConnectedComponents() ([][]K, error) {

	state := &sccState[K]{
		components: make([][]K, 0),
		stack:      make([]K, 0),
		onStack:    make(map[K]bool),
		visited:    make(map[K]bool),
		lowlink:    make(map[K]int),
		index:      make(map[K]int),
	}

	for hash := range d.vertices {
		if ok, _ := state.visited[hash]; !ok {
			d.findSCC(hash, state)
		}
	}

	return state.components, nil
}

func (d *directed[K, T]) findSCC(vertexHash K, state *sccState[K]) {
	state.stack = append(state.stack, vertexHash)
	state.onStack[vertexHash] = true
	state.visited[vertexHash] = true
	state.index[vertexHash] = state.time
	state.lowlink[vertexHash] = state.time

	state.time++

	for adjancency := range d.outEdges[vertexHash] {
		if ok, _ := state.visited[adjancency]; !ok {
			d.findSCC(adjancency, state)

			smallestLowlink := math.Min(
				float64(state.lowlink[vertexHash]),
				float64(state.lowlink[adjancency]),
			)
			state.lowlink[vertexHash] = int(smallestLowlink)
		} else {
			// If the adjacent vertex already is on the stack, the edge joining the current and the
			// adjacent vertex is a back edge. Therefore, update the vertex' lowlink value to the
			// index of the adjacent vertex if it is smaller than the lowlink value.
			if ok, _ := state.onStack[adjancency]; ok {
				smallestLowlink := math.Min(
					float64(state.lowlink[vertexHash]),
					float64(state.index[adjancency]),
				)
				state.lowlink[vertexHash] = int(smallestLowlink)
			}
		}
	}

	// If the lowlink value of the vertex is equal to its DFS index, this is th head vertex of a
	// strongly connected component, shaped by this vertex and the vertices on the stack.
	if state.lowlink[vertexHash] == state.index[vertexHash] {
		var hash K
		var component []K

		for hash != vertexHash {
			hash = state.stack[len(state.stack)-1]
			state.stack = state.stack[:len(state.stack)-1]
			state.onStack[hash] = false

			component = append(component, hash)
		}

		state.components = append(state.components, component)
	}
}

// ShortestPath computes the shortest path between two vertices and returns the hashes of the
// vertices forming that path. The current implementation uses Dijkstra with a priority queue.
func (d *directed[K, T]) ShortestPath(source, target T) ([]K, error) {
	sourceHash := d.hash(source)
	targetHash := d.hash(target)

	return d.ShortestPathByHashes(sourceHash, targetHash)
}

func (d *directed[K, T]) ShortestPathByHashes(sourceHash, targetHash K) ([]K, error) {
	weights := make(map[K]float64)
	visited := make(map[K]bool)
	predecessors := make(map[K]K)

	weights[sourceHash] = 0
	visited[sourceHash] = true

	queue := newPriorityQueue[K]()

	for hash := range d.vertices {
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
			if _, ok := d.inEdges[vertex]; !ok {
				return nil, fmt.Errorf("vertex %v is not reachable from vertex %v", targetHash, sourceHash)
			}
		}

		outEdges, ok := d.outEdges[vertex]
		if !ok {
			continue
		}

		for successor, edge := range outEdges {
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

func (d *directed[K, T]) AdjacencyMap() map[K]map[K]Edge[K] {
	adjacencyMap := make(map[K]map[K]Edge[K])

	// Create an entry for each vertex to guarantee that all vertices are contained and its
	// adjacencies can be safely accessed without a preceding check.
	for vertexHash := range d.vertices {
		adjacencyMap[vertexHash] = make(map[K]Edge[K])
	}

	for vertexHash, outEdges := range d.outEdges {
		for adjacencyHash, edge := range outEdges {
			adjacencyMap[vertexHash][adjacencyHash] = Edge[K]{
				Source: vertexHash,
				Target: adjacencyHash,
				Weight: edge.Weight,
			}
		}
	}

	return adjacencyMap
}

func (d *directed[K, T]) edgesAreEqual(a, b Edge[T]) bool {
	aSourceHash := d.hash(a.Source)
	aTargetHash := d.hash(a.Target)
	bSourceHash := d.hash(b.Source)
	bTargetHash := d.hash(b.Target)

	return aSourceHash == bSourceHash && aTargetHash == bTargetHash
}

func (d *directed[K, T]) addEdge(sourceHash, targetHash K, edge Edge[T]) {
	if _, ok := d.edges[sourceHash]; !ok {
		d.edges[sourceHash] = make(map[K]Edge[T])
	}

	d.edges[sourceHash][targetHash] = edge

	if _, ok := d.outEdges[sourceHash]; !ok {
		d.outEdges[sourceHash] = make(map[K]Edge[T])
	}

	d.outEdges[sourceHash][targetHash] = edge

	if _, ok := d.inEdges[targetHash]; !ok {
		d.inEdges[targetHash] = make(map[K]Edge[T])
	}

	d.inEdges[targetHash][sourceHash] = edge
}

func (d *directed[K, T]) predecessors(vertexHash K) []K {
	var predecessorHashes []K

	inEdges, ok := d.inEdges[vertexHash]
	if !ok {
		return predecessorHashes
	}

	for hash := range inEdges {
		predecessorHashes = append(predecessorHashes, hash)
	}

	return predecessorHashes
}

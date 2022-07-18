package graph

import (
	"fmt"
	"math"
)

type directed[K comparable, T any] struct {
	hash   Hash[K, T]
	traits *traits
	store  Store[K, T]
}

func newDirected[K comparable, T any](
	hash Hash[K, T],
	store Store[K, T],
	traits *traits,
) *directed[K, T] {
	return &directed[K, T]{
		hash:   hash,
		traits: traits,
		store:  store,
	}
}

func (d *directed[K, T]) Vertex(value T) {
	d.store.AddVertex(d.hash(value), value)
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
	_, ok := d.store.GetVertex(sourceHash)
	if !ok {
		return fmt.Errorf("could not find source vertex with hash %v", sourceHash)
	}

	_, ok = d.store.GetVertex(targetHash)
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

	edge := Edge[K]{
		Source: sourceHash,
		Target: targetHash,
		Weight: weight,
	}

	d.store.AddEdge(sourceHash, targetHash, edge)

	return nil
}

func (d *directed[K, T]) GetEdge(source, target T) (Edge[K], bool) {
	sourceHash := d.hash(source)
	targetHash := d.hash(target)

	return d.GetEdgeByHashes(sourceHash, targetHash)
}

func (d *directed[K, T]) GetEdgeByHashes(sourceHash, targetHash K) (Edge[K], bool) {
	return d.store.GetEdge(sourceHash, targetHash)
}

func (d *directed[K, T]) DFS(start T, visit func(value T) bool) error {
	startHash := d.hash(start)

	return d.DFSByHash(startHash, visit)
}

func (d *directed[K, T]) DFSByHash(startHash K, visit func(value T) bool) error {
	if _, ok := d.store.GetVertex(startHash); !ok {
		return fmt.Errorf("could not find start vertex with hash %v", startHash)
	}

	stack := make([]K, 0)
	visited := make(map[K]bool)

	stack = append(stack, startHash)

	for len(stack) > 0 {
		currentHash := stack[len(stack)-1]
		currentVertex, _ := d.store.GetVertex(currentHash)

		stack = stack[:len(stack)-1]

		if _, ok := visited[currentHash]; !ok {
			// Stop traversing the graph if the visit function returns true.
			if visit(*currentVertex) {
				break
			}
			visited[currentHash] = true

			edges, _ := d.store.GetEdgesBySource(currentHash) // TODO: error
			for _, edge := range edges {
				stack = append(stack, edge.Target)
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
	if _, ok := d.store.GetVertex(startHash); !ok {
		return fmt.Errorf("could not find start vertex with hash %v", startHash)
	}

	queue := make([]K, 0)
	visited := make(map[K]bool)

	visited[startHash] = true
	queue = append(queue, startHash)

	for len(queue) > 0 {
		currentHash := queue[0]
		currentVertex, _ := d.store.GetVertex(currentHash)

		queue = queue[1:]

		// Stop traversing the graph if the visit function returns true.
		if visit(*currentVertex) {
			break
		}

		edges, _ := d.store.GetEdgesBySource(currentHash) // TODO: error
		for _, edge := range edges {
			if _, ok := visited[edge.Target]; !ok {
				visited[edge.Target] = true
				queue = append(queue, edge.Target)
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
	_, ok := d.store.GetVertex(sourceHash)
	if !ok {
		return false, fmt.Errorf("could not find source vertex with hash %v", sourceHash)
	}

	_, ok = d.store.GetVertex(targetHash)
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
	if _, ok := d.store.GetVertex(vertexHash); !ok {
		return 0, fmt.Errorf("could not find vertex with hash %v", vertexHash)
	}

	degree := 0

	if inEdges, ok := d.store.GetEdgesByTarget(vertexHash); ok {
		degree += len(inEdges)
	}
	if outEdges, ok := d.store.GetEdgesBySource(vertexHash); ok {
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

	hashes, _ := d.store.ListVertices() // TODO: error
	for _, hash := range hashes {
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

	edges, _ := d.store.GetEdgesBySource(vertexHash) // TODO: error
	for _, edge := range edges {
		adjancency := edge.Target
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

	hashes, _ := d.store.ListVertices() // TODO: error
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
			if _, ok := d.store.GetEdgesByTarget(vertex); !ok {
				return nil, fmt.Errorf("vertex %v is not reachable from vertex %v", targetHash, sourceHash)
			}
		}

		outEdges, ok := d.store.GetEdgesBySource(vertex)
		if !ok {
			continue
		}

		for _, edge := range outEdges {
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

func (d *directed[K, T]) edgesAreEqual(a, b Edge[K]) bool {
	return a.Source == b.Source && a.Target == b.Target
}

func (d *directed[K, T]) predecessors(vertexHash K) []K {
	var predecessorHashes []K

	inEdges, ok := d.store.GetEdgesByTarget(vertexHash)
	if !ok {
		return predecessorHashes
	}

	for _, edge := range inEdges {
		predecessorHashes = append(predecessorHashes, edge.Source)
	}

	return predecessorHashes
}

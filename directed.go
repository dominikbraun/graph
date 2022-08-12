package graph

import (
	"fmt"
	"math"
)

type directed[K comparable, T any] struct {
	hash     Hash[K, T]
	traits   *Traits
	vertices map[K]T
	edges    map[K]map[K]Edge[T]
	outEdges map[K]map[K]Edge[T]
	inEdges  map[K]map[K]Edge[T]
}

func newDirected[K comparable, T any](hash Hash[K, T], traits *Traits) *directed[K, T] {
	return &directed[K, T]{
		hash:     hash,
		traits:   traits,
		vertices: make(map[K]T),
		edges:    make(map[K]map[K]Edge[T]),
		outEdges: make(map[K]map[K]Edge[T]),
		inEdges:  make(map[K]map[K]Edge[T]),
	}
}

func (d *directed[K, T]) Traits() *Traits {
	return d.traits
}

func (d *directed[K, T]) AddVertex(value T) {
	hash := d.hash(value)
	d.vertices[hash] = value
}

func (d *directed[K, T]) Vertex(hash K) (T, error) {
	vertex, ok := d.vertices[hash]
	if !ok {
		return vertex, fmt.Errorf("vertex with hash %v doesn't exist", hash)
	}

	return vertex, nil
}

func (d *directed[K, T]) AddEdge(sourceHash, targetHash K, options ...func(*EdgeProperties)) error {
	source, ok := d.vertices[sourceHash]
	if !ok {
		return fmt.Errorf("could not find source vertex with hash %v", sourceHash)
	}

	target, ok := d.vertices[targetHash]
	if !ok {
		return fmt.Errorf("could not find target vertex with hash %v", targetHash)
	}

	if _, ok := d.Edge(sourceHash, targetHash); ok {
		return fmt.Errorf("an edge between vertices %v and %v already exists", sourceHash, targetHash)
	}

	// If the graph was declared to be acyclic, permit the creation of a cycle.
	if d.traits.IsAcyclic {
		createsCycle, err := CreatesCycle[K, T](d, sourceHash, targetHash)
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

	d.addEdge(sourceHash, targetHash, edge)

	return nil
}

func (d *directed[K, T]) Edge(sourceHash, targetHash K) (Edge[T], bool) {
	sourceEdges, ok := d.edges[sourceHash]
	if !ok {
		return Edge[T]{}, false
	}

	if edge, ok := sourceEdges[targetHash]; ok {
		return edge, true
	}

	return Edge[T]{}, false
}

func (d *directed[K, T]) Degree(vertexHash K) (int, error) {
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
func (d *directed[K, T]) ShortestPath(sourceHash, targetHash K) ([]K, error) {
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
				Properties: EdgeProperties{
					Weight:     edge.Properties.Weight,
					Attributes: edge.Properties.Attributes,
				},
			}
		}
	}

	return adjacencyMap
}

func (d *directed[K, T]) Predecessors(vertex K) (map[K]Edge[K], error) {
	if _, ok := d.vertices[vertex]; !ok {
		return nil, fmt.Errorf("vertex with hash %v doesn't exist", vertex)
	}

	predecessors := make(map[K]Edge[K])

	for predecessor, edge := range d.inEdges[vertex] {
		predecessors[predecessor] = Edge[K]{
			Source:     vertex,
			Target:     predecessor,
			Properties: edge.Properties,
		}
	}

	return predecessors, nil
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

package graph

import "fmt"

type directed[K comparable, T any] struct {
	hash       Hash[K, T]
	properties *properties
	vertices   map[K]T
	edges      map[K]map[K]Edge[T]
	outEdges   map[K]map[K]Edge[T]
	inEdges    map[K]map[K]Edge[T]
}

func newDirected[K comparable, T any](hash Hash[K, T], properties *properties) *directed[K, T] {
	return &directed[K, T]{
		hash:       hash,
		properties: properties,
		vertices:   make(map[K]T),
		edges:      make(map[K]map[K]Edge[T]),
		outEdges:   make(map[K]map[K]Edge[T]),
		inEdges:    make(map[K]map[K]Edge[T]),
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
	if !ok && d.properties.isDirected {
		return Edge[T]{}, false
	}

	if edge, ok := sourceEdges[targetHash]; ok {
		return edge, true
	}

	if !d.properties.isDirected {
		targetEdges, ok := d.edges[targetHash]
		if !ok {
			return Edge[T]{}, false
		}

		if edge, ok := targetEdges[sourceHash]; ok {
			return edge, true
		}
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

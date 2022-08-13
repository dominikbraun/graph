package graph

import (
	"fmt"
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

func (d *directed[K, T]) AdjacencyMap() (map[K]map[K]Edge[K], error) {
	adjacencyMap := make(map[K]map[K]Edge[K])

	// Create an entry for each vertex to guarantee that all vertices are contained and its
	// adjacent vertices can be safely accessed without a preceding check.
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

	return adjacencyMap, nil
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

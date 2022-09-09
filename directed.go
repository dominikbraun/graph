package graph

import (
	"errors"
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

func (d *directed[K, T]) AddVertex(value T) error {
	hash := d.hash(value)
	d.vertices[hash] = value

	return nil
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

	if _, err := d.Edge(sourceHash, targetHash); !errors.Is(err, ErrEdgeNotFound) {
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

func (d *directed[K, T]) Edge(sourceHash, targetHash K) (Edge[T], error) {
	sourceEdges, ok := d.edges[sourceHash]
	if !ok {
		return Edge[T]{}, ErrEdgeNotFound
	}

	edge, ok := sourceEdges[targetHash]
	if !ok {
		return Edge[T]{}, ErrEdgeNotFound
	}

	return edge, nil
}

func (d *directed[K, T]) RemoveEdge(source, target K) error {
	if _, err := d.Edge(source, target); err != nil {
		return fmt.Errorf("failed to find edge from %v to %v: %w", source, target, err)
	}

	delete(d.edges[source], target)
	delete(d.inEdges[target], source)
	delete(d.outEdges[source], target)

	return nil
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

func (d *directed[K, T]) PredecessorMap() (map[K]map[K]Edge[K], error) {
	predecessors := make(map[K]map[K]Edge[K])

	for vertexHash := range d.vertices {
		predecessors[vertexHash] = make(map[K]Edge[K])
	}

	for vertexHash, inEdges := range d.inEdges {
		for predecessorHash, edge := range inEdges {
			predecessors[vertexHash][predecessorHash] = Edge[K]{
				Source: predecessorHash,
				Target: vertexHash,
				Properties: EdgeProperties{
					Attributes: edge.Properties.Attributes,
					Weight:     edge.Properties.Weight,
				},
			}
		}
	}

	return predecessors, nil
}

func (d *directed[K, T]) Clone() (Graph[K, T], error) {
	traits := &Traits{
		IsDirected: d.traits.IsDirected,
		IsAcyclic:  d.traits.IsAcyclic,
		IsWeighted: d.traits.IsWeighted,
		IsRooted:   d.traits.IsRooted,
	}

	vertices := make(map[K]T)

	for hash, vertex := range d.vertices {
		vertices[hash] = vertex
	}

	return &directed[K, T]{
		hash:     d.hash,
		traits:   traits,
		vertices: vertices,
		edges:    cloneEdges(d.edges),
		outEdges: cloneEdges(d.outEdges),
		inEdges:  cloneEdges(d.inEdges),
	}, nil
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

func cloneEdges[K comparable, T any](input map[K]map[K]Edge[T]) map[K]map[K]Edge[T] {
	edges := make(map[K]map[K]Edge[T])

	for hash, neighbours := range input {
		edges[hash] = make(map[K]Edge[T])

		for neighbourHash, edge := range neighbours {
			attributes := make(map[string]string)

			for key, value := range edge.Properties.Attributes {
				attributes[key] = value
			}

			edges[hash][neighbourHash] = Edge[T]{
				Source: edge.Source,
				Target: edge.Target,
				Properties: EdgeProperties{
					Attributes: attributes,
					Weight:     edge.Properties.Weight,
				},
			}
		}
	}

	return edges
}

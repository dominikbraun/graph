package graph

import (
	"fmt"
	"sync"
)

// Store represents a storage for vertices and edges. The graph library provides an in-memory store
// by default and accepts any Store implementation to work with - for example, an SQL store.
//
// When implementing your own Store, make sure the individual methods and their behavior adhere to
// this documentation. Otherwise, the graphs aren't guaranteed to behave as expected.
type Store[K comparable, T any] interface {
	// AddVertex should add the given vertex with the given hash value and vertex properties to the
	// graph. If the vertex already exists, it is up to you whether ErrVertexAlreadyExists or no
	// error should be returned.
	AddVertex(hash K, value T, properties VertexProperties) error

	// Vertex should return the vertex and vertex properties with the given hash value. If the
	// vertex doesn't exist, ErrVertexNotFound should be returned.
	Vertex(hash K) (T, VertexProperties, error)

	// RemoveVertex should remove the vertex with the given hash value. If the vertex doesn't
	// exist, ErrVertexNotFound should be returned. If the vertex has edges to other vertices,
	// ErrVertexHasEdges should be returned.
	RemoveVertex(hash K) error

	// ListVertices should return all vertices in the graph in a slice.
	ListVertices() ([]K, error)

	// VertexCount should return the number of vertices in the graph. This should be equal to the
	// length of the slice returned by ListVertices.
	VertexCount() (int, error)

	// AddEdge should add an edge between the vertices with the given source and target hashes.
	//
	// If either vertex doesn't exit, ErrVertexNotFound should be returned for the respective
	// vertex. If the edge already exists, ErrEdgeAlreadyExists should be returned.
	AddEdge(sourceHash, targetHash K, edge Edge[K]) error

	// UpdateEdge should update the edge between the given vertices with the data of the given
	// Edge instance. If the edge doesn't exist, ErrEdgeNotFound should be returned.
	UpdateEdge(sourceHash, targetHash K, edge Edge[K]) error

	// RemoveEdge should remove the edge between the vertices with the given source and target
	// hashes.
	//
	// If either vertex doesn't exist, it is up to you whether ErrVertexNotFound or no error should
	// be returned. If the edge doesn't exist, it is up to you whether ErrEdgeNotFound or no error
	// should be returned.
	RemoveEdge(sourceHash, targetHash K) error

	// Edge should return the edge joining the vertices with the given hash values. It should
	// exclusively look for an edge between the source and the target vertex, not vice versa. The
	// graph implementation does this for undirected graphs itself.
	//
	// Note that unlike Graph.Edge, this function is supposed to return an Edge[K], i.e. an edge
	// that only contains the vertex hashes instead of the vertices themselves.
	//
	// If the edge doesn't exist, ErrEdgeNotFound should be returned.
	Edge(sourceHash, targetHash K) (Edge[K], error)

	// ListEdges should return all edges in the graph in a slice.
	ListEdges() ([]Edge[K], error)

	// EdgeCount should return the number of edges in the graph. This should be equal to the
	// length of the slice returned by ListEdges.
	EdgeCount() (int, error)
}

type MemoryStore[K comparable, T any] struct {
	lock             sync.RWMutex
	vertices         map[K]T
	vertexProperties map[K]VertexProperties

	// outEdges and inEdges store all outgoing and ingoing edges for all vertices. For O(1) access,
	// these edges themselves are stored in maps whose keys are the hashes of the target vertices.
	outEdges  map[K]map[K]Edge[K] // source -> target
	inEdges   map[K]map[K]Edge[K] // target -> source
	edgeCount int
}

func NewMemoryStore[K comparable, T any]() Store[K, T] {
	return &MemoryStore[K, T]{
		vertices:         make(map[K]T),
		vertexProperties: make(map[K]VertexProperties),
		outEdges:         make(map[K]map[K]Edge[K]),
		inEdges:          make(map[K]map[K]Edge[K]),
	}
}

func (s *MemoryStore[K, T]) AddVertex(k K, t T, p VertexProperties) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.vertices[k]; ok {
		return ErrVertexAlreadyExists
	}

	s.vertices[k] = t
	s.vertexProperties[k] = p

	return nil
}

func (s *MemoryStore[K, T]) ListVertices() ([]K, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	hashes := make([]K, 0, len(s.vertices))
	for k := range s.vertices {
		hashes = append(hashes, k)
	}

	return hashes, nil
}

func (s *MemoryStore[K, T]) VertexCount() (int, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return len(s.vertices), nil
}

func (s *MemoryStore[K, T]) Vertex(k K) (T, VertexProperties, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	v, ok := s.vertices[k]
	if !ok {
		return v, VertexProperties{}, ErrVertexNotFound
	}

	p := s.vertexProperties[k]

	return v, p, nil
}

func (s *MemoryStore[K, T]) RemoveVertex(k K) error {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if _, ok := s.vertices[k]; !ok {
		return ErrVertexNotFound
	}

	if edges, ok := s.inEdges[k]; ok {
		if len(edges) > 0 {
			return ErrVertexHasEdges
		}
		delete(s.inEdges, k)
	}

	if edges, ok := s.outEdges[k]; ok {
		if len(edges) > 0 {
			return ErrVertexHasEdges
		}
		delete(s.outEdges, k)
	}

	delete(s.vertices, k)
	delete(s.vertexProperties, k)

	return nil
}

func (s *MemoryStore[K, T]) AddEdge(sourceHash, targetHash K, edge Edge[K]) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.outEdges[sourceHash]; !ok {
		s.outEdges[sourceHash] = make(map[K]Edge[K])
	}

	s.outEdges[sourceHash][targetHash] = edge

	if _, ok := s.inEdges[targetHash]; !ok {
		s.inEdges[targetHash] = make(map[K]Edge[K])
	}

	s.inEdges[targetHash][sourceHash] = edge

	s.edgeCount++

	return nil
}

func (s *MemoryStore[K, T]) UpdateEdge(sourceHash, targetHash K, edge Edge[K]) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	targetEdges, ok := s.outEdges[sourceHash]
	if !ok {
		return ErrEdgeNotFound
	}

	_, ok = targetEdges[targetHash]
	if !ok {
		return ErrEdgeNotFound
	}

	s.outEdges[sourceHash][targetHash] = edge
	s.inEdges[targetHash][sourceHash] = edge

	return nil
}

func (s *MemoryStore[K, T]) RemoveEdge(sourceHash, targetHash K) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.inEdges[targetHash], sourceHash)
	delete(s.outEdges[sourceHash], targetHash)

	s.edgeCount--

	return nil
}

func (s *MemoryStore[K, T]) Edge(sourceHash, targetHash K) (Edge[K], error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	sourceEdges, ok := s.outEdges[sourceHash]
	if !ok {
		return Edge[K]{}, ErrEdgeNotFound
	}

	edge, ok := sourceEdges[targetHash]
	if !ok {
		return Edge[K]{}, ErrEdgeNotFound
	}

	return edge, nil
}

func (s *MemoryStore[K, T]) EdgeCount() (int, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.edgeCount, nil
}

func (s *MemoryStore[K, T]) ListEdges() ([]Edge[K], error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	res := make([]Edge[K], 0, s.edgeCount)
	for _, edges := range s.outEdges {
		for _, edge := range edges {
			res = append(res, edge)
		}
	}
	return res, nil
}

// CreatesCycle is a fastpath version of [CreatesCycle] that avoids calling
// [PredecessorMap], which generates large amounts of garbage to collect.
//
// Because CreatesCycle doesn't need to modify the PredecessorMap, we can use
// inEdges instead to compute the same thing without creating any copies.
func (s *MemoryStore[K, T]) CreatesCycle(source, target K) (bool, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if _, ok := s.vertices[source]; !ok {
		return false, fmt.Errorf("could not get vertex with hash %v", source)
	}

	if _, ok := s.vertices[target]; !ok {
		return false, fmt.Errorf("could not get vertex with hash %v", target)
	}

	if source == target {
		return true, nil
	}

	stack := newStack[K]()
	visited := make(map[K]struct{})

	stack.push(source)

	for !stack.isEmpty() {
		currentHash, _ := stack.pop()

		if _, ok := visited[currentHash]; !ok {
			// If the adjacent vertex also is the target vertex, the target is a
			// parent of the source vertex. An edge would introduce a cycle.
			if currentHash == target {
				return true, nil
			}

			visited[currentHash] = struct{}{}

			for adjacency := range s.inEdges[currentHash] {
				stack.push(adjacency)
			}
		}
	}

	return false, nil
}

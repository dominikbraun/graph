package graph

import (
	"sync"
)

type memoryStore[K comparable, T any] struct {
	lock             sync.RWMutex
	vertices         map[K]T
	vertexProperties map[K]VertexProperties
	outEdges         map[K]map[K]Edge[K] // source -> target
	inEdges          map[K]map[K]Edge[K] // target -> source
}

func newMemoryStore[K comparable, T any]() Store[K, T] {
	return &memoryStore[K, T]{
		vertices:         make(map[K]T),
		vertexProperties: make(map[K]VertexProperties),
		outEdges:         make(map[K]map[K]Edge[K]),
		inEdges:          make(map[K]map[K]Edge[K]),
	}
}

func (s *memoryStore[K, T]) AddVertex(k K, t T, p VertexProperties) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.vertices[k]; ok {
		return ErrVertexAlreadyExists
	}

	s.vertices[k] = t
	s.vertexProperties[k] = p

	return nil
}

func (s *memoryStore[K, T]) ListVertices() ([]K, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var hashes []K
	for k := range s.vertices {
		hashes = append(hashes, k)
	}

	return hashes, nil
}

func (s *memoryStore[K, T]) VertexCount() (int, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return len(s.vertices), nil
}

func (s *memoryStore[K, T]) Vertex(k K) (T, VertexProperties, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var v T
	var ok bool
	v, ok = s.vertices[k]
	if !ok {
		return v, VertexProperties{}, ErrVertexNotFound
	}

	p := s.vertexProperties[k]
	return v, p, nil
}

func (s *memoryStore[K, T]) AddEdge(sourceHash, targetHash K, edge Edge[K]) error {
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

	return nil
}

func (s *memoryStore[K, T]) RemoveEdge(sourceHash, targetHash K) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.inEdges[sourceHash], targetHash)
	delete(s.outEdges[sourceHash], targetHash)
	return nil
}

func (s *memoryStore[K, T]) Edge(sourceHash, targetHash K) (Edge[K], error) {
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

func (s *memoryStore[K, T]) GetEdgesBySource(sourceHash K) ([]Edge[K], error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	sourceEdges, ok := s.outEdges[sourceHash]
	if !ok {
		return nil, ErrEdgeNotFound
	}

	sourceEdgesArray := make([]Edge[K], 0, len(sourceEdges))
	for _, edge := range sourceEdges {
		sourceEdgesArray = append(sourceEdgesArray, edge)
	}

	return sourceEdgesArray, nil
}

func (s *memoryStore[K, T]) GetEdgesByTarget(targetHash K) ([]Edge[K], error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	targetEdges, ok := s.inEdges[targetHash]
	if !ok {
		return nil, ErrEdgeNotFound
	}

	targetEdgesArray := make([]Edge[K], 0, len(targetEdges))
	for _, edge := range targetEdges {
		targetEdgesArray = append(targetEdgesArray, edge)
	}

	return targetEdgesArray, nil
}

func (s *memoryStore[K, T]) AdjacencyMap() (map[K]map[K]Edge[K], error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	m := make(map[K]map[K]Edge[K])

	for hash := range s.vertices {
		m[hash] = make(map[K]Edge[K])
	}

	for sourceHash, targetEdges := range s.outEdges {
		m[sourceHash] = make(map[K]Edge[K])
		for targetHash, edge := range targetEdges {
			m[sourceHash][targetHash] = edge
		}
	}

	return m, nil
}

func (s *memoryStore[K, T]) PredecessorMap() (map[K]map[K]Edge[K], error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	predecessors := make(map[K]map[K]Edge[K])

	for vertexHash := range s.vertices {
		predecessors[vertexHash] = make(map[K]Edge[K])
	}

	for vertexHash, inEdges := range s.inEdges {
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

func (s *memoryStore[K, T]) ListEdges() ([]Edge[K], error) {
	res := make([]Edge[K], 0)
	for _, edges := range s.outEdges {
		for _, edge := range edges {
			res = append(res, edge)
		}
	}
	return res, nil
}

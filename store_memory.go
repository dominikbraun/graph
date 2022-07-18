package graph

import "sync"

type memoryStore[K comparable, T any] struct {
	hash     Hash[K, T]
	lock     sync.RWMutex
	vertices map[K]T
	// edges    map[K]map[K]Edge[K]
	outEdges map[K]map[K]Edge[K] // source -> target
	inEdges  map[K]map[K]Edge[K] // target -> source
}

func newMemoryStore[K comparable, T any](hash Hash[K, T]) Store[K, T] {
	return &memoryStore[K, T]{
		hash:     hash,
		vertices: make(map[K]T),
		// edges:    make(map[K]map[K]Edge[K]),
		outEdges: make(map[K]map[K]Edge[K]),
		inEdges:  make(map[K]map[K]Edge[K]),
	}
}

func (s *memoryStore[K, T]) AddVertex(t T) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.vertices[s.hash(t)] = t

	return nil
}

func (s *memoryStore[K, T]) GetAllVertexHashes() ([]K, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var hashes []K
	for k := range s.vertices {
		hashes = append(hashes, k)
	}
	return hashes, true
}

func (s *memoryStore[K, T]) GetVertex(k K) (*T, bool) { // TODO: error
	s.lock.RLock()
	defer s.lock.RUnlock()

	v, ok := s.vertices[k]
	if !ok {
		return nil, false
	}

	return &v, true
}

func (s *memoryStore[K, T]) AddEdge(sourceHash, targetHash K, edge Edge[K]) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	// if _, ok := s.edges[sourceHash]; !ok {
	// 	s.edges[sourceHash] = make(map[K]Edge[K])
	// }

	// s.edges[sourceHash][targetHash] = edge

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

func (s *memoryStore[K, T]) GetEdge(sourceHash, targetHash K) (Edge[K], bool) { // TODO: error
	s.lock.RLock()
	defer s.lock.RUnlock()

	sourceEdges, ok := s.outEdges[sourceHash]
	if !ok {
		return Edge[K]{}, false
	}

	if edge, ok := sourceEdges[targetHash]; ok {
		return edge, true
	}

	return Edge[K]{}, false
}

func (s *memoryStore[K, T]) GetEdgesBySource(sourceHash K) ([]Edge[K], bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	sourceEdges, ok := s.outEdges[sourceHash]
	if !ok {
		return nil, false
	}

	sourceEdgesArray := make([]Edge[K], 0, len(sourceEdges))
	for _, edge := range sourceEdges {
		sourceEdgesArray = append(sourceEdgesArray, edge)
	}

	return sourceEdgesArray, true
}

func (s *memoryStore[K, T]) GetEdgesByTarget(targetHash K) ([]Edge[K], bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	targetEdges, ok := s.inEdges[targetHash]
	if !ok {
		return nil, false
	}

	targetEdgesArray := make([]Edge[K], 0, len(targetEdges))
	for _, edge := range targetEdges {
		targetEdgesArray = append(targetEdgesArray, edge)
	}

	return targetEdgesArray, true
}

package graph

import "fmt"

var errNotFound = fmt.Errorf("not found")

type Store[K comparable, T any] interface {
	AddVertex(t T) error
	GetVertex(k K) (*T, bool)
	GetAllVertexHashes() ([]K, bool)
	AddEdge(sourceHash, targetHash K, edge Edge[K]) error
	GetEdge(source, target K) (Edge[K], bool)
	GetEdgeTargetHashes(source K) ([]K, bool)
	GetEdgeSourceHashes(target K) ([]K, bool)
	GetEdgeTargets(source K) (map[K]Edge[K], bool)
	GetEdgeSources(target K) (map[K]Edge[K], bool)
}

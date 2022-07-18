package graph

import "fmt"

var errNotFound = fmt.Errorf("not found")

type Store[K comparable, T any] interface {
	AddVertex(k K, t T) error
	GetVertex(k K) (*T, bool)
	GetAllVertexHashes() ([]K, bool)
	AddEdge(sourceHash, targetHash K, edge Edge[K]) error
	GetEdge(source, target K) (Edge[K], bool)
	GetEdgesBySource(source K) ([]Edge[K], bool)
	GetEdgesByTarget(target K) ([]Edge[K], bool)
}

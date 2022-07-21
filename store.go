package graph

import "fmt"

var ErrNotFound = fmt.Errorf("not found")

type Store[K comparable, T any] interface {
	AddVertex(k K, t T) error
	GetVertex(k K) (*T, error)
	ListVertices() ([]K, error)
	AddEdge(sourceHash, targetHash K, edge Edge[K]) error
	GetEdge(source, target K) (*Edge[K], error)
	GetEdgesBySource(source K) ([]Edge[K], error)
	GetEdgesByTarget(target K) ([]Edge[K], error)
}

package graph

type Store[K comparable, T any] interface {
	AddVertex(hash K, value T, properties VertexProperties) error
	Vertex(hash K) (T, VertexProperties, error)
	ListVertices() ([]K, error)
	CountVertices() (int, error)
	AddEdge(sourceHash, targetHash K, edge Edge[K]) error
	RemoveEdge(sourceHash, targetHash K) error
	Edge(sourceHash, targetHash K) (Edge[K], error)
	GetEdgesBySource(sourceHash K) ([]Edge[K], error)
	GetEdgesByTarget(targetHash K) ([]Edge[K], error)
	AdjacencyMap() (map[K]map[K]Edge[K], error)
	PredecessorMap() (map[K]map[K]Edge[K], error)
}

package graph

type Store[K comparable, T any] interface {
	AddVertex(hash K, value T, properties VertexProperties) error
	Vertex(hash K) (T, VertexProperties, error)
	ListVertices() ([]K, error)
	VertexCount() (int, error)
	AddEdge(sourceHash, targetHash K, edge Edge[K]) error
	RemoveEdge(sourceHash, targetHash K) error
	Edge(sourceHash, targetHash K) (Edge[K], error)
	ListEdges() ([]Edge[K], error)
}

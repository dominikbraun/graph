// Package graph is a library for creating generic graph data structures and
// modifying, analyzing, and visualizing them.
//
// # Hashes
//
// A graph consists of vertices of type T, which are identified by a hash value
// of type K. The hash value for a given vertex is obtained using the hashing
// function passed to [New]. A hashing function takes a T and returns a K.
//
// For primitive types like integers, you may use a predefined hashing function
// such as [IntHash] – a function that takes an integer and uses that integer as
// the hash value at the same time:
//
//	g := graph.New(graph.IntHash)
//
// For storing custom data types, you need to provide your own hashing function.
// This example takes a City instance and returns its name as the hash value:
//
//	cityHash := func(c City) string {
//		return c.Name
//	}
//
// Creating a graph using this hashing function will yield a graph of vertices
// of type City identified by hash values of type string.
//
//	g := graph.New(cityHash)
//
// # Operations
//
// Adding vertices to a graph of integers is simple. [graph.Graph.AddVertex]
// takes a vertex and adds it to the graph.
//
//	g := graph.New(graph.IntHash)
//
//	_ = g.AddVertex(1)
//	_ = g.AddVertex(2)
//
// Most functions accept and return only hash values instead of entire instances
// of the vertex type T. For example, [graph.Graph.AddEdge] creates an edge
// between two vertices and accepts the hash values of those vertices. Because
// this graph uses the [IntHash] hashing function, the vertex values and hash
// values are the same.
//
//	_ = g.AddEdge(1, 2)
//
// All operations that modify the graph itself are methods of [Graph]. All other
// operations are top-level functions of by this library.
//
// For detailed usage examples, take a look at the README.
package graph

import (
	"errors"
	"fmt"
)

var (
	ErrVertexNotFound      = errors.New("vertex not found")
	ErrVertexAlreadyExists = errors.New("vertex already exists")
	ErrEdgeNotFound        = errors.New("edge not found")
	ErrEdgeAlreadyExists   = errors.New("edge already exists")
	ErrEdgeCreatesCycle    = errors.New("edge would create a cycle")
	ErrVertexHasEdges      = errors.New("vertex has edges")
)

// Graph represents a generic graph data structure consisting of vertices of
// type T identified by a hash of type K.
type Graph[K comparable, T any] interface {
	// Traits returns the graph's traits. Those traits must be set when creating
	// a graph using New.
	Traits() *Traits

	// AddVertex creates a new vertex in the graph. If the vertex already exists
	// in the graph, ErrVertexAlreadyExists will be returned.
	//
	// AddVertex accepts a variety of functional options to set further edge
	// details such as the weight or an attribute:
	//
	//	_ = graph.AddVertex("A", "B", graph.VertexWeight(4), graph.VertexAttribute("label", "my-label"))
	//
	AddVertex(value T, options ...func(*VertexProperties)) error

	// AddVerticesFrom adds all vertices along with their properties from the
	// given graph to the receiving graph.
	//
	// All vertices will be added until an error occurs. If one of the vertices
	// already exists, ErrVertexAlreadyExists will be returned.
	AddVerticesFrom(g Graph[K, T]) error

	// Vertex returns the vertex with the given hash or ErrVertexNotFound if it
	// doesn't exist.
	Vertex(hash K) (T, error)

	// VertexWithProperties returns the vertex with the given hash along with
	// its properties or ErrVertexNotFound if it doesn't exist.
	VertexWithProperties(hash K) (T, VertexProperties, error)

	// RemoveVertex removes the vertex with the given hash value from the graph.
	//
	// The vertex is not allowed to have edges and thus must be disconnected.
	// Potential edges must be removed first. Otherwise, ErrVertexHasEdges will
	// be returned. If the vertex doesn't exist, ErrVertexNotFound is returned.
	RemoveVertex(hash K) error

	// AddEdge creates an edge between the source and the target vertex.
	//
	// If either vertex cannot be found, ErrVertexNotFound will be returned. If
	// the edge already exists, ErrEdgeAlreadyExists will be returned. If cycle
	// prevention has been activated using PreventCycles and if adding the edge
	// would create a cycle, ErrEdgeCreatesCycle will be returned.
	//
	// AddEdge accepts functional options to set further edge properties such as
	// the weight or an attribute:
	//
	//	_ = g.AddEdge("A", "B", graph.EdgeWeight(4), graph.EdgeAttribute("label", "my-label"))
	//
	AddEdge(sourceHash, targetHash K, options ...func(*EdgeProperties)) error

	// AddEdgesFrom adds all edges along with their properties from the given
	// graph to the receiving graph.
	//
	// All vertices that the edges are joining have to exist already. If needed,
	// these vertices can be added using AddVerticesFrom first. Depending on the
	// situation, it also might make sense to clone the entire original graph.
	AddEdgesFrom(g Graph[K, T]) error

	// Edge returns the edge joining two given vertices or ErrEdgeNotFound if
	// the edge doesn't exist. In an undirected graph, an edge with swapped
	// source and target vertices does match.
	Edge(sourceHash, targetHash K) (Edge[T], error)

	// Edges returns a slice of all edges in the graph. These edges are of type
	// Edge[K] and hence will contain the vertex hashes, not the vertex values.
	Edges() ([]Edge[K], error)

	// UpdateEdge updates the edge joining the two given vertices with the data
	// provided in the given functional options. Valid functional options are:
	// - EdgeWeight: Sets a new weight for the edge properties.
	// - EdgeAttribute: Adds a new attribute to the edge properties.
	// - EdgeAttributes: Sets a new attributes map for the edge properties.
	// - EdgeData: Sets a new Data field for the edge properties.
	//
	// UpdateEdge accepts the same functional options as AddEdge. For example,
	// setting the weight of an edge (A,B) to 10 would look as follows:
	//
	//	_ = g.UpdateEdge("A", "B", graph.EdgeWeight(10))
	//
	// Removing a particular edge attribute is not possible at the moment. A
	// workaround is to create a new map without the respective element and
	// overwrite the existing attributes using the EdgeAttributes option.
	UpdateEdge(source, target K, options ...func(properties *EdgeProperties)) error

	// RemoveEdge removes the edge between the given source and target vertices.
	// If the edge cannot be found, ErrEdgeNotFound will be returned.
	RemoveEdge(source, target K) error

	// AdjacencyMap computes an adjacency map with all vertices in the graph.
	//
	// There is an entry for each vertex. Each of those entries is another map
	// whose keys are the hash values of the adjacent vertices. The value is an
	// Edge instance that stores the source and target hash values along with
	// the edge metadata.
	//
	// For a directed graph with two edges AB and AC, AdjacencyMap would return
	// the following map:
	//
	//	map[string]map[string]Edge[string]{
	//		"A": map[string]Edge[string]{
	//			"B": {Source: "A", Target: "B"},
	//			"C": {Source: "A", Target: "C"},
	//		},
	//		"B": map[string]Edge[string]{},
	//		"C": map[string]Edge[string]{},
	//	}
	//
	// This design makes AdjacencyMap suitable for a wide variety of algorithms.
	AdjacencyMap() (map[K]map[K]Edge[K], error)

	// PredecessorMap computes a predecessor map with all vertices in the graph.
	//
	// It has the same map layout and does the same thing as AdjacencyMap, but
	// for ingoing instead of outgoing edges of each vertex.
	//
	// For a directed graph with two edges AB and AC, PredecessorMap would
	// return the following map:
	//
	//	map[string]map[string]Edge[string]{
	//		"A": map[string]Edge[string]{},
	//		"B": map[string]Edge[string]{
	//			"A": {Source: "A", Target: "B"},
	//		},
	//		"C": map[string]Edge[string]{
	//			"A": {Source: "A", Target: "C"},
	//		},
	//	}
	//
	// For an undirected graph, PredecessorMap is the same as AdjacencyMap. This
	// is because there is no distinction between "outgoing" and "ingoing" edges
	// in an undirected graph.
	PredecessorMap() (map[K]map[K]Edge[K], error)

	// Clone creates a deep copy of the graph and returns that cloned graph.
	//
	// The cloned graph will use the default in-memory store for storing the
	// vertices and edges. If you want to utilize a custom store instead, create
	// a new graph using NewWithStore and use AddVerticesFrom and AddEdgesFrom.
	Clone() (Graph[K, T], error)

	// Order returns the number of vertices in the graph.
	Order() (int, error)

	// Size returns the number of edges in the graph.
	Size() (int, error)
}

// Edge represents an edge that joins two vertices. Even though these edges are
// always referred to as source and target, whether the graph is directed or not
// is determined by its traits.
type Edge[T any] struct {
	Source     T
	Target     T
	Properties EdgeProperties
}

// EdgeProperties represents a set of properties that each edge possesses. They
// can be set when adding a new edge using the corresponding functional options:
//
//	g.AddEdge("A", "B", graph.EdgeWeight(2), graph.EdgeAttribute("color", "red"))
//
// The example above will create an edge with a weight of 2 and an attribute
// "color" with value "red".
type EdgeProperties struct {
	Attributes map[string]string
	Weight     int
	Data       any
}

// Hash is a hashing function that takes a vertex of type T and returns a hash
// value of type K.
//
// Every graph has a hashing function and uses that function to retrieve the
// hash values of its vertices. You can either use one of the predefined hashing
// functions or provide your own one for custom data types:
//
//	cityHash := func(c City) string {
//		return c.Name
//	}
//
// The cityHash function returns the city name as a hash value. The types of T
// and K, in this case City and string, also define the types of the graph.
type Hash[K comparable, T any] func(T) K

// New creates a new graph with vertices of type T, identified by hash values of
// type K. These hash values will be obtained using the provided hash function.
//
// The graph will use the default in-memory store for persisting vertices and
// edges. To use a different [Store], use [NewWithStore].
func New[K comparable, T any](hash Hash[K, T], options ...func(*Traits)) Graph[K, T] {
	return NewWithStore(hash, newMemoryStore[K, T](), options...)
}

// NewWithStore creates a new graph same as [New] but uses the provided store
// instead of the default memory store.
func NewWithStore[K comparable, T any](hash Hash[K, T], store Store[K, T], options ...func(*Traits)) Graph[K, T] {
	var p Traits

	for _, option := range options {
		option(&p)
	}

	if p.IsDirected {
		return newDirected(hash, &p, store)
	}

	return newUndirected(hash, &p, store)
}

// NewLike creates a graph that is "like" the given graph: It has the same type,
// the same hashing function, and the same traits. The new graph is independent
// of the original graph and uses the default in-memory storage.
//
//	g := graph.New(graph.IntHash, graph.Directed())
//	h := graph.NewLike(g)
//
// In the example above, h is a new directed graph of integers derived from g.
func NewLike[K comparable, T any](g Graph[K, T]) Graph[K, T] {
	copyTraits := func(t *Traits) {
		t.IsDirected = g.Traits().IsDirected
		t.IsAcyclic = g.Traits().IsAcyclic
		t.IsWeighted = g.Traits().IsWeighted
		t.IsRooted = g.Traits().IsRooted
		t.PreventCycles = g.Traits().PreventCycles
	}

	gr := tograph(g)

	return New(gr.hash, copyTraits)
}

// StringHash is a hashing function that accepts a string and uses that exact
// string as a hash value. Using it as Hash will yield a Graph[string, string].
func StringHash(v string) string {
	return v
}

// IntHash is a hashing function that accepts an integer and uses that exact
// integer as a hash value. Using it as Hash will yield a Graph[int, int].
func IntHash(v int) int {
	return v
}

// EdgeWeight returns a function that sets the weight of an edge to the given
// weight. This is a functional option for the [graph.Graph.Edge] and
// [graph.Graph.AddEdge] methods.
func EdgeWeight(weight int) func(*EdgeProperties) {
	return func(e *EdgeProperties) {
		e.Weight = weight
	}
}

// EdgeAttribute returns a function that adds the given key-value pair to the
// attributes of an edge. This is a functional option for the [graph.Graph.Edge]
// and [graph.Graph.AddEdge] methods.
func EdgeAttribute(key, value string) func(*EdgeProperties) {
	return func(e *EdgeProperties) {
		e.Attributes[key] = value
	}
}

// EdgeAttributes returns a function that sets the given map as the attributes
// of an edge. This is a functional option for the [graph.Graph.AddEdge] and
// [graph.Graph.UpdateEdge] methods.
func EdgeAttributes(attributes map[string]string) func(*EdgeProperties) {
	return func(e *EdgeProperties) {
		e.Attributes = attributes
	}
}

// EdgeData returns a function that sets the data of an edge to the given value.
// This is a functional option for the [graph.Graph.Edge] and
// [graph.Graph.AddEdge] methods.
func EdgeData(data any) func(*EdgeProperties) {
	return func(e *EdgeProperties) {
		e.Data = data
	}
}

// VertexProperties represents a set of properties that each vertex has. They
// can be set when adding a vertex using the corresponding functional options:
//
//	_ = g.AddVertex("A", "B", graph.VertexWeight(2), graph.VertexAttribute("color", "red"))
//
// The example above will create a vertex with a weight of 2 and an attribute
// "color" with value "red".
type VertexProperties struct {
	Attributes map[string]string
	Weight     int
}

// VertexWeight returns a function that sets the weight of a vertex to the given
// weight. This is a functional option for the [graph.Graph.Vertex] and
// [graph.Graph.AddVertex] methods.
func VertexWeight(weight int) func(*VertexProperties) {
	return func(e *VertexProperties) {
		e.Weight = weight
	}
}

// VertexAttribute returns a function that adds the given key-value pair to the
// vertex attributes. This is a functional option for the [graph.Graph.Vertex]
// and [graph.Graph.AddVertex] methods.
func VertexAttribute(key, value string) func(*VertexProperties) {
	return func(e *VertexProperties) {
		e.Attributes[key] = value
	}
}

// VertexAttributes returns a function that sets the given map as the attributes
// of a vertex. This is a functional option for the [graph.Graph.AddVertex] methods.
func VertexAttributes(attributes map[string]string) func(*VertexProperties) {
	return func(e *VertexProperties) {
		e.Attributes = attributes
	}
}

// graph is a partial Graph[K, T] implementation where all the operations are
// common to both directed and undirected graphs.
type graph[K comparable, T any] struct {
	hash   Hash[K, T]
	traits *Traits
	store  Store[K, T]
}

func (g *graph[K, T]) AddEdge(sourceHash, targetHash K, options ...func(*EdgeProperties)) error {
	edge, err := g.addEdgeWithOptions(sourceHash, targetHash, options...)
	if err != nil || g.traits.IsDirected { // if we're directed or in error, we're done
		return err
	}

	// add reverse edge of undirected
	rEdge := Edge[K]{
		Source: edge.Target,
		Target: edge.Source,
		Properties: EdgeProperties{
			Weight:     edge.Properties.Weight,
			Attributes: edge.Properties.Attributes,
			Data:       edge.Properties.Data,
		},
	}

	return g.store.AddEdge(targetHash, sourceHash, rEdge)
}

func (g *graph[K, T]) addEdge(sourceHash, targetHash K, edge Edge[K]) error {
	return g.store.AddEdge(sourceHash, targetHash, edge)
}

func (g *graph[K, T]) addEdgeWithOptions(sourceHash, targetHash K, options ...func(*EdgeProperties)) (Edge[K], error) {
	if _, _, err := g.store.Vertex(sourceHash); err != nil {
		return Edge[K]{}, fmt.Errorf("could not find source vertex with hash %v: %w", sourceHash, err)
	}

	if _, _, err := g.store.Vertex(targetHash); err != nil {
		return Edge[K]{}, fmt.Errorf("could not find target vertex with hash %v: %w", targetHash, err)
	}

	//nolint:govet // False positive.
	if _, err := g.Edge(sourceHash, targetHash); !errors.Is(err, ErrEdgeNotFound) {
		return Edge[K]{}, ErrEdgeAlreadyExists
	}

	// If the user opted in to preventing cycles, run a cycle check.
	if g.traits.PreventCycles {
		createsCycle, err := g.createsCycle(sourceHash, targetHash)
		if err != nil {
			return Edge[K]{}, fmt.Errorf("check for cycles: %w", err)
		}
		if createsCycle {
			return Edge[K]{}, ErrEdgeCreatesCycle
		}
	}

	edge := Edge[K]{
		Source: sourceHash,
		Target: targetHash,
		Properties: EdgeProperties{
			Attributes: make(map[string]string),
		},
	}

	for _, option := range options {
		option(&edge.Properties)
	}

	if err := g.store.AddEdge(sourceHash, targetHash, edge); err != nil {
		return Edge[K]{}, fmt.Errorf("failed to add edge: %w", err)
	}

	return edge, nil
}

func (g *graph[K, T]) addEdgesFrom(other *graph[K, T]) error {
	edges, err := other.Edges()
	if err != nil {
		return fmt.Errorf("failed to get edges: %w", err)
	}

	for _, edge := range edges {
		if err := g.AddEdge(copyEdge(edge)); err != nil {
			return fmt.Errorf("failed to add (%v, %v): %w", edge.Source, edge.Target, err)
		}
	}

	return nil
}

func (g *graph[K, T]) clone() (*graph[K, T], error) {
	traits := &Traits{
		IsDirected:    g.traits.IsDirected,
		IsAcyclic:     g.traits.IsAcyclic,
		IsWeighted:    g.traits.IsWeighted,
		IsRooted:      g.traits.IsRooted,
		PreventCycles: g.traits.PreventCycles,
	}

	graph := &graph[K, T]{
		hash:   g.hash,
		traits: traits,
		store:  newMemoryStore[K, T](),
	}

	if err := graph.addVerticesFrom(g); err != nil {
		return nil, fmt.Errorf("failed to add vertices: %w", err)
	}

	if err := graph.addEdgesFrom(g); err != nil {
		return nil, fmt.Errorf("failed to add edges: %w", err)
	}

	return graph, nil
}
func (g *graph[K, T]) Edge(sourceHash, targetHash K) (Edge[T], error) {
	edge, err := g.store.Edge(sourceHash, targetHash)
	if !g.traits.IsDirected {
		// In an undirected graph, since multigraphs aren't supported, the edge AB
		// is the same as BA. Therefore, if source[target] cannot be found, this
		// function also looks for target[source].
		if errors.Is(err, ErrEdgeNotFound) {
			edge, err = g.store.Edge(targetHash, sourceHash)
		}
	}

	if err != nil {
		return Edge[T]{}, err
	}

	sourceVertex, _, err := g.store.Vertex(sourceHash)
	if err != nil {
		return Edge[T]{}, err
	}

	targetVertex, _, err := g.store.Vertex(targetHash)
	if err != nil {
		return Edge[T]{}, err
	}

	return Edge[T]{
		Source: sourceVertex,
		Target: targetVertex,
		Properties: EdgeProperties{
			Weight:     edge.Properties.Weight,
			Attributes: edge.Properties.Attributes,
			Data:       edge.Properties.Data,
		},
	}, nil
}

type tuple[K comparable] struct {
	source, target K
}

func (g *graph[K, T]) Edges() ([]Edge[K], error) {
	storedEdges, err := g.store.ListEdges()
	if g.traits.IsDirected {
		return storedEdges, err
	}
	// An undirected graph creates each edge twice internally: The edge (A,B) is
	// stored both as (A,B) and (B,A). The Edges method is supposed to return
	// one of these two edges, because from an outside perspective, it only is
	// a single edge.
	//
	// To achieve this, Edges keeps track of already-added edges. For each edge,
	// it also checks if the reversed edge has already been added - e.g., for
	// an edge (A,B), Edges checks if the edge has been added as (B,A).
	//
	// These reversed edges are built as a custom tuple type, which is then used
	// as a map key for access in O(1) time. It looks scarier than it is.
	edges := make([]Edge[K], 0, len(storedEdges)/2)

	added := make(map[tuple[K]]struct{})

	for _, storedEdge := range storedEdges {
		reversedEdge := tuple[K]{
			source: storedEdge.Target,
			target: storedEdge.Source,
		}
		if _, ok := added[reversedEdge]; ok {
			continue
		}

		edges = append(edges, storedEdge)

		addedEdge := tuple[K]{
			source: storedEdge.Source,
			target: storedEdge.Target,
		}

		added[addedEdge] = struct{}{}
	}

	return edges, nil
}

func (g *graph[K, T]) RemoveEdge(source, target K) error {
	if _, err := g.Edge(source, target); err != nil {
		return err
	}

	if err := g.store.RemoveEdge(source, target); err != nil {
		return fmt.Errorf("failed to remove edge from %v to %v: %w", source, target, err)
	}

	return nil
}

func (g *graph[K, T]) updateEdge(source, target K, options ...func(properties *EdgeProperties)) (Edge[K], error) {
	existingEdge, err := g.store.Edge(source, target)
	if err != nil {
		return Edge[K]{}, err
	}

	for _, option := range options {
		option(&existingEdge.Properties)
	}

	return existingEdge, g.store.UpdateEdge(source, target, existingEdge)
}

func (g *graph[K, T]) Traits() *Traits {
	return g.traits
}

func (g *graph[K, T]) AddVertex(value T, options ...func(*VertexProperties)) error {
	hash := g.hash(value)
	properties := VertexProperties{
		Weight:     0,
		Attributes: make(map[string]string),
	}

	for _, option := range options {
		option(&properties)
	}

	return g.store.AddVertex(hash, value, properties)

}

func (g *graph[K, T]) addVerticesFrom(other *graph[K, T]) error {
	adjacencyMap, err := other.AdjacencyMap()
	if err != nil {
		return fmt.Errorf("failed to get adjacency map: %w", err)
	}

	for hash := range adjacencyMap {
		vertex, properties, err := other.VertexWithProperties(hash)
		if err != nil {
			return fmt.Errorf("failed to get vertex %v: %w", hash, err)
		}

		if err = g.AddVertex(vertex, copyVertexProperties(properties)); err != nil {
			return fmt.Errorf("failed to add vertex %v: %w", hash, err)
		}
	}

	return nil
}

func (g *graph[K, T]) Vertex(hash K) (T, error) {
	vertex, _, err := g.store.Vertex(hash)
	return vertex, err
}

func (g *graph[K, T]) VertexWithProperties(hash K) (T, VertexProperties, error) {
	vertex, properties, err := g.store.Vertex(hash)
	if err != nil {
		return vertex, VertexProperties{}, err
	}

	return vertex, properties, nil
}

func (g *graph[K, T]) RemoveVertex(hash K) error {
	return g.store.RemoveVertex(hash)
}

func (g *graph[K, T]) AdjacencyMap() (map[K]map[K]Edge[K], error) {
	vertices, err := g.store.ListVertices()
	if err != nil {
		return nil, fmt.Errorf("failed to list vertices: %w", err)
	}

	edges, err := g.store.ListEdges()
	if err != nil {
		return nil, fmt.Errorf("failed to list edges: %w", err)
	}

	m := make(map[K]map[K]Edge[K], len(vertices))

	for _, vertex := range vertices {
		m[vertex] = make(map[K]Edge[K])
	}

	for _, edge := range edges {
		m[edge.Source][edge.Target] = edge
	}

	return m, nil
}

func (g *graph[K, T]) PredecessorMap() (map[K]map[K]Edge[K], error) {
	if !g.traits.IsDirected {
		return g.AdjacencyMap()
	}

	vertices, err := g.store.ListVertices()
	if err != nil {
		return nil, fmt.Errorf("failed to list vertices: %w", err)
	}

	edges, err := g.store.ListEdges()
	if err != nil {
		return nil, fmt.Errorf("failed to list edges: %w", err)
	}

	m := make(map[K]map[K]Edge[K], len(vertices))

	for _, vertex := range vertices {
		m[vertex] = make(map[K]Edge[K])
	}

	for _, edge := range edges {
		if _, ok := m[edge.Target]; !ok {
			m[edge.Target] = make(map[K]Edge[K])
		}
		m[edge.Target][edge.Source] = edge
	}

	return m, nil
}

func (g *graph[K, T]) Order() (int, error) {
	return g.store.VertexCount()
}

func (g *graph[K, T]) Size() (int, error) {
	edgeCount, err := g.store.EdgeCount()

	if g.traits.IsDirected {
		return edgeCount, err
	}
	// Divide by 2 since every add edge operation on undirected graph is counted
	// twice.
	return edgeCount / 2, err
}

func (g *graph[K, T]) createsCycle(source, target K) (bool, error) {
	// If the underlying store implements CreatesCycle, use that fast path.
	if cc, ok := g.store.(interface {
		CreatesCycle(source, target K) (bool, error)
	}); ok {
		return cc.CreatesCycle(source, target)
	}

	// Slow path.
	return CreatesCycle[K, T](fromGraph(g), source, target)
}

func (g *graph[K, T]) edgesAreEqual(a, b Edge[T]) bool {
	aSourceHash := g.hash(a.Source)
	aTargetHash := g.hash(a.Target)
	bSourceHash := g.hash(b.Source)
	bTargetHash := g.hash(b.Target)

	if aSourceHash == bSourceHash && aTargetHash == bTargetHash {
		return true
	}

	if !g.traits.IsDirected {
		return aSourceHash == bTargetHash && aTargetHash == bSourceHash
	}

	return false
}

// copyEdge returns an argument list suitable for the Graph.AddEdge method. This
// argument list is derived from the given edge, hence the name copyEdge.
//
// The last argument is a custom functional option that sets the edge properties
// to the properties of the original edge.
func copyEdge[K comparable](edge Edge[K]) (K, K, func(properties *EdgeProperties)) {
	copyProperties := func(p *EdgeProperties) {
		for k, v := range edge.Properties.Attributes {
			p.Attributes[k] = v
		}
		p.Weight = edge.Properties.Weight
		p.Data = edge.Properties.Data
	}

	return edge.Source, edge.Target, copyProperties
}

// tograph converts a Graph interface to a graph instance.
func tograph[K comparable, T any](g Graph[K, T]) graph[K, T] {
	var other graph[K, T]
	if g.Traits().IsDirected {
		other = g.(*directed[K, T]).graph
	} else {
		other = g.(*undirected[K, T]).graph
	}
	return other
}

// fromGraph converts a graph instance to a Graph interface.
func fromGraph[K comparable, T any](g *graph[K, T]) Graph[K, T] {
	var gr Graph[K, T]
	if g.Traits().IsDirected {
		gr = &directed[K, T]{graph: *g}
	} else {
		gr = &undirected[K, T]{graph: *g}
	}
	return gr
}

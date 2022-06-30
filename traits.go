package graph

type traits struct {
	isDirected bool
	isAcyclic  bool
	isWeighted bool
	isRooted   bool
}

// Directed creates a directed graph. This has implications on graph traversal and the order of
// arguments of the Edge and EdgeByHashes functions.
func Directed() func(*traits) {
	return func(t *traits) {
		t.isDirected = true
	}
}

// Acyclic creates an acyclic graph, which won't allow adding directed edges resulting in a cycle.
func Acyclic() func(*traits) {
	return func(t *traits) {
		t.isAcyclic = true
	}
}

// Weighted creates a weighted graph. To set weights, use the WeightedEdge and WeightedEdgeByHashes
// functions.
func Weighted() func(*traits) {
	return func(t *traits) {
		t.isWeighted = true
	}
}

// Rooted creates a rooted graph. This is particularly common for building tree data structures.
func Rooted() func(*traits) {
	return func(t *traits) {
		t.isRooted = true
	}
}

// Tree is an alias for Acyclic and Rooted, since most trees in Computer Science are rooted trees.
func Tree() func(*traits) {
	return func(t *traits) {
		Acyclic()(t)
		Rooted()(t)
	}
}

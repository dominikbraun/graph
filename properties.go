package graph

type properties struct {
	isDirected bool
	isAcyclic  bool
	isWeighted bool
	isRooted   bool
}

// Directed creates a directed graph. This has implications on graph traversal and the order of
// arguments of the Edge and EdgeByHashes functions.
func Directed() func(*properties) {
	return func(p *properties) {
		p.isDirected = true
	}
}

// Acyclic creates an acyclic graph, which won't allow adding directed edges resulting in a cycle.
func Acyclic() func(*properties) {
	return func(p *properties) {
		p.isAcyclic = true
	}
}

// Weighted creates a weighted graph. To set weights, use the WeightedEdge and WeightedEdgeByHashes
// functions.
func Weighted() func(*properties) {
	return func(p *properties) {
		p.isWeighted = true
	}
}

// Rooted creates a rooted graph. This is particularly common for building tree data structures.
func Rooted() func(*properties) {
	return func(p *properties) {
		p.isRooted = true
	}
}

// Tree is an alias for Acyclic and Rooted, since most trees in Computer Science are rooted trees.
func Tree() func(*properties) {
	return func(p *properties) {
		Acyclic()(p)
		Rooted()(p)
	}
}

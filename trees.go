package graph

type unionFind[K comparable] struct {
	parents map[K]K
}

func newUnionFind[K comparable](initialVertices []K) *unionFind[K] {
	u := &unionFind[K]{
		parents: make(map[K]K),
	}

	for _, vertex := range initialVertices {
		u.parents[vertex] = vertex
	}

	return u
}

func (u *unionFind[K]) union(vertex1, vertex2 K) {
	root1 := u.find(vertex1)
	root2 := u.find(vertex2)

	if root1 == root2 {
		return
	}

	u.parents[root2] = root1
}

func (u *unionFind[K]) find(vertex K) K {
	root := vertex

	for u.parents[root] != root {
		root = u.parents[root]
	}

	// Perform a path compression in order to optimize of future find calls.
	current := vertex

	for u.parents[current] != root {
		parent := u.parents[vertex]
		u.parents[vertex] = root
		current = parent
	}

	return root
}

package graph

type directed[K comparable, T any] struct {
	graph[K, T]
}

func newDirected[K comparable, T any](hash Hash[K, T], traits *Traits, store Store[K, T]) *directed[K, T] {
	traits.IsDirected = true
	return &directed[K, T]{
		graph: graph[K, T]{
			hash:   hash,
			traits: traits,
			store:  store,
		},
	}
}

func (d *directed[K, T]) AddVerticesFrom(g Graph[K, T]) error {
	other := tograph(g)
	return d.addVerticesFrom(&other)

}

func (d *directed[K, T]) AddEdgesFrom(g Graph[K, T]) error {
	other := tograph(g)
	return d.addEdgesFrom(&other)
}

func (d *directed[K, T]) UpdateEdge(source, target K, options ...func(properties *EdgeProperties)) error {
	_, err := d.updateEdge(source, target, options...)

	return err
}

func (d *directed[K, T]) Clone() (Graph[K, T], error) {
	gclone, err := d.clone()
	if err != nil {
		return nil, err
	}

	dir := &directed[K, T]{
		graph: *gclone,
	}
	return dir, nil
}

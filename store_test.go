package graph

import (
	"errors"
	"testing"
)

func TestRemoveEdge(t *testing.T) {
	noerr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	build := func(nodes []string, edges [][]string) Store[string, string] {
		store := NewMemoryStore[string, string]()
		for _, n := range nodes {
			noerr(store.AddVertex(n, n, VertexProperties{}))
		}
		for _, e := range edges {
			noerr(store.AddEdge(e[0], e[1], Edge[string]{
				Source: e[0],
				Target: e[1],
			}))
		}
		return store
	}

	t.Run("normal remove", func(t *testing.T) {
		g := build([]string{
			"a", "b", "c",
		}, [][]string{
			{"a", "b"},
			{"a", "c"},
		})
		noerr(g.RemoveEdge("a", "b"))
		noerr(g.RemoveEdge("a", "c"))
		noerr(g.RemoveVertex("a"))
	})
	t.Run("remove edge has in-edges", func(t *testing.T) {
		g := build([]string{
			"a", "b", "c",
		}, [][]string{
			{"a", "b"},
			{"a", "c"},
		})
		if err := g.RemoveVertex("b"); !errors.Is(err, ErrVertexHasEdges) {
			t.Fail()
		}
	})
	t.Run("remove edge has out-edges", func(t *testing.T) {
		g := build([]string{
			"a", "b", "c",
		}, [][]string{
			{"a", "b"},
			{"a", "c"},
		})
		if err := g.RemoveVertex("a"); !errors.Is(err, ErrVertexHasEdges) {
			t.Fail()
		}
	})
}

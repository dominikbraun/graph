package draw

import (
	"testing"

	"github.com/dominikbraun/graph"
)

func TestGenerateDOT(t *testing.T) {
	tests := map[string]struct {
		graph    graph.Graph[int, int]
		vertices []int
		edges    []graph.Edge[int]
		expected description
	}{
		"3-vertex directed graph": {
			graph:    graph.New(graph.IntHash, graph.Directed()),
			vertices: []int{1, 2, 3},
			edges: []graph.Edge[int]{
				{Source: 1, Target: 2},
				{Source: 1, Target: 3},
			},
			expected: description{
				GraphType:    "digraph",
				EdgeOperator: "->",
				Statements: []statement{
					{Source: 1, Target: 2},
					{Source: 1, Target: 3},
					{Source: 2},
					{Source: 3},
				},
			},
		},
	}

	for name, test := range tests {
		for _, vertex := range test.vertices {
			test.graph.Vertex(vertex)
		}

		for _, edge := range test.edges {
			if err := test.graph.Edge(edge.Source, edge.Target); err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		description := generateDOT(test.graph)

		if description.GraphType != test.expected.GraphType {
			t.Errorf("%s: graph type expectancy doesn't match: expected %v, got %v", name, test.expected.GraphType, description.GraphType)
		}

		if description.EdgeOperator != test.expected.EdgeOperator {
			t.Errorf("%s: edge operator expectancy doesn't match: expected %v, got %v", name, test.expected.EdgeOperator, description.EdgeOperator)
		}

		if !slicesAreEqual(description.Statements, test.expected.Statements, statementsAreEqual) {
			t.Errorf("%s: statements expectancy doesn't match: expected %v, got %v", name, test.expected.Statements, description.Statements)
		}
	}
}

func TestRenderDOT(t *testing.T) {
}

func slicesAreEqual[T any](a []T, b []T, equals func(a, b T) bool) bool {
	if len(a) != len(b) {
		return false
	}

	for _, aValue := range a {
		found := false
		for _, bValue := range b {
			if equals(aValue, bValue) {
				found = true
			}
		}
		if !found {
			return false
		}
	}

	return true
}

func statementsAreEqual(a, b statement) bool {
	return a.Source == b.Source &&
		a.Target == b.Target &&
		a.Weight == b.Weight &&
		a.Label == b.Label
}

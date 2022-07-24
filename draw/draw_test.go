package draw

import (
	"bytes"
	"strings"
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
	tests := map[string]struct {
		description description
		expected    string
	}{
		"3-vertex directed graph": {
			description: description{
				GraphType:    "digraph",
				EdgeOperator: "->",
				Statements: []statement{
					{Source: 1, Target: 2},
					{Source: 1, Target: 3},
					{Source: 2},
					{Source: 3},
				},
			},
			expected: `strict digraph {
				1 -> 2 [ weight=0 ];
				1 -> 3 [ weight=0 ];
				2 ;
				3 ;
			}`,
		},
		"custom edge attributes": {
			description: description{
				GraphType:    "digraph",
				EdgeOperator: "->",
				Statements: []statement{
					{
						Source: 1,
						Target: 2,
						Attributes: map[string]string{
							"color": "red",
						},
					},
					{
						Source: 1,
						Target: 3,
						Attributes: map[string]string{
							"color": "blue",
						},
					},
					{Source: 2},
					{Source: 3},
				},
			},
			expected: `strict digraph {
				1 -> 2 [ color="red", weight=0 ];
				1 -> 3 [ color="blue", weight=0 ];
				2 ;
				3 ;
			}`,
		},
	}

	for name, test := range tests {
		buf := new(bytes.Buffer)
		renderDOT(buf, test.description)

		output := normalizeOutput(buf.String())
		expected := normalizeOutput(test.expected)

		if output != expected {
			t.Errorf("%s: DOT output expectancy doesn't match: expected %v, got %v", name, expected, output)
		}
	}
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

func normalizeOutput(output string) string {
	replacer := strings.NewReplacer(" ", "", "\n", "", "\t", "")

	return replacer.Replace(output)
}

func statementsAreEqual(a, b statement) bool {
	if len(a.Attributes) != len(b.Attributes) {
		return false
	}

	for aKey, aValue := range a.Attributes {
		bValue, ok := b.Attributes[aKey]
		if !ok {
			return false
		}
		if aValue != bValue {
			return false
		}
	}

	return a.Source == b.Source &&
		a.Target == b.Target &&
		a.Weight == b.Weight
}

package draw

import (
	"bytes"
	"strings"
	"testing"

	"github.com/dominikbraun/graph"
)

func TestGenerateDOT(t *testing.T) {
	tests := map[string]struct {
		graph    graph.Graph[string, string]
		vertices []string
		edges    []graph.Edge[string]
		expected description
	}{
		"3-vertex directed graph": {
			graph:    graph.New(graph.StringHash, graph.Directed()),
			vertices: []string{"1", "2", "3"},
			edges: []graph.Edge[string]{
				{Source: "1", Target: "2"},
				{Source: "1", Target: "3"},
			},
			expected: description{
				GraphType:    "digraph",
				EdgeOperator: "->",
				Statements: []statement{
					{Source: "1", Target: "2"},
					{Source: "1", Target: "3"},
					{Source: "2"},
					{Source: "3"},
				},
			},
		},
		"3-vertex directed, weighted graph with weights and attributes": {
			graph:    graph.New(graph.StringHash, graph.Directed(), graph.Weighted()),
			vertices: []string{"1", "2", "3"},
			edges: []graph.Edge[string]{
				{
					Source: "1",
					Target: "2",
					Properties: graph.EdgeProperties{
						Weight: 10,
						Attributes: map[string]string{
							"color": "red",
						},
					},
				},
				{Source: "1", Target: "3"},
			},
			expected: description{
				GraphType:    "digraph",
				EdgeOperator: "->",
				Statements: []statement{
					{
						Source: "1",
						Target: "2",
						Weight: 10,
						Attributes: map[string]string{
							"color": "red",
						},
					},
					{Source: "1", Target: "3"},
					{Source: "2"},
					{Source: "3"},
				},
			},
		},
	}

	for name, test := range tests {
		for _, vertex := range test.vertices {
			_ = test.graph.AddVertex(vertex)
		}

		for _, edge := range test.edges {
			var err error
			if len(edge.Properties.Attributes) == 0 {
				err = test.graph.AddEdge(edge.Source, edge.Target, graph.EdgeWeight(edge.Properties.Weight))
			}
			// If there are edge attributes, iterate over them and call EdgeAttribute for each
			// entry. An edge should only have one attribute so that AddEdge is invoked once.
			for key, value := range edge.Properties.Attributes {
				err = test.graph.AddEdge(edge.Source, edge.Target, graph.EdgeWeight(edge.Properties.Weight), graph.EdgeAttribute(key, value))
			}
			if err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		desc, _ := generateDOT(test.graph)

		if desc.GraphType != test.expected.GraphType {
			t.Errorf("%s: graph type expectancy doesn't match: expected %v, got %v", name, test.expected.GraphType, desc.GraphType)
		}

		if desc.EdgeOperator != test.expected.EdgeOperator {
			t.Errorf("%s: edge operator expectancy doesn't match: expected %v, got %v", name, test.expected.EdgeOperator, desc.EdgeOperator)
		}

		if !slicesAreEqual(desc.Statements, test.expected.Statements, statementsAreEqual) {
			t.Errorf("%s: statements expectancy doesn't match: expected %v, got %v", name, test.expected.Statements, desc.Statements)
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
				"1" -> "2" [ weight=0 ];
				"1" -> "3" [ weight=0 ];
				"2" ;
				"3" ;
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
				"1" -> "2" [ color="red", weight=0 ];
				"1" -> "3" [ color="blue", weight=0 ];
				"2" ;
				"3" ;
			}`,
		},
		"vertices containing special characters": {
			description: description{
				GraphType:    "digraph",
				EdgeOperator: "->",
				Statements: []statement{
					{Source: "/home", Target: "projects/graph"},
					{Source: "/home", Target: ".config"},
					{Source: ".config", Target: "my file.txt"},
				},
			},
			expected: `strict digraph {
				"/home" -> "projects/graph" [ weight=0 ];
				"/home" -> ".config" [ weight=0 ];
				".config" -> "my file.txt" [ weight=0 ];
			}`,
		},
	}

	for name, test := range tests {
		buf := new(bytes.Buffer)
		_ = renderDOT(buf, test.description)

		output := normalizeOutput(buf.String())
		expected := normalizeOutput(test.expected)

		if output != expected {
			t.Errorf("%s: DOT output expectancy doesn't match: expected %v, got %v", name, expected, output)
		}
	}
}

func slicesAreEqual[T any](a, b []T, equals func(a, b T) bool) bool {
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

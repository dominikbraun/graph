package draw

import (
	"bytes"
	"strings"
	"testing"

	"github.com/dominikbraun/graph"
)

func TestGenerateDOT(t *testing.T) {
	tests := map[string]struct {
		graph            graph.Graph[string, string]
		attributes       map[string]string
		vertices         []string
		vertexProperties map[string]graph.VertexProperties
		edges            []graph.Edge[string]
		expected         description
	}{
		"3-vertex directed graph": {
			graph:      graph.New(graph.StringHash, graph.Directed()),
			attributes: map[string]string{},
			vertices:   []string{"1", "2", "3"},
			edges: []graph.Edge[string]{
				{Source: "1", Target: "2"},
				{Source: "1", Target: "3"},
			},
			expected: description{
				GraphType:    "digraph",
				Attributes:   map[string]string{},
				EdgeOperator: "->",
				Statements: []statement{
					{Source: "1", Target: "2"},
					{Source: "1", Target: "3"},
					{Source: "1"},
					{Source: "2"},
					{Source: "3"},
				},
			},
		},
		"3-vertex directed, weighted graph with weights and attributes": {
			graph:      graph.New(graph.StringHash, graph.Directed(), graph.Weighted()),
			attributes: map[string]string{},
			vertices:   []string{"1", "2", "3"},
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
				Attributes:   map[string]string{},
				EdgeOperator: "->",
				Statements: []statement{
					{
						Source:     "1",
						Target:     "2",
						EdgeWeight: 10,
						EdgeAttributes: map[string]string{
							"color": "red",
						},
					},
					{Source: "1", Target: "3"},
					{Source: "1"},
					{Source: "2"},
					{Source: "3"},
				},
			},
		},
		"vertices with attributes": {
			graph:      graph.New(graph.StringHash, graph.Directed(), graph.Weighted()),
			attributes: map[string]string{},
			vertices:   []string{"1", "2"},
			vertexProperties: map[string]graph.VertexProperties{
				"1": {
					Attributes: map[string]string{
						"color": "red",
					},
					Weight: 10,
				},
				"2": {
					Attributes: map[string]string{
						"color": "blue",
					},
					Weight: 20,
				},
			},
			edges: []graph.Edge[string]{
				{Source: "1", Target: "2"},
			},
			expected: description{
				GraphType:    "digraph",
				Attributes:   map[string]string{},
				EdgeOperator: "->",
				Statements: []statement{
					{
						Source: "1",
						SourceAttributes: map[string]string{
							"color": "red",
						},
						SourceWeight: 10,
					},
					{
						Source: "2",
						SourceAttributes: map[string]string{
							"color": "blue",
						},
						SourceWeight: 20,
					},
					{
						Source: "1",
						Target: "2",
					},
				},
			},
		},
		"directed graph with attributes": {
			graph: graph.New(graph.StringHash, graph.Directed()),
			attributes: map[string]string{
				"label":     "my-graph",
				"normalize": "true",
				"compound":  "false",
			},
			vertices: []string{"1", "2"},
			edges: []graph.Edge[string]{
				{Source: "1", Target: "2"},
			},
			expected: description{
				GraphType: "digraph",
				Attributes: map[string]string{
					"label":     "my-graph",
					"normalize": "true",
					"compound":  "false",
				},
				EdgeOperator: "->",
				Statements: []statement{
					{
						Source: "1",
					},
					{
						Source: "2",
					},
					{
						Source: "1",
						Target: "2",
					},
				},
			},
		},
	}

	for name, test := range tests {
		for _, vertex := range test.vertices {
			if test.vertexProperties == nil {
				_ = test.graph.AddVertex(vertex)
				continue
			}
			// If there are vertex attributes, iterate over them and call
			// VertexAttribute for each entry. A vertex should only have one
			// attribute so that AddVertex is invoked once.
			for key, value := range test.vertexProperties[vertex].Attributes {
				weight := test.vertexProperties[vertex].Weight
				// ToDo: Clarify how multiple functional options and attributes can be tested.
				_ = test.graph.AddVertex(vertex, graph.VertexWeight(weight), graph.VertexAttribute(key, value))
			}
		}

		for _, edge := range test.edges {
			var err error
			if len(edge.Properties.Attributes) == 0 {
				err = test.graph.AddEdge(edge.Source, edge.Target, graph.EdgeWeight(edge.Properties.Weight))
			}
			// If there are edge attributes, iterate over them and call
			// EdgeAttribute for each entry. An edge should only have one
			// attribute so that AddEdge is invoked once.
			for key, value := range edge.Properties.Attributes {
				err = test.graph.AddEdge(edge.Source, edge.Target, graph.EdgeWeight(edge.Properties.Weight), graph.EdgeAttribute(key, value))
			}
			if err != nil {
				t.Fatalf("%s: failed to add edge: %s", name, err.Error())
			}
		}

		desc, _ := generateDOT(test.graph)

		// Add the graph attributes manually instead of using the functional
		// option. This is the reason why I dislike them more and more.
		desc.Attributes = test.attributes

		if desc.GraphType != test.expected.GraphType {
			t.Errorf("%s: graph type expectancy doesn't match: expected %v, got %v", name, test.expected.GraphType, desc.GraphType)
		}

		if desc.EdgeOperator != test.expected.EdgeOperator {
			t.Errorf("%s: edge operator expectancy doesn't match: expected %v, got %v", name, test.expected.EdgeOperator, desc.EdgeOperator)
		}

		if !slicesAreEqual(desc.Statements, test.expected.Statements, statementsAreEqual) {
			t.Errorf("%s: statements expectancy doesn't match: expected %v, got %v", name, test.expected.Statements, desc.Statements)
		}

		stringsAreEqual := func(a, b string) bool {
			return a == b
		}

		if !mapsAreEqual(desc.Attributes, test.expected.Attributes, stringsAreEqual) {
			t.Errorf("%s: attributes expectancy doesn't match: expected %v, got %v", name, test.expected.Attributes, desc.Attributes)
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
				Attributes:   map[string]string{},
				EdgeOperator: "->",
				Statements: []statement{
					{Source: 1, Target: 2},
					{Source: 1, Target: 3},
					{Source: 1},
					{Source: 2},
					{Source: 3},
				},
			},
			expected: `strict digraph {
				"1" -> "2" [ weight=0 ];
				"1" -> "3" [ weight=0 ];
				"1" [ weight=0 ];
				"2" [ weight=0 ];
				"3" [ weight=0 ];
			}`,
		},
		"custom edge attributes": {
			description: description{
				GraphType:    "digraph",
				Attributes:   map[string]string{},
				EdgeOperator: "->",
				Statements: []statement{
					{
						Source: 1,
						Target: 2,
						EdgeAttributes: map[string]string{
							"color": "red",
						},
					},
					{
						Source: 1,
						Target: 3,
						EdgeAttributes: map[string]string{
							"color": "blue",
						},
					},
					{Source: 1},
					{Source: 2},
					{Source: 3},
				},
			},
			expected: `strict digraph {
				"1" -> "2" [ color="red", weight=0 ];
				"1" -> "3" [ color="blue", weight=0 ];
				"1" [ weight=0 ];
				"2" [ weight=0 ];
				"3" [ weight=0 ];
			}`,
		},
		"vertices containing special characters": {
			description: description{
				GraphType:    "digraph",
				Attributes:   map[string]string{},
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
		"vertices with attributes": {
			description: description{
				GraphType:    "digraph",
				Attributes:   map[string]string{},
				EdgeOperator: "->",
				Statements: []statement{
					{
						Source: "1",
						SourceAttributes: map[string]string{
							"color": "red",
						},
						SourceWeight: 10,
					},
					{
						Source: "2",
						SourceAttributes: map[string]string{
							"color": "blue",
						},
						SourceWeight: 20,
					},
					{
						Source: "1",
						Target: "2",
					},
				},
			},
			expected: `strict digraph {
				"1" [ color="red", weight=10 ];
				"2" [ color="blue", weight=20 ];
				"1" -> "2" [ weight=0 ];
			}`,
		},
		"3-vertex directed graph with attributes": {
			description: description{
				GraphType: "digraph",
				Attributes: map[string]string{
					"label": "my-graph",
				},
				EdgeOperator: "->",
				Statements: []statement{
					{Source: 1, Target: 2},
					{Source: 1, Target: 3},
					{Source: 1},
					{Source: 2},
					{Source: 3},
				},
			},
			expected: `strict digraph {
				label="my-graph";
				"1" -> "2" [ weight=0 ];
				"1" -> "3" [ weight=0 ];
				"1" [ weight=0 ];
				"2" [ weight=0 ];
				"3" [ weight=0 ];
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

func mapsAreEqual[K comparable, V any](a, b map[K]V, equals func(a, b V) bool) bool {
	if len(a) != len(b) {
		return false
	}

	for aKey, aValue := range a {
		bValue, ok := b[aKey]
		if !ok {
			return false
		}
		if !equals(aValue, bValue) {
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
	if len(a.EdgeAttributes) != len(b.EdgeAttributes) {
		return false
	}

	for aKey, aValue := range a.EdgeAttributes {
		bValue, ok := b.EdgeAttributes[aKey]
		if !ok {
			return false
		}
		if aValue != bValue {
			return false
		}
	}

	if len(a.SourceAttributes) != len(b.SourceAttributes) {
		return false
	}

	for aKey, aValue := range a.SourceAttributes {
		bValue, ok := b.SourceAttributes[aKey]
		if !ok {
			return false
		}
		if aValue != bValue {
			return false
		}
	}

	return a.Source == b.Source &&
		a.Target == b.Target &&
		a.EdgeWeight == b.EdgeWeight &&
		a.SourceWeight == b.SourceWeight
}

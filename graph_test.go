package graph

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	directedType := reflect.TypeOf(&directed[int, int]{})
	undirectedType := reflect.TypeOf(&undirected[int, int]{})

	tests := map[string]struct {
		expectedType reflect.Type
		options      []func(*Traits)
	}{
		"no options": {
			options:      []func(*Traits){},
			expectedType: undirectedType,
		},
		"directed option": {
			options:      []func(*Traits){Directed()},
			expectedType: directedType,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			graph := New(IntHash, test.options...)
			actualType := reflect.TypeOf(graph)

			if actualType != test.expectedType {
				t.Errorf("graph type expectancy doesn't match: expected %v, got %v", test.expectedType, actualType)
			}
		})
	}
}

func TestNewLike(t *testing.T) {
	tests := map[string]struct {
		g        Graph[int, int]
		vertices []int
	}{
		"new directed graph of integers": {
			g:        New(IntHash, Directed()),
			vertices: []int{1, 2, 3},
		},
		"new undirected weighted graph of integers": {
			g:        New(IntHash, Weighted()),
			vertices: []int{1, 2, 3},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			for _, vertex := range test.vertices {
				_ = test.g.AddVertex(vertex)
			}

			h := NewLike(test.g)

			if len(test.vertices) > 0 {
				if _, err := h.Vertex(test.vertices[0]); err == nil {
					t.Errorf("expected vertex %v not to exist in h", test.vertices[0])
				}
			}

			if test.g.Traits().IsDirected {
				actual, ok := h.(*directed[int, int])
				if !ok {
					t.Fatalf("type assertion to *directed failed")
				}

				expected := test.g.(*directed[int, int])

				if actual.hash(42) != expected.hash(42) {
					t.Errorf("expected hash %v, got %v", expected.hash, actual.hash)
				}
			} else {
				actual, ok := h.(*undirected[int, int])
				if !ok {
					t.Fatalf("type assertion to *undirected failed")
				}

				expected := test.g.(*undirected[int, int])

				if actual.hash(42) != expected.hash(42) {
					t.Errorf("expected hash %v, got %v", expected.hash, actual.hash)
				}
			}

			if !traitsAreEqual(h.Traits(), test.g.Traits()) {
				t.Errorf("expected traits %+v, got %+v", test.g.Traits(), h.Traits())
			}
		})
	}
}

func TestStringHash(t *testing.T) {
	tests := map[string]struct {
		value        string
		expectedHash string
	}{
		"string value": {
			value:        "London",
			expectedHash: "London",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			hash := StringHash(test.value)

			if hash != test.expectedHash {
				t.Errorf("hash expectancy doesn't match: expected %v, got %v", test.expectedHash, hash)
			}
		})
	}
}

func TestIntHash(t *testing.T) {
	tests := map[string]struct {
		value        int
		expectedHash int
	}{
		"int value": {
			value:        3,
			expectedHash: 3,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			hash := IntHash(test.value)

			if hash != test.expectedHash {
				t.Errorf("hash expectancy doesn't match: expected %v, got %v", test.expectedHash, hash)
			}
		})
	}
}

func TestEdgeWeight(t *testing.T) {
	tests := map[string]struct {
		expected EdgeProperties
		weight   int
	}{
		"weight 4": {
			weight: 4,
			expected: EdgeProperties{
				Weight: 4,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			properties := EdgeProperties{}

			EdgeWeight(test.weight)(&properties)

			if properties.Weight != test.expected.Weight {
				t.Errorf("weight expectation doesn't match: expected %v, got %v", test.expected.Weight, properties.Weight)
			}
		})
	}
}

func TestEdgeAttribute(t *testing.T) {
	tests := map[string]struct {
		key      string
		value    string
		expected EdgeProperties
	}{
		"attribute label=my-label": {
			key:   "label",
			value: "my-label",
			expected: EdgeProperties{
				Attributes: map[string]string{
					"label": "my-label",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			properties := EdgeProperties{
				Attributes: make(map[string]string),
			}

			EdgeAttribute(test.key, test.value)(&properties)

			value, ok := properties.Attributes[test.key]
			if !ok {
				t.Errorf("attribute expectaton doesn't match: key %v doesn't exist", test.key)
			}

			expectedValue := test.expected.Attributes[test.key]

			if value != expectedValue {
				t.Errorf("value expectation doesn't match: expected %v, got %v", expectedValue, value)
			}
		})

	}
}

func TestEdgeAttributes(t *testing.T) {
	tests := map[string]struct {
		attributes map[string]string
		expected   map[string]string
	}{
		"attribute label=my-label": {
			attributes: map[string]string{
				"label": "my-label",
			},
			expected: map[string]string{
				"label": "my-label",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			properties := EdgeProperties{
				Attributes: make(map[string]string),
			}

			EdgeAttributes(test.attributes)(&properties)

			if !mapsAreEqual(test.expected, properties.Attributes) {
				t.Errorf("expected %v, got %v", test.expected, properties.Attributes)
			}
		})
	}
}

func TestVertexAttribute(t *testing.T) {
	tests := map[string]struct {
		key      string
		value    string
		expected VertexProperties
	}{
		"attribute label=my-label": {
			key:   "label",
			value: "my-label",
			expected: VertexProperties{
				Attributes: map[string]string{
					"label": "my-label",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			properties := VertexProperties{
				Attributes: make(map[string]string),
			}

			VertexAttribute(test.key, test.value)(&properties)

			value, ok := properties.Attributes[test.key]
			if !ok {
				t.Errorf("attribute expectaton doesn't match: key %v doesn't exist", test.key)
			}

			expectedValue := test.expected.Attributes[test.key]

			if value != expectedValue {
				t.Errorf("value expectation doesn't match: expected %v, got %v", expectedValue, value)
			}
		})

	}
}

func TestVertexAttributes(t *testing.T) {
	tests := map[string]struct {
		attributes map[string]string
		expected   map[string]string
	}{
		"attribute label=my-label": {
			attributes: map[string]string{
				"label": "my-label",
			},
			expected: map[string]string{
				"label": "my-label",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			properties := VertexProperties{
				Attributes: make(map[string]string),
			}

			VertexAttributes(test.attributes)(&properties)

			if !mapsAreEqual(test.expected, properties.Attributes) {
				t.Errorf("expected %v, got %v", test.expected, properties.Attributes)
			}
		})
	}
}

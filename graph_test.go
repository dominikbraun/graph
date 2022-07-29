package graph

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	directedType := reflect.TypeOf(&directed[int, int]{})
	undirectedType := reflect.TypeOf(&undirected[int, int]{})

	tests := map[string]struct {
		options      []func(*Traits)
		expectedType reflect.Type
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
		graph := New(IntHash, test.options...)
		actualType := reflect.TypeOf(graph)

		if actualType != test.expectedType {
			t.Errorf("%s: graph type expectancy doesn't match: expected %v, got %v", name, test.expectedType, actualType)
		}
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
		hash := StringHash(test.value)

		if hash != test.expectedHash {
			t.Errorf("%s: hash expectancy doesn't match: expected %v, got %v", name, test.expectedHash, hash)
		}
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
		hash := IntHash(test.value)

		if hash != test.expectedHash {
			t.Errorf("%s: hash expectancy doesn't match: expected %v, got %v", name, test.expectedHash, hash)
		}
	}
}

func TestEdgeWeight(t *testing.T) {
	tests := map[string]struct {
		weight   int
		expected edgeProperties
	}{
		"weight 4": {
			weight: 4,
			expected: edgeProperties{
				Weight: 4,
			},
		},
	}

	for name, test := range tests {
		properties := edgeProperties{}

		EdgeWeight(test.weight)(&properties)

		if properties.Weight != test.expected.Weight {
			t.Errorf("%s: weight expectation doesn't match: expected %v, got %v", name, test.expected.Weight, properties.Weight)
		}
	}
}

func TestEdgeAttribute(t *testing.T) {
	tests := map[string]struct {
		key      string
		value    string
		expected edgeProperties
	}{
		"attribute label=mylabel": {
			key:   "label",
			value: "mylabel",
			expected: edgeProperties{
				Attributes: map[string]string{
					"label": "mylabel",
				},
			},
		},
	}

	for name, test := range tests {
		properties := edgeProperties{
			Attributes: make(map[string]string),
		}

		EdgeAttribute(test.key, test.value)(&properties)

		value, ok := properties.Attributes[test.key]
		if !ok {
			t.Errorf("%s: attribute expectaton doesn't match: key %v doesn't exist", name, test.key)
		}

		expectedValue := test.expected.Attributes[test.key]

		if value != expectedValue {
			t.Errorf("%s: value expectation doesn't match: expected %v, got %v", name, expectedValue, value)
		}
	}
}

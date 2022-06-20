package graph

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	directedType := reflect.TypeOf(&directed[int, int]{})
	undirectedType := reflect.TypeOf(&undirected[int, int]{})

	tests := map[string]struct {
		options      []func(*properties)
		expectedType reflect.Type
	}{
		"no options": {
			options:      []func(*properties){},
			expectedType: undirectedType,
		},
		"directed option": {
			options:      []func(*properties){Directed()},
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

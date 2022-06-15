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

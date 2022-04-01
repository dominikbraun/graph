package graph

import (
	"testing"
)

func TestSet_Values(t *testing.T) {
	tests := map[string]struct {
		set      Set[int]
		expected []int
	}{
		"set with 3 elements": {
			set: Set[int]{
				1: 100,
				2: 150,
				3: 200,
			},
			expected: []int{100, 150, 200},
		},
		"empty set": {
			set:      Set[int]{},
			expected: []int{},
		},
	}

	equals := func(a, b int) bool {
		return a == b
	}

	for name, test := range tests {
		actual := test.set.Values()

		if !slicesAreEqual(actual, test.expected, equals) {
			t.Fatalf("%s: expected %v, got %v", name, test.expected, actual)
		}
	}
}

func TestGraph_Vertices(t *testing.T) {
	tests := map[string]struct {
		graph    Graph[int]
		expected []int
	}{
		"graph with 3 vertices": {
			graph: Graph[int]{
				vertices: Set[int]{
					1: 100,
					2: 150,
					3: 200,
				},
			},
			expected: []int{100, 150, 200},
		},
	}

	equals := func(a, b int) bool {
		return a == b
	}

	for name, test := range tests {
		actual := test.graph.Vertices()

		if !slicesAreEqual(actual, test.expected, equals) {
			t.Fatalf("%s: expected %v, got %v", name, test.expected, actual)
		}
	}
}

func TestGraph_Edges(t *testing.T) {
	tests := map[string]struct {
		graph    Graph[int]
		expected []Pair[int]
	}{
		"graph with 3 edges": {
			graph: Graph[int]{
				edges: Set[Pair[int]]{
					Pair[int]{A: 100, B: 150}: Pair[int]{A: 100, B: 150},
					Pair[int]{A: 150, B: 200}: Pair[int]{A: 150, B: 200},
					Pair[int]{A: 200, B: 100}: Pair[int]{A: 200, B: 100},
				},
			},
			expected: []Pair[int]{
				{A: 100, B: 150},
				{A: 150, B: 200},
				{A: 200, B: 100},
			},
		},
	}

	equals := func(a, b Pair[int]) bool {
		return a.Equals(b)
	}

	for name, test := range tests {
		actual := test.graph.Edges()

		if !slicesAreEqual(actual, test.expected, equals) {
			t.Fatalf("%s: expected %v, got %v", name, test.expected, actual)
		}
	}
}

func TestPair_Equals(t *testing.T) {
	tests := map[string]struct {
		a, b     Pair[int]
		expected bool
	}{
		"equal unordered pairs": {
			a:        Pair[int]{A: 100, B: 150},
			b:        Pair[int]{A: 150, B: 100},
			expected: true,
		},
		"equal ordered pairs": {
			a:        Pair[int]{A: 100, B: 150},
			b:        Pair[int]{A: 100, B: 150},
			expected: true,
		},
		"unequal pairs": {
			a:        Pair[int]{A: 100, B: 150},
			b:        Pair[int]{A: 150, B: 200},
			expected: false,
		},
	}

	for name, test := range tests {
		actual := test.a.Equals(test.b)

		if actual != test.expected {
			t.Fatalf("%s: expected %v, got %v", name, test.expected, actual)
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

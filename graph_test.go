package graph

import "testing"

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

	for name, test := range tests {
		actual := test.set.Values()

		if !slicesAreEqual(actual, test.expected) {
			t.Fatalf("%s: expected %v, got %v", name, test.expected, actual)
		}
	}
}

func slicesAreEqual[T comparable](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	for _, aValue := range a {
		found := false
		for _, bValue := range b {
			if aValue == bValue {
				found = true
			}
		}
		if !found {
			return false
		}
	}

	return true
}

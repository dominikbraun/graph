package graph

import "testing"

func TestDirected(t *testing.T) {
	tests := map[string]struct {
		expected *traits
	}{
		"directed graph": {
			expected: &traits{
				isDirected: true,
			},
		},
	}

	for name, test := range tests {
		p := &traits{}

		Directed()(p)

		if !traitsAreEqual(test.expected, p) {
			t.Errorf("%s: trait expectation doesn't match: expected %v, got %v", name, test.expected, p)
		}
	}
}

func TestAcyclic(t *testing.T) {
	tests := map[string]struct {
		expected *traits
	}{
		"acyclic graph": {
			expected: &traits{
				isAcyclic: true,
			},
		},
	}

	for name, test := range tests {
		p := &traits{}

		Acyclic()(p)

		if !traitsAreEqual(test.expected, p) {
			t.Errorf("%s: trait expectation doesn't match: expected %v, got %v", name, test.expected, p)
		}
	}
}

func TestWeighted(t *testing.T) {
	tests := map[string]struct {
		expected *traits
	}{
		"weighted graph": {
			expected: &traits{
				isWeighted: true,
			},
		},
	}

	for name, test := range tests {
		p := &traits{}

		Weighted()(p)

		if !traitsAreEqual(test.expected, p) {
			t.Errorf("%s: trait expectation doesn't match: expected %v, got %v", name, test.expected, p)
		}
	}
}

func TestRooted(t *testing.T) {
	tests := map[string]struct {
		expected *traits
	}{
		"rooted graph": {
			expected: &traits{
				isRooted: true,
			},
		},
	}

	for name, test := range tests {
		p := &traits{}

		Rooted()(p)

		if !traitsAreEqual(test.expected, p) {
			t.Errorf("%s: trait expectation doesn't match: expected %v, got %v", name, test.expected, p)
		}
	}
}

func TestTree(t *testing.T) {
	tests := map[string]struct {
		expected *traits
	}{
		"tree graph": {
			expected: &traits{
				isAcyclic: true,
				isRooted:  true,
			},
		},
	}

	for name, test := range tests {
		p := &traits{}

		Tree()(p)

		if !traitsAreEqual(test.expected, p) {
			t.Errorf("%s: trait expectation doesn't match: expected %v, got %v", name, test.expected, p)
		}
	}
}

func traitsAreEqual(a, b *traits) bool {
	return a.isAcyclic == b.isAcyclic &&
		a.isDirected == b.isDirected &&
		a.isRooted == b.isRooted &&
		a.isWeighted == b.isWeighted
}

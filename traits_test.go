package graph

import "testing"

func TestDirected(t *testing.T) {
	tests := map[string]struct {
		expected *Traits
	}{
		"directed graph": {
			expected: &Traits{
				IsDirected: true,
			},
		},
	}

	for name, test := range tests {
		p := &Traits{}

		Directed()(p)

		if !traitsAreEqual(test.expected, p) {
			t.Errorf("%s: trait expectation doesn't match: expected %v, got %v", name, test.expected, p)
		}
	}
}

func TestAcyclic(t *testing.T) {
	tests := map[string]struct {
		expected *Traits
	}{
		"acyclic graph": {
			expected: &Traits{
				IsAcyclic: true,
			},
		},
	}

	for name, test := range tests {
		p := &Traits{}

		Acyclic()(p)

		if !traitsAreEqual(test.expected, p) {
			t.Errorf("%s: trait expectation doesn't match: expected %v, got %v", name, test.expected, p)
		}
	}
}

func TestWeighted(t *testing.T) {
	tests := map[string]struct {
		expected *Traits
	}{
		"weighted graph": {
			expected: &Traits{
				IsWeighted: true,
			},
		},
	}

	for name, test := range tests {
		p := &Traits{}

		Weighted()(p)

		if !traitsAreEqual(test.expected, p) {
			t.Errorf("%s: trait expectation doesn't match: expected %v, got %v", name, test.expected, p)
		}
	}
}

func TestRooted(t *testing.T) {
	tests := map[string]struct {
		expected *Traits
	}{
		"rooted graph": {
			expected: &Traits{
				IsRooted: true,
			},
		},
	}

	for name, test := range tests {
		p := &Traits{}

		Rooted()(p)

		if !traitsAreEqual(test.expected, p) {
			t.Errorf("%s: trait expectation doesn't match: expected %v, got %v", name, test.expected, p)
		}
	}
}

func TestTree(t *testing.T) {
	tests := map[string]struct {
		expected *Traits
	}{
		"tree graph": {
			expected: &Traits{
				IsAcyclic: true,
				IsRooted:  true,
			},
		},
	}

	for name, test := range tests {
		p := &Traits{}

		Tree()(p)

		if !traitsAreEqual(test.expected, p) {
			t.Errorf("%s: trait expectation doesn't match: expected %v, got %v", name, test.expected, p)
		}
	}
}

func TestPreventCycles(t *testing.T) {
	tests := map[string]struct {
		expected *Traits
	}{
		"prevent cycles": {
			expected: &Traits{
				IsAcyclic:     true,
				PreventCycles: true,
			},
		},
	}

	for name, test := range tests {
		p := &Traits{}

		PreventCycles()(p)

		if !traitsAreEqual(test.expected, p) {
			t.Errorf("%s: trait expectation doesn't match: expected %v, got %v", name, test.expected, p)
		}
	}
}

func traitsAreEqual(a, b *Traits) bool {
	return a.IsAcyclic == b.IsAcyclic &&
		a.IsDirected == b.IsDirected &&
		a.IsRooted == b.IsRooted &&
		a.IsWeighted == b.IsWeighted &&
		a.PreventCycles == b.PreventCycles
}

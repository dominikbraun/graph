package graph

import "testing"

func TestDirected(t *testing.T) {
	tests := map[string]struct {
		expected *properties
	}{
		"directed graph": {
			expected: &properties{
				isDirected: true,
			},
		},
	}

	for name, test := range tests {
		p := &properties{}

		Directed()(p)

		if !propertiesAreEqual(test.expected, p) {
			t.Errorf("%s: property expectation doesn't match: expected %v, got %v", name, test.expected, p)
		}
	}
}

func TestAcyclic(t *testing.T) {
	tests := map[string]struct {
		expected *properties
	}{
		"acyclic graph": {
			expected: &properties{
				isAcyclic: true,
			},
		},
	}

	for name, test := range tests {
		p := &properties{}

		Acyclic()(p)

		if !propertiesAreEqual(test.expected, p) {
			t.Errorf("%s: property expectation doesn't match: expected %v, got %v", name, test.expected, p)
		}
	}
}

func TestWeighted(t *testing.T) {
	tests := map[string]struct {
		expected *properties
	}{
		"weighted graph": {
			expected: &properties{
				isWeighted: true,
			},
		},
	}

	for name, test := range tests {
		p := &properties{}

		Weighted()(p)

		if !propertiesAreEqual(test.expected, p) {
			t.Errorf("%s: property expectation doesn't match: expected %v, got %v", name, test.expected, p)
		}
	}
}

func TestRooted(t *testing.T) {
	tests := map[string]struct {
		expected *properties
	}{
		"rooted graph": {
			expected: &properties{
				isRooted: true,
			},
		},
	}

	for name, test := range tests {
		p := &properties{}

		Rooted()(p)

		if !propertiesAreEqual(test.expected, p) {
			t.Errorf("%s: property expectation doesn't match: expected %v, got %v", name, test.expected, p)
		}
	}
}

func TestTree(t *testing.T) {
	tests := map[string]struct {
		expected *properties
	}{
		"tree graph": {
			expected: &properties{
				isAcyclic: true,
				isRooted:  true,
			},
		},
	}

	for name, test := range tests {
		p := &properties{}

		Tree()(p)

		if !propertiesAreEqual(test.expected, p) {
			t.Errorf("%s: property expectation doesn't match: expected %v, got %v", name, test.expected, p)
		}
	}
}

func propertiesAreEqual(a, b *properties) bool {
	return a.isAcyclic == b.isAcyclic &&
		a.isDirected == b.isDirected &&
		a.isRooted == b.isRooted &&
		a.isWeighted == b.isWeighted
}

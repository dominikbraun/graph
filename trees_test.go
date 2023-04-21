package graph

import "testing"

func TestUnionFind_add(t *testing.T) {
	tests := map[string]struct {
		vertex         int
		expectedParent int
	}{
		"add vertex 1": {
			vertex:         1,
			expectedParent: 1,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			u := newUnionFind[int]()
			u.add(test.vertex)

			if u.parents[test.vertex] != test.expectedParent {
				t.Errorf("expected parent %v, got %v", test.expectedParent, u.parents[test.vertex])
			}
		})
	}
}

func TestUnionFind_union(t *testing.T) {
	tests := map[string]struct {
		vertices        []int
		args            [2]int
		expectedParents map[int]int
	}{
		"merge 1 and 2": {
			vertices: []int{1, 2, 3, 4},
			args:     [2]int{1, 2},
			expectedParents: map[int]int{
				1: 1,
				2: 1,
				3: 3,
				4: 4,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			u := newUnionFind(test.vertices...)
			u.union(test.args[0], test.args[1])

			for expectedKey, expectedValue := range test.expectedParents {
				actualValue, ok := u.parents[expectedKey]
				if !ok {
					t.Fatalf("expected key %v not found", expectedKey)
				}
				if actualValue != expectedValue {
					t.Fatalf("expected value %v for %v, got %v", expectedValue, expectedKey, actualValue)
				}
			}
		})
	}
}

func TestUnionFind_find(t *testing.T) {
	tests := map[string]struct {
		parents  map[int]int
		arg      int
		expected int
	}{
		"find 1 in sets 1, 2, 3": {
			parents: map[int]int{
				1: 1,
				2: 2,
				3: 3,
			},
			arg:      1,
			expected: 1,
		},
		"find 3 in sets 1-2, 3-4": {
			parents: map[int]int{
				1: 1,
				2: 1,
				3: 3,
				4: 3,
			},
			arg:      3,
			expected: 3,
		},
		"find 4 in sets 1-2, 3-4": {
			parents: map[int]int{
				1: 1,
				2: 1,
				3: 3,
				4: 3,
			},
			arg:      4,
			expected: 3,
		},
		"find 3 in set 1-2-3": {
			parents: map[int]int{
				1: 1,
				2: 1,
				3: 2,
			},
			arg:      3,
			expected: 1,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			u := unionFind[int]{
				parents: test.parents,
			}

			actual := u.find(test.arg)

			if actual != test.expected {
				t.Errorf("expected %v, got %v", test.expected, actual)
			}
		})
	}
}

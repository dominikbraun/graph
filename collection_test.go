package graph

import (
	"testing"
)

func TestPriorityQueue_Push(t *testing.T) {
	tests := map[string]struct {
		items              []int
		priorities         []int
		expectedItems      []int
		expectedPriorities []int
	}{
		"queue with 5 elements": {
			items:              []int{10, 20, 30, 40, 50},
			priorities:         []int{6, 8, 2, 7, 5},
			expectedItems:      []int{20, 40, 10, 50, 30},
			expectedPriorities: []int{8, 7, 6, 5, 2},
		},
	}

	for name, test := range tests {
		queue := newPriorityQueue[int]()

		for i, item := range test.items {
			queue.Push(item, test.priorities[i])
		}

		if len(queue.items) != len(test.expectedItems) {
			t.Fatalf("%s: item length expectancy doesn't match: expected %v, got %v", name, len(test.expectedItems), len(queue.items))
		}

		if len(queue.priorities) != len(test.expectedPriorities) {
			t.Fatalf("%s: priorities length expectancy doesn't match: expected %v, got %v", name, len(test.expectedPriorities), len(queue.priorities))
		}

		for i, expectedItem := range test.expectedItems {
			if queue.items[i] != expectedItem {
				t.Errorf("%s: item expectancy doesn't match: expected %v at index %d, got %v", name, expectedItem, i, queue.items[i])
			}
		}

		for i, expectedPriority := range test.expectedPriorities {
			if queue.priorities[i] != expectedPriority {
				t.Errorf("%s: priority expectancy doesn't match: expected %v at index %d, got %v", name, expectedPriority, i, queue.priorities[i])
			}
		}
	}
}

func TestPriorityQueue_Pop(t *testing.T) {
	tests := map[string]struct {
		items        []int
		priorities   []int
		expectedItem int
		shouldFail   bool
	}{}

	for name, test := range tests {
		queue := newPriorityQueue[int]()

		for i, item := range test.items {
			queue.Push(item, test.priorities[i])
		}

		item, err := queue.Pop()

		if test.shouldFail != (err != nil) {
			t.Fatalf("%s: error expectancy doesn't match: expected %v, got %v (error: %v)", name, test.shouldFail, (err != nil), err)
		}

		if item != test.expectedItem {
			t.Errorf("%s: item expectancy doesn't match: expected %v, got %v", name, test.expectedItem, item)
		}
	}
}

func TestPriorityQueue_Len(t *testing.T) {
	tests := map[string]struct {
		items       []int
		priorities  []int
		expectedLen int
	}{}

	for name, test := range tests {
		queue := newPriorityQueue[int]()

		for i, item := range test.items {
			queue.Push(item, test.priorities[i])
		}

		len := queue.Len()

		if len != test.expectedLen {
			t.Errorf("%s: length expectancy doesn't match: expected %v, got %v", name, test.expectedLen, len)
		}
	}
}

func TestPriorityQueue_insertItemAt(t *testing.T) {
	tests := map[string]struct {
		items         []int
		priorities    []int
		insertItem    int
		insertIndex   int
		expectedItems []int
	}{
		"insert in the middle of the queue": {
			items:         []int{10, 20, 30, 40},
			priorities:    []int{1, 2, 3, 4},
			insertItem:    25,
			insertIndex:   2,
			expectedItems: []int{10, 20, 25, 30, 40},
		},
		"insert at the start of the queue": {
			items:         []int{10, 20, 30, 40},
			priorities:    []int{1, 2, 3, 4},
			insertItem:    5,
			insertIndex:   0,
			expectedItems: []int{5, 10, 20, 30, 40},
		},
		"insert at the end of the queue": {
			items:         []int{10, 20, 30, 40},
			priorities:    []int{1, 2, 3, 4},
			insertItem:    50,
			insertIndex:   4,
			expectedItems: []int{10, 20, 30, 40, 50},
		},
	}

	for name, test := range tests {
		queue := &priorityQueue[int]{
			items:      test.items,
			priorities: test.priorities,
		}

		queue.insertItemAt(test.insertItem, 0, test.insertIndex)

		if len(queue.items) != len(test.expectedItems) {
			t.Fatalf("%s: item length expectancy doesn't match: expected %v, got %v", name, len(test.items), len(queue.items))
		}

		for i, expectedItem := range test.expectedItems {
			if queue.items[i] != expectedItem {
				t.Errorf("%s: item expectancy doesn't match: expected %v at index %d, got %v", name, expectedItem, i, queue.items[i])
			}
		}
	}
}

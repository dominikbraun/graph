package graph

import (
	"testing"
)

func TestPriorityQueue_Push(t *testing.T) {
	tests := map[string]struct {
		items                 []int
		priorities            []float64
		expectedPriorityItems []priorityItem[int]
	}{
		"queue with 5 elements": {
			items:      []int{10, 20, 30, 40, 50},
			priorities: []float64{6, 8, 2, 7, 5},
			expectedPriorityItems: []priorityItem[int]{
				{value: 20, priority: 8},
				{value: 40, priority: 7},
				{value: 10, priority: 6},
				{value: 50, priority: 5},
				{value: 30, priority: 2},
			},
		},
	}

	for name, test := range tests {
		queue := newPriorityQueue[int]()

		for i, item := range test.items {
			queue.Push(item, test.priorities[i])
		}

		if len(queue.items) != len(test.expectedPriorityItems) {
			t.Fatalf("%s: item length expectancy doesn't match: expected %v, got %v", name, len(test.expectedPriorityItems), len(queue.items))
		}

		for i, expectedPriorityItem := range test.expectedPriorityItems {
			if queue.items[i] != expectedPriorityItem {
				t.Errorf("%s: item doesn't match: expected %v at index %d, got %v", name, expectedPriorityItem, i, queue.items[i])
			}
		}
	}
}

func TestPriorityQueue_Pop(t *testing.T) {
	tests := map[string]struct {
		items        []int
		priorities   []float64
		expectedItem int
		shouldFail   bool
	}{
		"queue with 5 item": {
			items:        []int{10, 20, 30, 40, 50},
			priorities:   []float64{6, 8, 2, 7, 5},
			expectedItem: 30,
			shouldFail:   false,
		},
		"queue with 1 item": {
			items:        []int{10},
			priorities:   []float64{6},
			expectedItem: 10,
			shouldFail:   false,
		},
		"empty queue": {
			items:      []int{},
			priorities: []float64{},
			shouldFail: true,
		},
	}

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

func TestPriorityQueue_DecreasePriority(t *testing.T) {
	tests := map[string]struct {
		items                 []priorityItem[int]
		decreaseItem          int
		decreasePriority      float64
		expectedPriorityItems []priorityItem[int]
	}{
		"decrease 30 to priority 5": {
			items: []priorityItem[int]{
				{value: 40, priority: 40},
				{value: 30, priority: 30},
				{value: 20, priority: 20},
				{value: 10, priority: 10},
			},
			decreaseItem:     30,
			decreasePriority: 5,
			expectedPriorityItems: []priorityItem[int]{
				{value: 40, priority: 40},
				{value: 20, priority: 20},
				{value: 10, priority: 10},
				{value: 30, priority: 5},
			},
		},
		"decrease a non-existent item": {
			items: []priorityItem[int]{
				{value: 40, priority: 40},
				{value: 30, priority: 30},
				{value: 20, priority: 20},
				{value: 10, priority: 10},
			},
			decreaseItem:     50,
			decreasePriority: 10,
			expectedPriorityItems: []priorityItem[int]{
				{value: 40, priority: 40},
				{value: 30, priority: 30},
				{value: 20, priority: 20},
				{value: 10, priority: 10},
			},
		},
	}

	for name, test := range tests {
		queue := &priorityQueue[int]{
			items: test.items,
		}

		queue.DecreasePriority(test.decreaseItem, test.decreasePriority)

		if len(queue.items) != len(test.expectedPriorityItems) {
			t.Fatalf("%s: item length expectancy doesn't match: expected %v, got %v", name, len(test.expectedPriorityItems), len(queue.items))
		}

		for i, expectedPriorityItem := range test.expectedPriorityItems {
			if queue.items[i] != expectedPriorityItem {
				t.Errorf("%s: item doesn't match: expected %v at index %d, got %v", name, expectedPriorityItem, i, queue.items[i])
			}
		}
	}
}

func TestPriorityQueue_Len(t *testing.T) {
	tests := map[string]struct {
		items       []int
		priorities  []float64
		expectedLen int
	}{
		"queue with 5 item": {
			items:       []int{10, 20, 30, 40, 50},
			priorities:  []float64{6, 8, 2, 7, 5},
			expectedLen: 5,
		},
		"queue with 1 item": {
			items:       []int{10},
			priorities:  []float64{6},
			expectedLen: 1,
		},
		"empty queue": {
			items:       []int{},
			priorities:  []float64{},
			expectedLen: 0,
		},
	}

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
		items                 []priorityItem[int]
		insertItem            int
		insertPriority        float64
		insertIndex           int
		expectedPriorityItems []priorityItem[int]
	}{
		"insert in the middle of the queue": {
			items: []priorityItem[int]{
				{value: 10, priority: 10},
				{value: 20, priority: 20},
				{value: 30, priority: 30},
				{value: 40, priority: 40},
			},
			insertItem:     25,
			insertPriority: 25,
			insertIndex:    2,
			expectedPriorityItems: []priorityItem[int]{
				{value: 10, priority: 10},
				{value: 20, priority: 20},
				{value: 25, priority: 25},
				{value: 30, priority: 30},
				{value: 40, priority: 40},
			},
		},
		"insert at the start of the queue": {
			items: []priorityItem[int]{
				{value: 10, priority: 10},
				{value: 20, priority: 20},
				{value: 30, priority: 30},
				{value: 40, priority: 40},
			},
			insertItem:     5,
			insertPriority: 5,
			insertIndex:    0,
			expectedPriorityItems: []priorityItem[int]{
				{value: 5, priority: 5},
				{value: 10, priority: 10},
				{value: 20, priority: 20},
				{value: 30, priority: 30},
				{value: 40, priority: 40},
			},
		},
		"insert at the end of the queue": {
			items: []priorityItem[int]{
				{value: 10, priority: 10},
				{value: 20, priority: 20},
				{value: 30, priority: 30},
				{value: 40, priority: 40},
			},
			insertItem:     50,
			insertPriority: 50,
			insertIndex:    4,
			expectedPriorityItems: []priorityItem[int]{
				{value: 10, priority: 10},
				{value: 20, priority: 20},
				{value: 30, priority: 30},
				{value: 40, priority: 40},
				{value: 50, priority: 50},
			},
		},
	}

	for name, test := range tests {
		queue := &priorityQueue[int]{
			items: test.items,
		}

		queue.insertItemAt(test.insertItem, test.insertPriority, test.insertIndex)

		if len(queue.items) != len(test.expectedPriorityItems) {
			t.Fatalf("%s: item length expectancy doesn't match: expected %v, got %v", name, len(test.expectedPriorityItems), len(queue.items))
		}

		for i, expectedPriorityItem := range test.expectedPriorityItems {
			if queue.items[i] != expectedPriorityItem {
				t.Errorf("%s: item doesn't match: expected %v at index %d, got %v", name, expectedPriorityItem, i, queue.items[i])
			}
		}
	}
}

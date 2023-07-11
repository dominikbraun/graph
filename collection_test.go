package graph

import (
	"reflect"
	"testing"
)

func TestPriorityQueue_Push(t *testing.T) {
	tests := map[string]struct {
		items                 []int
		priorities            []float64
		expectedPriorityItems []*priorityItem[int]
	}{
		"queue with 5 elements": {
			items:      []int{10, 20, 30, 40, 50},
			priorities: []float64{6, 8, 2, 7, 5},
			expectedPriorityItems: []*priorityItem[int]{
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

		if queue.Len() != len(test.expectedPriorityItems) {
			t.Fatalf("%s: item length expectancy doesn't match: expected %v, got %v", name, len(test.expectedPriorityItems), queue.Len())
		}

		popped := make([]int, queue.Len())

		for queue.Len() > 0 {
			item, _ := queue.Pop()
			popped = append(popped, item)
		}

		n := len(popped)

		for i, item := range test.expectedPriorityItems {
			poppedItem := popped[n-1-i]
			if item.value != poppedItem {
				t.Errorf("%s: item doesn't match: expected %v at index %d, got %v", name, item.value, i, poppedItem)
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

func TestPriorityQueue_UpdatePriority(t *testing.T) {
	tests := map[string]struct {
		items                 []*priorityItem[int]
		expectedPriorityItems []*priorityItem[int]
		decreaseItem          int
		decreasePriority      float64
	}{
		"decrease 30 to priority 5": {
			items: []*priorityItem[int]{
				{value: 40, priority: 40},
				{value: 30, priority: 30},
				{value: 20, priority: 20},
				{value: 10, priority: 10},
			},
			decreaseItem:     30,
			decreasePriority: 5,
			expectedPriorityItems: []*priorityItem[int]{
				{value: 40, priority: 40},
				{value: 20, priority: 20},
				{value: 10, priority: 10},
				{value: 30, priority: 5},
			},
		},
		"decrease a non-existent item": {
			items: []*priorityItem[int]{
				{value: 40, priority: 40},
				{value: 30, priority: 30},
				{value: 20, priority: 20},
				{value: 10, priority: 10},
			},
			decreaseItem:     50,
			decreasePriority: 10,
			expectedPriorityItems: []*priorityItem[int]{
				{value: 40, priority: 40},
				{value: 30, priority: 30},
				{value: 20, priority: 20},
				{value: 10, priority: 10},
			},
		},
		"increase 10 to priority 100": {
			items: []*priorityItem[int]{
				{value: 40, priority: 40},
				{value: 30, priority: 30},
				{value: 20, priority: 20},
				{value: 10, priority: 10},
			},
			decreaseItem:     10,
			decreasePriority: 100,
			expectedPriorityItems: []*priorityItem[int]{
				{value: 10, priority: 100},
				{value: 40, priority: 40},
				{value: 30, priority: 30},
				{value: 20, priority: 20},
			},
		},
	}

	for name, test := range tests {
		queue := newPriorityQueue[int]()

		for _, item := range test.items {
			queue.Push(item.value, item.priority)
		}

		queue.UpdatePriority(test.decreaseItem, test.decreasePriority)

		if queue.Len() != len(test.expectedPriorityItems) {
			t.Fatalf("%s: item length expectancy doesn't match: expected %v, got %v", name, len(test.expectedPriorityItems), queue.Len())
		}

		popped := make([]int, queue.Len())

		for queue.Len() > 0 {
			item, _ := queue.Pop()
			popped = append(popped, item)
		}

		n := len(popped)

		for i, item := range test.expectedPriorityItems {
			poppedItem := popped[n-1-i]
			if item.value != poppedItem {
				t.Errorf("%s: item doesn't match: expected %v at index %d, got %v", name, item.value, i, poppedItem)
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

		n := queue.Len()

		if n != test.expectedLen {
			t.Errorf("%s: length expectancy doesn't match: expected %v, got %v", name, test.expectedLen, n)
		}
	}
}

func TestStack_push(t *testing.T) {
	type args[T comparable] struct {
		t T
	}
	type testCase[T comparable] struct {
		name     string
		elements []int
		args     args[T]
	}
	tests := []testCase[int]{
		{
			"push 1",
			[]int{1, 2, 3, 4, 5, 6},
			args[int]{
				t: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newStack[int]()

			for _, element := range tt.elements {
				s.push(element)
			}

			s.push(tt.args.t)
		})
	}
}

func TestStack_pop(t *testing.T) {
	type testCase[T comparable] struct {
		name     string
		elements []int
		want     T
		wantErr  bool
	}
	tests := []testCase[int]{
		{
			"pop element",
			[]int{1, 2, 3, 4, 5, 6},
			6,
			false,
		},
		{
			"pop element from empty stack",
			[]int{},
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newStack[int]()

			for _, element := range tt.elements {
				s.push(element)
			}

			got, ok := s.pop()
			if ok == tt.wantErr {
				t.Errorf("pop() bool = %v, wantErr %v", ok, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pop() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStack_top(t *testing.T) {
	type testCase[T comparable] struct {
		name     string
		elements []int
		want     T
		wantErr  bool
	}
	tests := []testCase[int]{
		{
			"top element",
			[]int{1, 2, 3, 4, 5, 6},
			6,
			false,
		},
		{
			"top element of empty stack",
			[]int{},
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newStack[int]()

			for _, element := range tt.elements {
				s.push(element)
			}

			got, ok := s.top()
			if ok == tt.wantErr {
				t.Errorf("top() bool = %v, wantErr %v", ok, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("top() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStack_isEmpty(t *testing.T) {
	type testCase[T comparable] struct {
		name     string
		elements []int
		want     bool
	}
	tests := []testCase[int]{
		{
			"empty",
			[]int{},
			true,
		},
		{
			"not empty",
			[]int{1, 2, 3, 4, 5, 6},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newStack[int]()

			for _, element := range tt.elements {
				s.push(element)
			}

			if got := s.isEmpty(); got != tt.want {
				t.Errorf("isEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStack_forEach(t *testing.T) {
	type args[T comparable] struct {
		f func(T)
	}
	type testCase[T comparable] struct {
		name     string
		elements []int
		args     args[T]
	}
	tests := []testCase[int]{
		{
			name:     "forEach",
			elements: []int{1, 2, 3, 4, 5, 6},
			args: args[int]{
				f: func(i int) {
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newStack[int]()

			for _, element := range tt.elements {
				s.push(element)
			}

			s.forEach(tt.args.f)
		})
	}
}

func TestStack_contains(t *testing.T) {
	type testCase[T comparable] struct {
		name     string
		elements []int
		arg      T
		expected bool
	}
	tests := []testCase[int]{
		{
			name:     "contains 6",
			elements: []int{1, 2, 3, 4, 5, 6},
			arg:      6,
			expected: true,
		},
		{
			name:     "contains 7",
			elements: []int{1, 2, 3, 4, 5, 6},
			arg:      7,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newStack[int]()

			for _, element := range tt.elements {
				s.push(element)
			}

			got := s.contains(tt.arg)
			if got != tt.expected {
				t.Errorf("contains() = %v, want %v", got, tt.expected)
			}
		})
	}
}

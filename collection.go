package graph

import (
	"errors"
)

// priorityQueue is a priority queue implementation for minimum priorities. It maintains a list of
// items and a list of their corresponding priorities. Both are descendently ordered.
//
// For example, a priorityQueue that stores the items and priorities ("A", 5), ("B", 2), ("C", 3)
// looks as follows:
//
//	items: []string{"A", "C", "B"}
//	priorities: []int{5, 3, 2}
//
// Pulling an item from the queue will remove the least-priotized item, i.e. the last one.
type priorityQueue[T comparable] struct {
	items      []T
	priorities []int
}

func newPriorityQueue[T comparable]() *priorityQueue[T] {
	return &priorityQueue[T]{
		items:      make([]T, 0),
		priorities: make([]int, 0),
	}
}

// Push pushes a new item with the given priority into the queue. Because Push keeps track of the
// descendant order of items and priorities, it has O(n) insertion time.
func (p *priorityQueue[T]) Push(item T, priority int) {
	index := p.Len() - 1

	for i := p.Len(); i > 0; i-- {
		currentPriority := p.priorities[i-1]
		if currentPriority > priority {
			index = i
			break
		}
	}

	p.insertItemAt(item, priority, index)
}

// Pop returns the item with the smallest priority from the queue and removes that item. Returns an
// error if the priority queue is empty, which can be tested using Len first.s
func (p *priorityQueue[T]) Pop() (T, error) {
	var leastPriorityItem T

	if p.Len() == 0 {
		return leastPriorityItem, errors.New("priority queue is empty")
	}

	leastPriorityItem = p.items[p.Len()-1]

	p.items = p.items[:p.Len()-1]
	p.priorities = p.priorities[:p.Len()-1]

	return leastPriorityItem, nil
}

// Len returns the current length of the priority queue, i.e. the number of items in the queue.
func (p *priorityQueue[T]) Len() int {
	return len(p.items)
}

func (p *priorityQueue[T]) insertItemAt(item T, priority, index int) {
	if p.Len() == 0 || p.Len() == index {
		p.items = append(p.items, item)
		p.priorities = append(p.priorities, priority)
		return
	}

	p.items = append(p.items[:index+1], p.items[index:]...)
	p.priorities = append(p.priorities[:index+1], p.priorities[index:]...)

	p.items[index] = item
	p.priorities[index] = priority
}

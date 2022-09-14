package graph

import (
	"container/heap"
	"errors"
)

// priorityQueue is a priority queue implementation for minimum priorities, meaning that smaller
// values will be prioritized. It maintains a descendingly ordered list of priority items.
//
// This implementation is built on top of heap.Interface with some adjustments to comply with our
// generic usage.
type priorityQueue[T comparable] struct {
	items *priorityItems[T]

	// The map used to look up item that already pushed to the queue. Especially useful when we do
	// UpdatePriority operation.
	lookup map[T]*priorityItem[T]
}

type priorityItems[T comparable] []*priorityItem[T]

// priorityItem is an item in the priority queue, consisting of a priority and an actual value.
type priorityItem[T comparable] struct {
	value    T
	priority float64

	// The index field is used and operated internally by heap.Interface to re-organize items in the
	// queue.
	index int
}

func newPriorityQueue[T comparable]() *priorityQueue[T] {
	return &priorityQueue[T]{
		items:  &priorityItems[T]{},
		lookup: map[T]*priorityItem[T]{},
	}
}

func (p *priorityQueue[T]) Len() int {
	return p.items.Len()
}

// Push pushes a new item with the given priority into the queue.
func (p *priorityQueue[T]) Push(item T, priority float64) {
	_, ok := p.lookup[item]
	if ok {
		return
	}

	newItem := &priorityItem[T]{
		value:    item,
		priority: priority,
		index:    0,
	}
	heap.Push(p.items, newItem)
	p.lookup[item] = newItem
}

// Pop returns the item with the smallest priority from the queue and removes that item.
func (p *priorityQueue[T]) Pop() (T, error) {
	if len(*p.items) == 0 {
		var zeroVal T
		return zeroVal, errors.New("priority queue is empty")
	}

	item := heap.Pop(p.items).(*priorityItem[T])
	delete(p.lookup, item.value)
	return item.value, nil
}

// UpdatePriority updates the priority of a given item to the given priority. The item must be
// pushed into the queue first. If the item doesn't exist, nothing happens.
func (p *priorityQueue[T]) UpdatePriority(item T, priority float64) {
	targetItem, ok := p.lookup[item]
	if !ok {
		return
	}

	targetItem.priority = priority
	heap.Fix(p.items, targetItem.index)
}

// Len and the rest methods below are required to implement priority queue on top of heap.Interface.
func (pi *priorityItems[T]) Len() int {
	return len(*pi)
}

func (pi *priorityItems[T]) Less(i, j int) bool {
	return (*pi)[i].priority < (*pi)[j].priority
}

func (pi *priorityItems[T]) Swap(i, j int) {
	(*pi)[i], (*pi)[j] = (*pi)[j], (*pi)[i]
	(*pi)[i].index = i
	(*pi)[j].index = j
}

func (pi *priorityItems[T]) Push(x interface{}) {
	it := x.(*priorityItem[T])
	it.index = len(*pi)
	*pi = append(*pi, it)
}

func (pi *priorityItems[T]) Pop() interface{} {
	old := *pi
	item := old[len(old)-1]
	*pi = old[0 : len(old)-1]

	return item
}

package graph

import (
	"errors"
	"sort"
)

// priorityQueue is a priority queue implementation for minimum priorities, meaning that smaller
// values will be prioritized. It maintains an descendently ordered list of priority items.
//
// This is still a naive implementation, which is to be replaced with a binary heap implementation.
type priorityQueue[T comparable] struct {
	items []priorityItem[T]
}

// priorityItem is an item in the priority queue, consiting of a priority and an actual value.
type priorityItem[T comparable] struct {
	value    T
	priority float64
}

func newPriorityQueue[T comparable]() *priorityQueue[T] {
	return &priorityQueue[T]{
		items: make([]priorityItem[T], 0),
	}
}

// Push pushes a new item with the given priority into the queue. Because Push keeps track of the
// descendant order of items and priorities, it has O(n) insertion time.
func (p *priorityQueue[T]) Push(item T, priority float64) {
	index := p.Len() - 1

	for i := p.Len(); i > 0; i-- {
		currentItem := p.items[i-1]
		if currentItem.priority > priority {
			index = i
			break
		}
	}

	p.insertItemAt(item, priority, index)
}

// Pop returns the item with the smallest priority from the queue and removes that item. Returns an
// error if the priority queue is empty, which can be tested using Len first.s
func (p *priorityQueue[T]) Pop() (T, error) {
	var priorityItem priorityItem[T]

	if p.Len() == 0 {
		return priorityItem.value, errors.New("priority queue is empty")
	}

	priorityItem = p.items[p.Len()-1]
	p.items = p.items[:p.Len()-1]

	return priorityItem.value, nil
}

// DecreasePriority decreases the priority of a given item to the given priority. The item must be
// pushed into the queue first. If the item doesn't exist, nothing happens.
//
// With the current implementation, DecreasePriority causes the items in the queue to be re-sorted.
func (p *priorityQueue[T]) DecreasePriority(item T, priority float64) {
	for i, currentItem := range p.items {
		if currentItem.value == item {
			p.items[i].priority = priority
		}
	}

	sort.Slice(p.items, func(i, j int) bool {
		return p.items[i].priority > p.items[j].priority
	})
}

// Len returns the current length of the priority queue, i.e. the number of items in the queue.
func (p *priorityQueue[T]) Len() int {
	return len(p.items)
}

func (p *priorityQueue[T]) insertItemAt(item T, priority float64, index int) {
	priorityItem := priorityItem[T]{
		value:    item,
		priority: priority,
	}

	if p.Len() == 0 || p.Len() == index {
		p.items = append(p.items, priorityItem)
		return
	}

	p.items = append(p.items[:index+1], p.items[index:]...)

	p.items[index] = priorityItem
}

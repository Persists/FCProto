package queue

import (
	"sync"
)

// Generic Struct
type Queue[T comparable] struct {
	head *Node[T]
	tail *Node[T]

	// Mutual exclusion lock
	lock sync.Mutex
	// Cond is used to pause mulitple goroutines and wait
	cond *sync.Cond
}

// Node struct
type Node[T comparable] struct {
	value T

	next *Node[T]
}

// Initialize ConcurrentQueue
func New[T comparable]() *Queue[T] {
	q := &Queue[T]{}
	q.cond = sync.NewCond(&q.lock)
	return q
}

// Enqueue adds an item to the queue
func (q *Queue[T]) Enqueue(item T) {
	q.lock.Lock()
	defer q.lock.Unlock()

	// Create a new node
	node := &Node[T]{value: item}

	if q.head == nil {
		q.head = node
		q.tail = node
	} else {
		q.tail.next = node
		q.tail = node
	}

	q.cond.Signal()
}

// Dequeue removes an item from the queue and returns it, it blocks until an item is available
func (q *Queue[T]) Dequeue() T {
	q.lock.Lock()
	defer q.lock.Unlock()

	for q.head == nil {
		q.cond.Wait()
	}

	node := q.head
	q.head = q.head.next

	if q.head == nil {
		q.tail = nil
	}

	return node.value
}

func (q *Queue[T]) Len() int {
	q.lock.Lock()
	defer q.lock.Unlock()

	len := 0
	for node := q.head; node != nil; node = node.next {
		len++
	}

	return len
}

// IsEmpty checks if the queue is empty
func (q *Queue[T]) IsEmpty() bool {
	return q.head == nil
}

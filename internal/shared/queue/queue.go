package queue

import (
	"sync"

	"github.com/Persists/fcproto/internal/shared/models"
)

// the queue is implemented using a linked list

// Queue
type Queue struct {
	head *node
	tail *node

	mu sync.Mutex
}

// node
type node struct {
	message models.Message

	next *node
}

// NewQueue creates a new queue
func NewQueue() *Queue {
	return &Queue{}
}

// Enqueue adds a message to the queue
func (q *Queue) Enqueue(msg models.Message) {
	q.mu.Lock()
	defer q.mu.Unlock()
	n := &node{message: msg}

	if q.tail == nil {
		q.head = n
		q.tail = n
		return
	}

	q.tail.next = n
	q.tail = n
}

// Dequeue removes a message from the queue
func (q *Queue) Dequeue() (models.Message, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.head == nil {
		return models.Message{}, false
	}

	n := q.head
	q.head = n.next

	if q.head == nil {
		q.tail = nil
	}

	return n.message, true
}

// IsEmpty checks if the queue is empty
func (q *Queue) IsEmpty() bool {
	return q.head == nil
}

// Size returns the size of the queue
func (q *Queue) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	n := q.head
	size := 0
	for n != nil {
		size++
		n = n.next
	}
	return size
}

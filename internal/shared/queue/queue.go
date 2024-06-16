package queue

import "github.com/Persists/fcproto/internal/shared/models"

// the queue is implemented using a linked list

// Queue
type Queue struct {
	head *node
	tail *node
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

// Peek returns the message at the head of the queue
func (q *Queue) Peek() (models.Message, bool) {
	if q.head == nil {
		return models.Message{}, false
	}

	return q.head.message, true
}

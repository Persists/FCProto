package queue_test

import (
	"testing"

	"github.com/Persists/fcproto/internal/shared/models"
	"github.com/Persists/fcproto/internal/shared/queue"
)

func TestQueue(t *testing.T) {
	q := queue.NewQueue()

	msg1 := models.Message{Timestamp: 1, Content: "msg1"}
	msg2 := models.Message{Timestamp: 2, Content: "msg2"}
	msg3 := models.Message{Timestamp: 3, Content: "msg3"}

	q.Enqueue(msg1)
	q.Enqueue(msg2)

	m, ok := q.Peek()
	if !ok {
		t.Error("Peek failed")
	}

	if m != msg1 {
		t.Errorf("Expected message %v, got %v", msg1, m)
	}

	m, ok = q.Dequeue()
	if !ok {
		t.Error("Dequeue failed")
	}

	if m != msg1 {
		t.Errorf("Expected message %v, got %v", msg1, m)
	}

	m, ok = q.Dequeue()
	if !ok {
		t.Error("Dequeue failed")
	}

	if m != msg2 {
		t.Errorf("Expected message %v, got %v", msg2, m)
	}

	_, ok = q.Dequeue()
	if ok {
		t.Error("Dequeue should have failed")
	}

	q.Enqueue(msg3)

	m, ok = q.Dequeue()
	if !ok {
		t.Error("Dequeue failed")
	}

	if m != msg3 {
		t.Errorf("Expected message %v, got %v", msg3, m)
	}

	_, ok = q.Peek()
	if ok {
		t.Error("Peek should have failed")
	}
}

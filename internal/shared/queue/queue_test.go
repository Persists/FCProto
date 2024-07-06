package queue

import (
	"testing"
	"time"

	"github.com/Persists/fcproto/internal/shared/models"
	"github.com/Persists/fcproto/internal/shared/queue"
)

func TestQueueQueuesItems(t *testing.T) {
	q := queue.NewQueue[models.SensorMessage]()

	q.Enqueue(models.SensorMessage{Timestamp: 1, Content: "msg1"})
	q.Enqueue(models.SensorMessage{Timestamp: 2, Content: "msg2"})
	q.Enqueue(models.SensorMessage{Timestamp: 3, Content: "msg3"})

	time.Sleep(50 * time.Millisecond)

	item := q.Dequeue()
	if item.Timestamp != 1 || item.Content != "msg1" {
		t.Errorf("Expected msg1 but got %v", item)
	}

	item = q.Dequeue()
	if item.Timestamp != 2 || item.Content != "msg2" {
		t.Errorf("Expected msg2 but got %v", item)
	}

	item = q.Dequeue()
	if item.Timestamp != 3 || item.Content != "msg3" {
		t.Errorf("Expected msg3 but got %v", item)
	}

	go func() {
		time.Sleep(50 * time.Millisecond)
		q.Enqueue(models.SensorMessage{Timestamp: 0, Content: ""})
	}()

	item = q.Dequeue()
	if item.Timestamp != 0 || item.Content != "" {
		t.Errorf("Expected empty message but got %v", item)
	}
}

func TestQueueBlockingDequeue(t *testing.T) {
	q := queue.NewQueue[models.SensorMessage]()

	// Channel to signal when Dequeue has returned a value
	done := make(chan struct{})

	go func() {
		q.Dequeue()
		close(done)
	}()

	go func() {
		time.Sleep(50 * time.Millisecond)
		q.Enqueue(models.SensorMessage{Timestamp: 4, Content: "msg4"})
	}()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Error("Dequeue did not unblock after an item was enqueued")
	}
}

func TestQueueEnqueueUnblocksDequeue(t *testing.T) {
	q := queue.NewQueue[models.SensorMessage]()

	// Channel to signal when Dequeue has returned a value
	done := make(chan struct{})

	go func() {
		item := q.Dequeue() // This should block until an item is enqueued
		if item.Timestamp != 1 || item.Content != "msg1" {
			t.Errorf("Expected msg1 but got %v", item)
		}
		close(done) // Close the channel to signal that Dequeue has returned
	}()

	time.Sleep(50 * time.Millisecond)

	// Enqueue an item which should unblock the Dequeue call
	q.Enqueue(models.SensorMessage{Timestamp: 1, Content: "msg1"})

	select {
	case <-done:
		// Test passes, Dequeue unblocked as expected
	case <-time.After(200 * time.Millisecond):
		t.Error("Dequeue did not unblock after an item was enqueued")
	}
}

func BenchmarkQueueEnqueue(b *testing.B) {
	q := queue.NewQueue[models.SensorMessage]()

	count := 0

	for i := 0; i < b.N; i++ {
		q.Enqueue(models.SensorMessage{Timestamp: int64(i), Content: "msg"})

		count++
	}

	b.Log("count:", count)

	for !q.IsEmpty() {
		q.Dequeue()
		count--
	}

	if count != 0 {
		b.Errorf("Expected count to be 0 but got %d", count)
	}

	b.ReportAllocs()
}

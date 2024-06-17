package queue_test

import (
	"testing"
	"time"

	"github.com/Persists/fcproto/internal/shared/models"
	"github.com/Persists/fcproto/internal/shared/queue"
)

func TestQueueQueuesItems(t *testing.T) {
	q := queue.NewQueue[models.Message]()

	q.Enqueue(models.Message{Timestamp: 1, Content: "msg1"})
	q.Enqueue(models.Message{Timestamp: 2, Content: "msg2"})
	q.Enqueue(models.Message{Timestamp: 3, Content: "msg3"})

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
		q.Enqueue(models.Message{Timestamp: 0, Content: ""})
	}()

	item = q.Dequeue()
	if item.Timestamp != 0 || item.Content != "" {
		t.Errorf("Expected empty message but got %v", item)
	}
}

func TestQueueBlockingDequeue(t *testing.T) {
	q := queue.NewQueue[models.Message]()

	// Channel to signal when Dequeue has returned a value
	done := make(chan struct{})

	go func() {
		q.Dequeue()
		close(done)
	}()

	go func() {
		time.Sleep(50 * time.Millisecond)
		q.Enqueue(models.Message{Timestamp: 4, Content: "msg4"})
	}()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Error("Dequeue did not unblock after an item was enqueued")
	}
}

func TestQueueEnqueueUnblocksDequeue(t *testing.T) {
	q := queue.NewQueue[models.Message]()

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
	q.Enqueue(models.Message{Timestamp: 1, Content: "msg1"})

	select {
	case <-done:
		// Test passes, Dequeue unblocked as expected
	case <-time.After(200 * time.Millisecond):
		t.Error("Dequeue did not unblock after an item was enqueued")
	}
}

func TestQueueMultipleBlockingDequeue(t *testing.T) {
	q := queue.NewQueue[models.Message]()

	// Create channels to signal when Dequeue operations return values
	done1 := make(chan struct{})
	done2 := make(chan struct{})
	done3 := make(chan struct{})
	done4 := make(chan struct{})
	done5 := make(chan struct{})

	go func() {
		item := q.Dequeue() // This should block until an item is enqueued
		if item.Timestamp != 1 || item.Content != "msg1" {
			t.Errorf("Expected msg1 but got %v", item)
		}
		close(done1) // Signal Dequeue1 completion
	}()

	go func() {
		item := q.Dequeue() // This should block until a second item is enqueued
		if item.Timestamp != 2 || item.Content != "msg2" {
			t.Errorf("Expected msg2 but got %v", item)
		}
		close(done2) // Signal Dequeue2 completion
	}()

	go func() {
		item := q.Dequeue() // This should block until a third item is enqueued
		if item.Timestamp != 3 || item.Content != "msg3" {
			t.Errorf("Expected msg3 but got %v", item)
		}
		close(done3) // Signal Dequeue3 completion
	}()

	go func() {
		item := q.Dequeue() // This should block until a fourth item is enqueued
		if item.Timestamp != 4 || item.Content != "msg4" {
			t.Errorf("Expected msg4 but got %v", item)
		}
		close(done4) // Signal Dequeue4 completion
	}()

	go func() {
		item := q.Dequeue() // This should block until a fifth item is enqueued
		if item.Timestamp != 5 || item.Content != "msg5" {
			t.Errorf("Expected msg5 but got %v", item)
		}
		close(done5) // Signal Dequeue5 completion
	}()

	time.Sleep(50 * time.Millisecond) // Give time for Dequeue to block

	// Enqueue first item, should unblock the first Dequeue
	q.Enqueue(models.Message{Timestamp: 1, Content: "msg1"})

	select {
	case <-done1:
		// First Dequeue unblocked as expected
	case <-time.After(200 * time.Millisecond):
		t.Error("First Dequeue did not unblock after the first item was enqueued")
	}

	select {
	case <-done2:
		t.Error("Second Dequeue unblocked before the second item was enqueued")
	case <-time.After(50 * time.Millisecond):
		// Second Dequeue still blocked as expected
	}

	select {
	case <-done3:
		t.Error("Second Dequeue unblocked before the second item was enqueued")
	case <-time.After(50 * time.Millisecond):
		// Second Dequeue still blocked as expected
	}

	select {
	case <-done4:
		t.Error("Second Dequeue unblocked before the second item was enqueued")
	case <-time.After(50 * time.Millisecond):
		// Second Dequeue still blocked as expected
	}

	select {
	case <-done5:
		t.Error("Second Dequeue unblocked before the second item was enqueued")
	case <-time.After(50 * time.Millisecond):
		// Second Dequeue still blocked as expected
	}

	// Enqueue second item, should unblock the second Dequeue
	q.Enqueue(models.Message{Timestamp: 2, Content: "msg2"})

	select {
	case <-done2:
		// Second Dequeue unblocked as expected
	case <-time.After(200 * time.Millisecond):
		t.Error("Second Dequeue did not unblock after the second item was enqueued")
	}
}

func BenchmarkQueueEnqueue(b *testing.B) {
	q := queue.NewQueue[models.Message]()

	count := 0

	for i := 0; i < b.N; i++ {
		q.Enqueue(models.Message{Timestamp: int64(i), Content: "msg"})

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

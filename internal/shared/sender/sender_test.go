package sender

import (
	"github.com/Persists/fcproto/internal/shared/models"
	"testing"
)

// MockQueue is a mock implementation of Queue.
type MockQueue struct {
	Messages []models.Message
}

func (q *MockQueue) Enqueue(msg models.Message) {
	q.Messages = append(q.Messages, msg)
}

func (q *MockQueue) Dequeue() models.Message {
	if len(q.Messages) == 0 {
		return models.Message{}
	}
	msg := q.Messages[0]
	q.Messages = q.Messages[1:]
	return msg
}

// MockUtils is a mock implementation of Utils.
type MockUtils struct{}

func TestSender_Send(t *testing.T) {
	mockQueue := &MockQueue{}
	mockConfig := &MockConfig{}
	mockUtils := &MockUtils{}

	sender := NewSender(mockQueue, mockConfig, mockUtils)

	msg := models.Message{
		// Initialize message as needed for test.
	}

	sender.Send(msg)

	if len(mockQueue.Messages) != 1 {
		t.Errorf("Expected queue length 1, got %d", len(mockQueue.Messages))
	}

	// Additional assertions as needed.
}

func TestSender_Connect(t *testing.T) {
	mockQueue := &MockQueue{}
	mockConfig := &MockConfig{}
	mockUtils := &MockUtils{}

	sender := NewSender(mockQueue, mockConfig, mockUtils)

	err := sender.Connect()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Additional assertions as needed.
}

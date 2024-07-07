package models

import "testing"

func TestNewHeartbeatMessage(t *testing.T) {
	addr := "1234"
	heartbeatMessage := NewHeartbeatMessage(addr)

	if heartbeatMessage.CallbackPort != addr {
		t.Errorf("Expected %s, got %s", addr, heartbeatMessage.CallbackPort)
	}
}

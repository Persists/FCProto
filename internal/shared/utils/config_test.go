package utils

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	key := "TEST"
	fallback := "fallback"
	value := GetEnv(key, fallback)

	if value != fallback {
		t.Errorf("Expected %s but got %s", fallback, value)
	}

	os.Setenv(key, "test")
	value = GetEnv(key, fallback)

	if value != "test" {
		t.Errorf("Expected test but got %s", value)
	}

	os.Unsetenv(key)
}

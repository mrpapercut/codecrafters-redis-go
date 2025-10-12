package commands

import (
	"testing"
	"time"
)

func TestHandleSet(t *testing.T) {
	message := []byte("*3\r\n$3\r\nSET\r\n$6\r\nmy-key\r\n$8\r\nmy-value\r\n")
	expected := "+OK\r\n"

	response := HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	message = []byte("*2\r\n$3\r\nGET\r\n$6\r\nmy-key\r\n")
	expected = "$8\r\nmy-value\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

func TestHandleSetWithSecondExpiration(t *testing.T) {
	message := []byte("*5\r\n$3\r\nSET\r\n$17\r\nmy-expiration-key\r\n$8\r\nmy-value\r\n$2\r\nEX\r\n$1\r\n1\r\n")
	expected := "+OK\r\n"

	response := HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	message = []byte("*2\r\n$3\r\nGET\r\n$17\r\nmy-expiration-key\r\n")
	expected = "$8\r\nmy-value\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	time.Sleep(time.Duration(1) * time.Second)

	message = []byte("*2\r\n$3\r\nGET\r\n$17\r\nmy-expiration-key\r\n")
	expected = "$-1\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

func TestHandleSetWithMilliSecondExpiration(t *testing.T) {
	message := []byte("*5\r\n$3\r\nSET\r\n$19\r\nmy-expiration-key-2\r\n$8\r\nmy-value\r\n$2\r\nPX\r\n$3\r\n500\r\n")
	expected := "+OK\r\n"

	response := HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	message = []byte("*2\r\n$3\r\nGET\r\n$19\r\nmy-expiration-key-2\r\n")
	expected = "$8\r\nmy-value\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	time.Sleep(time.Duration(750) * time.Millisecond)

	message = []byte("*2\r\n$3\r\nGET\r\n$19\r\nmy-expiration-key-2\r\n")
	expected = "$-1\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

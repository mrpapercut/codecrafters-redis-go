package commands

import (
	"testing"
	"time"
)

func TestHandleSet(t *testing.T) {
	// SET set_key value
	message := []byte("*3\r\n$3\r\nSET\r\n$7\r\nset_key\r\n$5\r\nvalue\r\n")
	expected := "+OK\r\n"

	response := HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// GET set_key
	message = []byte("*2\r\n$3\r\nGET\r\n$7\r\nset_key\r\n")
	expected = "$5\r\nvalue\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

func TestHandleSetWithSecondExpiration(t *testing.T) {
	// SET set_key_exp_sec value EX 1
	message := []byte("*5\r\n$3\r\nSET\r\n$15\r\nset_key_exp_sec\r\n$5\r\nvalue\r\n$2\r\nEX\r\n$1\r\n1\r\n")
	expected := "+OK\r\n"

	response := HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// GET set_key_exp_sec
	message = []byte("*2\r\n$3\r\nGET\r\n$15\r\nset_key_exp_sec\r\n")
	expected = "$5\r\nvalue\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	time.Sleep(time.Duration(1) * time.Second)

	// GET set_key_exp_sec
	message = []byte("*2\r\n$3\r\nGET\r\n$15\r\nset_key_exp_sec\r\n")
	expected = "$-1\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

func TestHandleSetWithMilliSecondExpiration(t *testing.T) {
	// SET set_key_exp_ms value PX 500
	message := []byte("*5\r\n$3\r\nSET\r\n$14\r\nset_key_exp_ms\r\n$5\r\nvalue\r\n$2\r\nPX\r\n$3\r\n500\r\n")
	expected := "+OK\r\n"

	response := HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// GET set_key_exp_ms
	message = []byte("*2\r\n$3\r\nGET\r\n$14\r\nset_key_exp_ms\r\n")
	expected = "$5\r\nvalue\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	time.Sleep(time.Duration(750) * time.Millisecond)

	// GET set_key_exp_ms
	message = []byte("*2\r\n$3\r\nGET\r\n$14\r\nset_key_exp_ms\r\n")
	expected = "$-1\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

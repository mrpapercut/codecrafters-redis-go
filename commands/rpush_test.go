package commands

import "testing"

func TestHandleRPUSH(t *testing.T) {
	// RPUSH rpush_key element
	message := []byte("*3\r\n$5\r\nRPUSH\r\n$9\r\nrpush_key\r\n$7\r\nelement\r\n")
	expected := ":1\r\n"

	response := HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// RPUSH rpush_key "next element"
	message = []byte("*3\r\n$5\r\nRPUSH\r\n$9\r\nrpush_key\r\n$12\r\nnext element\r\n")
	expected = ":2\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

func TestHandleRPUSHMultipleElement(t *testing.T) {
	// RPUSH rpush_key_multi "element 1" "element 2" "element 3"
	message := []byte("*5\r\n$5\r\nRPUSH\r\n$15\r\nrpush_key_multi\r\n$9\r\nelement 1\r\n$9\r\nelement 2\r\n$9\r\nelement 3\r\n")
	expected := ":3\r\n"

	response := HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// RPUSH rpush_key_multi "element 4" "element 5" "element 6"
	message = []byte("*5\r\n$5\r\nRPUSH\r\n$15\r\nrpush_key_multi\r\n$9\r\nelement 4\r\n$9\r\nelement 5\r\n$9\r\nelement 6\r\n")
	expected = ":6\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

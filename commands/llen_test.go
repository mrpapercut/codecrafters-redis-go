package commands

import "testing"

func TestHandleLLEN(t *testing.T) {
	// RPUSH llen_key element
	message := []byte("*3\r\n$5\r\nRPUSH\r\n$8\r\nllen_key\r\n$7\r\nelement\r\n")
	expected := ":1\r\n"

	response := HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// RPUSH llen_key "next element"
	message = []byte("*3\r\n$5\r\nRPUSH\r\n$8\r\nllen_key\r\n$12\r\nnext element\r\n")
	expected = ":2\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// LLEN llen_key
	message = []byte("*2\r\n$4\r\nLLEN\r\n$8\r\nllen_key\r\n")
	expected = ":2\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

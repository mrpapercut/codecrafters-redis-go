package commands

import "testing"

func TestHandleRPUSH(t *testing.T) {
	message := []byte("*3\r\n$5\r\nRPUSH\r\n$8\r\nlist_key\r\n$7\r\nelement\r\n")
	expected := ":1\r\n"

	response, err := HandleCommand(message)
	if err != nil {
		t.Fatalf("error handling RPUSH command: %v", err)
	}

	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	message = []byte("*3\r\n$5\r\nRPUSH\r\n$8\r\nlist_key\r\n$12\r\nnext element\r\n")
	expected = ":2\r\n"

	response, err = HandleCommand(message)
	if err != nil {
		t.Fatalf("error handling RPUSH command: %v", err)
	}

	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

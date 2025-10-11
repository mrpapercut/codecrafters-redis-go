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

func TestHandleRPUSHMultipleElement(t *testing.T) {
	message := []byte("*5\r\n$5\r\nRPUSH\r\n$14\r\nmulti_list_key\r\n$9\r\nelement 1\r\n$9\r\nelement 2\r\n$9\r\nelement 3\r\n")
	expected := ":3\r\n"

	response, err := HandleCommand(message)
	if err != nil {
		t.Fatalf("error handling RPUSH command: %v", err)
	}

	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	message = []byte("*5\r\n$5\r\nRPUSH\r\n$14\r\nmulti_list_key\r\n$9\r\nelement 4\r\n$9\r\nelement 5\r\n$9\r\nelement 6\r\n")
	expected = ":6\r\n"

	response, err = HandleCommand(message)
	if err != nil {
		t.Fatalf("error handling RPUSH command: %v", err)
	}

	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

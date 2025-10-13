package commands

import "testing"

func TestHandleLRange(t *testing.T) {
	message := []byte("*7\r\n$5\r\nRPUSH\r\n$10\r\nlrange_key\r\n$1\r\na\r\n$1\r\nb\r\n$1\r\nc\r\n$1\r\nd\r\n$1\r\ne\r\n")
	expected := ":5\r\n"

	response := HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	message = []byte("*4\r\n$6\r\nLRANGE\r\n$10\r\nlrange_key\r\n$1\r\n0\r\n$1\r\n1\r\n")
	expected = "*2\r\n$1\r\na\r\n$1\r\nb\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	message = []byte("*4\r\n$6\r\nLRANGE\r\n$10\r\nlrange_key\r\n$1\r\n2\r\n$1\r\n4\r\n")
	expected = "*3\r\n$1\r\nc\r\n$1\r\nd\r\n$1\r\ne\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// Non-existent list
	message = []byte("*4\r\n$6\r\nLRANGE\r\n$12\r\nlrange_key_2\r\n$1\r\n0\r\n$1\r\n1\r\n")
	expected = "*0\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// 'start' greater than list length
	message = []byte("*4\r\n$6\r\nLRANGE\r\n$10\r\nlrange_key\r\n$1\r\n6\r\n$1\r\n7\r\n")
	expected = "*0\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// 'start' greater than 'stop'
	message = []byte("*4\r\n$6\r\nLRANGE\r\n$10\r\nlrange_key\r\n$1\r\n2\r\n$1\r\n1\r\n")
	expected = "*0\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// 'stop' greater than list length
	message = []byte("*4\r\n$6\r\nLRANGE\r\n$10\r\nlrange_key\r\n$1\r\n0\r\n$1\r\n6\r\n")
	expected = "*5\r\n$1\r\na\r\n$1\r\nb\r\n$1\r\nc\r\n$1\r\nd\r\n$1\r\ne\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

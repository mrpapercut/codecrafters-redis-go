package commands

import "testing"

func TestHandleLRange(t *testing.T) {
	// RPUSH lrange_key a b c d e
	message := []byte("*7\r\n$5\r\nRPUSH\r\n$10\r\nlrange_key\r\n$1\r\na\r\n$1\r\nb\r\n$1\r\nc\r\n$1\r\nd\r\n$1\r\ne\r\n")
	expected := ":5\r\n"

	response := HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// LRANGE lrange_key 0 1
	message = []byte("*4\r\n$6\r\nLRANGE\r\n$10\r\nlrange_key\r\n$1\r\n0\r\n$1\r\n1\r\n")
	expected = "*2\r\n$1\r\na\r\n$1\r\nb\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// LRANGE lrange_key 2 4
	message = []byte("*4\r\n$6\r\nLRANGE\r\n$10\r\nlrange_key\r\n$1\r\n2\r\n$1\r\n4\r\n")
	expected = "*3\r\n$1\r\nc\r\n$1\r\nd\r\n$1\r\ne\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// Non-existent list
	// LRANGE lrange_key_2 0 1
	message = []byte("*4\r\n$6\r\nLRANGE\r\n$12\r\nlrange_key_2\r\n$1\r\n0\r\n$1\r\n1\r\n")
	expected = "*0\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// 'start' greater than list length
	// LRANGE lrange_key 6 7
	message = []byte("*4\r\n$6\r\nLRANGE\r\n$10\r\nlrange_key\r\n$1\r\n6\r\n$1\r\n7\r\n")
	expected = "*0\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// 'start' greater than 'stop'
	// LRANGE lrange_key 2 1
	message = []byte("*4\r\n$6\r\nLRANGE\r\n$10\r\nlrange_key\r\n$1\r\n2\r\n$1\r\n1\r\n")
	expected = "*0\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// 'stop' greater than list length
	// LRANGE lrange_key 0 6
	message = []byte("*4\r\n$6\r\nLRANGE\r\n$10\r\nlrange_key\r\n$1\r\n0\r\n$1\r\n6\r\n")
	expected = "*5\r\n$1\r\na\r\n$1\r\nb\r\n$1\r\nc\r\n$1\r\nd\r\n$1\r\ne\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

func TestHandleLRangeNegativeIndex(t *testing.T) {
	// RPUSH lrange_neg a b c d e
	message := []byte("*7\r\n$5\r\nRPUSH\r\n$14\r\nlrange_key_neg\r\n$1\r\na\r\n$1\r\nb\r\n$1\r\nc\r\n$1\r\nd\r\n$1\r\ne\r\n")
	expected := ":5\r\n"

	response := HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// LRANGE lrange_neg -2 -1
	message = []byte("*4\r\n$6\r\nLRANGE\r\n$14\r\nlrange_key_neg\r\n$2\r\n-2\r\n$2\r\n-1\r\n")
	expected = "*2\r\n$1\r\nd\r\n$1\r\ne\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// LRANGE lrange_neg 0 -3
	message = []byte("*4\r\n$6\r\nLRANGE\r\n$14\r\nlrange_key_neg\r\n$1\r\n0\r\n$2\r\n-3\r\n")
	expected = "*3\r\n$1\r\na\r\n$1\r\nb\r\n$1\r\nc\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// LRANGE lrange_neg 0 -6
	message = []byte("*4\r\n$6\r\nLRANGE\r\n$14\r\nlrange_key_neg\r\n$1\r\n0\r\n$2\r\n-6\r\n")
	expected = "*0\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// LRANGE lrange_neg -6 2
	message = []byte("*4\r\n$6\r\nLRANGE\r\n$14\r\nlrange_key_neg\r\n$2\r\n-6\r\n$1\r\n2\r\n")
	expected = "*3\r\n$1\r\na\r\n$1\r\nb\r\n$1\r\nc\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

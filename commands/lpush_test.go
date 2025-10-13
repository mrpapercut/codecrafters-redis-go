package commands

import "testing"

func TestHandleLPush(t *testing.T) {
	// LPUSH lpush_key a
	message := []byte("*3\r\n$5\r\nLPUSH\r\n$9\r\nlpush_key\r\n$1\r\na\r\n")
	expected := ":1\r\n"

	response := HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// LPUSH lpush_key b
	message = []byte("*3\r\n$5\r\nLPUSH\r\n$9\r\nlpush_key\r\n$1\r\nb\r\n")
	expected = ":2\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// LLEN lpush_key
	message = []byte("*2\r\n$4\r\nLLEN\r\n$9\r\nlpush_key\r\n")
	expected = ":2\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// LRANGE lpush_key 0 -1
	message = []byte("*4\r\n$6\r\nLRANGE\r\n$9\r\nlpush_key\r\n$1\r\n0\r\n$2\r\n-1\r\n")
	expected = "*2\r\n$1\r\nb\r\n$1\r\na\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

func TestHandleLPushWithMultipleArguments(t *testing.T) {
	// LPUSH lpush_key_multi a b c
	message := []byte("*5\r\n$5\r\nLPUSH\r\n$15\r\nlpush_key_multi\r\n$1\r\na\r\n$1\r\nb\r\n$1\r\nc\r\n")
	expected := ":3\r\n"

	response := HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// LLEN lpush_key_multi
	message = []byte("*2\r\n$4\r\nLLEN\r\n$15\r\nlpush_key_multi\r\n")
	expected = ":3\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// LRANGE lpush_key_multi 0 -1
	message = []byte("*4\r\n$6\r\nLRANGE\r\n$15\r\nlpush_key_multi\r\n$1\r\n0\r\n$2\r\n-1\r\n")
	expected = "*3\r\n$1\r\nc\r\n$1\r\nb\r\n$1\r\na\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

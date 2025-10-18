package commands

import (
	"testing"
)

func TestHandleLPOP(t *testing.T) {
	// RPUSH lpop_key a b c d
	message := []byte("*6\r\n$5\r\nRPUSH\r\n$8\r\nlpop_key\r\n$1\r\na\r\n$1\r\nb\r\n$1\r\nc\r\n$1\r\nd\r\n")
	expected := ":4\r\n"

	response := HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// LPOP lpop_key
	message = []byte("*2\r\n$4\r\nLPOP\r\n$8\r\nlpop_key\r\n")
	expected = "$1\r\na\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// LLEN lpop_key
	message = []byte("*2\r\n$4\r\nLLEN\r\n$8\r\nlpop_key\r\n")
	expected = ":3\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// LRANGE lpop_key 0 -1
	message = []byte("*4\r\n$6\r\nLRANGE\r\n$8\r\nlpop_key\r\n$1\r\n0\r\n$2\r\n-1\r\n")
	expected = "*3\r\n$1\r\nb\r\n$1\r\nc\r\n$1\r\nd\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

func TestHandleLPOPMultiple(t *testing.T) {
	// RPUSH lpop_key_multi a b c d
	message := []byte("*6\r\n$5\r\nRPUSH\r\n$14\r\nlpop_key_multi\r\n$1\r\na\r\n$1\r\nb\r\n$1\r\nc\r\n$1\r\nd\r\n")
	expected := ":4\r\n"

	response := HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// LPOP lpop_key_multi 2
	message = []byte("*3\r\n$4\r\nLPOP\r\n$14\r\nlpop_key_multi\r\n$1\r\n2\r\n")
	expected = "*2\r\n$1\r\na\r\n$1\r\nb\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// LLEN lpop_key_multi
	message = []byte("*2\r\n$4\r\nLLEN\r\n$14\r\nlpop_key_multi\r\n")
	expected = ":2\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// LRANGE lpop_key_multi 0 -1
	message = []byte("*4\r\n$6\r\nLRANGE\r\n$14\r\nlpop_key_multi\r\n$1\r\n0\r\n$2\r\n-1\r\n")
	expected = "*2\r\n$1\r\nc\r\n$1\r\nd\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

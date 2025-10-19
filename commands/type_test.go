package commands

import "testing"

func TestHandleTypeNone(t *testing.T) {
	// TYPE type_key_nonexistant
	message := []byte("*2\r\n$4\r\nTYPE\r\n$20\r\ntype_key_nonexistant\r\n")
	expected := "+none\r\n"

	response := HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

func TestHandleTypeString(t *testing.T) {
	// SET type_key_string value
	message := []byte("*3\r\n$3\r\nSET\r\n$15\r\ntype_key_string\r\n$5\r\nvalue\r\n")
	expected := "+OK\r\n"

	response := HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// TYPE type_key_string
	message = []byte("*2\r\n$4\r\nTYPE\r\n$15\r\ntype_key_string\r\n")
	expected = "+string\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

func TestHandleTypeList(t *testing.T) {
	// LPUSH type_key_list type_value
	message := []byte("*3\r\n$5\r\nLPUSH\r\n$13\r\ntype_key_list\r\n$10\r\ntype_value\r\n")
	expected := ":1\r\n"

	response := HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// TYPE type_key_list
	message = []byte("*2\r\n$4\r\nTYPE\r\n$13\r\ntype_key_list\r\n")
	expected = "+list\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

func TestHandleTypeStream(t *testing.T) {
	// XADD type_key_stream 1-1 type_key type_value
	message := []byte("*5\r\n$4\r\nXADD\r\n$15\r\ntype_key_stream\r\n$3\r\n1-1\r\n$8\r\ntype_key\r\n$10\r\ntype_value\r\n")
	expected := "$3\r\n1-1\r\n"

	response := HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// TYPE type_key_stream
	message = []byte("*2\r\n$4\r\nTYPE\r\n$15\r\ntype_key_stream\r\n")
	expected = "+stream\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

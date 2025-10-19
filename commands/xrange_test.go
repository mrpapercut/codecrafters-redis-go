package commands

import "testing"

func TestHandleXRANGE(t *testing.T) {
	// XADD xrange_key 0-1 key1 value1
	message := []byte("*5\r\n$4\r\nXADD\r\n$10\r\nxrange_key\r\n$3\r\n0-1\r\n$4\r\nkey1\r\n$6\r\nvalue1\r\n")
	expected := "$3\r\n0-1\r\n"

	response := HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// XADD xrange_key 0-2 key2 value2
	message = []byte("*5\r\n$4\r\nXADD\r\n$10\r\nxrange_key\r\n$3\r\n0-2\r\n$4\r\nkey2\r\n$6\r\nvalue2\r\n")
	expected = "$3\r\n0-2\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// XADD xrange_key 0-3 key3 value3
	message = []byte("*5\r\n$4\r\nXADD\r\n$10\r\nxrange_key\r\n$3\r\n0-3\r\n$4\r\nkey3\r\n$6\r\nvalue3\r\n")
	expected = "$3\r\n0-3\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// XRANGE xrange_key 0-2 0-3
	message = []byte("*4\r\n$6\r\nXRANGE\r\n$10\r\nxrange_key\r\n$3\r\n0-2\r\n$3\r\n0-3\r\n")
	expected = "*2\r\n*2\r\n$3\r\n0-2\r\n*2\r\n$4\r\nkey2\r\n$6\r\nvalue2\r\n*2\r\n$3\r\n0-3\r\n*2\r\n$4\r\nkey3\r\n$6\r\nvalue3\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

func TestHandleXRANGEWithIncompleteIDs(t *testing.T) {
	// XADD xrange_key_incomplete 0-1 key1 value1
	message := []byte("*5\r\n$4\r\nXADD\r\n$21\r\nxrange_key_incomplete\r\n$3\r\n0-1\r\n$4\r\nkey1\r\n$6\r\nvalue1\r\n")
	expected := "$3\r\n0-1\r\n"

	response := HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// XADD xrange_key_incomplete 0-2 key2 value2
	message = []byte("*5\r\n$4\r\nXADD\r\n$21\r\nxrange_key_incomplete\r\n$3\r\n0-2\r\n$4\r\nkey2\r\n$6\r\nvalue2\r\n")
	expected = "$3\r\n0-2\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// XADD xrange_key_incomplete 1-0 key3 value3
	message = []byte("*5\r\n$4\r\nXADD\r\n$21\r\nxrange_key_incomplete\r\n$3\r\n1-0\r\n$4\r\nkey3\r\n$6\r\nvalue3\r\n")
	expected = "$3\r\n1-0\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// XRANGE xrange_key_incomplete 0 1
	message = []byte("*4\r\n$6\r\nXRANGE\r\n$21\r\nxrange_key_incomplete\r\n$1\r\n0\r\n$1\r\n1\r\n")
	expected = "*3\r\n*2\r\n$3\r\n0-1\r\n*2\r\n$4\r\nkey1\r\n$6\r\nvalue1\r\n*2\r\n$3\r\n0-2\r\n*2\r\n$4\r\nkey2\r\n$6\r\nvalue2\r\n*2\r\n$3\r\n1-0\r\n*2\r\n$4\r\nkey3\r\n$6\r\nvalue3\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

func TestHandleXRANGEWithSpecialIDs(t *testing.T) {
	// XADD xrange_key_special 0-1 key1 value1
	message := []byte("*5\r\n$4\r\nXADD\r\n$18\r\nxrange_key_special\r\n$3\r\n0-1\r\n$4\r\nkey1\r\n$6\r\nvalue1\r\n")
	expected := "$3\r\n0-1\r\n"

	response := HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// XADD xrange_key_special 0-2 key2 value2
	message = []byte("*5\r\n$4\r\nXADD\r\n$18\r\nxrange_key_special\r\n$3\r\n0-2\r\n$4\r\nkey2\r\n$6\r\nvalue2\r\n")
	expected = "$3\r\n0-2\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// XADD xrange_key_special 1-0 key3 value3
	message = []byte("*5\r\n$4\r\nXADD\r\n$18\r\nxrange_key_special\r\n$3\r\n1-0\r\n$4\r\nkey3\r\n$6\r\nvalue3\r\n")
	expected = "$3\r\n1-0\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// XRANGE xrange_key_special - +
	message = []byte("*4\r\n$6\r\nXRANGE\r\n$18\r\nxrange_key_special\r\n$1\r\n-\r\n$1\r\n+\r\n")
	expected = "*3\r\n*2\r\n$3\r\n0-1\r\n*2\r\n$4\r\nkey1\r\n$6\r\nvalue1\r\n*2\r\n$3\r\n0-2\r\n*2\r\n$4\r\nkey2\r\n$6\r\nvalue2\r\n*2\r\n$3\r\n1-0\r\n*2\r\n$4\r\nkey3\r\n$6\r\nvalue3\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

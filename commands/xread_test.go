package commands

import (
	"testing"
)

func TestHandleXREAD(t *testing.T) {
	// XADD xread_key 0-1 temperature 95
	message := []byte("*5\r\n$4\r\nXADD\r\n$9\r\nxread_key\r\n$3\r\n0-1\r\n$11\r\ntemperature\r\n$2\r\n95\r\n")
	expected := "$3\r\n0-1\r\n"

	response := HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// XADD xread_key 0-2 temperature 97
	message = []byte("*5\r\n$4\r\nXADD\r\n$9\r\nxread_key\r\n$3\r\n0-2\r\n$11\r\ntemperature\r\n$2\r\n97\r\n")
	expected = "$3\r\n0-2\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// XREAD streams xread_key 0-0
	message = []byte("*4\r\n$5\r\nXREAD\r\n$7\r\nstreams\r\n$9\r\nxread_key\r\n$3\r\n0-0\r\n")
	expected = "*1\r\n*2\r\n$9\r\nxread_key\r\n*2\r\n*2\r\n$3\r\n0-1\r\n*2\r\n$11\r\ntemperature\r\n$2\r\n95\r\n*2\r\n$3\r\n0-2\r\n*2\r\n$11\r\ntemperature\r\n$2\r\n97\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// XREAD streams stream_key 0-1
	message = []byte("*4\r\n$5\r\nXREAD\r\n$7\r\nstreams\r\n$9\r\nxread_key\r\n$3\r\n0-1\r\n")
	expected = "*1\r\n*2\r\n$9\r\nxread_key\r\n*1\r\n*2\r\n$3\r\n0-2\r\n*2\r\n$11\r\ntemperature\r\n$2\r\n97\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

func TestHandleXREADMultipleStreams(t *testing.T) {
	// XADD xread_key_multi 0-1 temperature 95
	message := []byte("*5\r\n$4\r\nXADD\r\n$15\r\nxread_key_multi\r\n$3\r\n0-1\r\n$11\r\ntemperature\r\n$2\r\n95\r\n")
	expected := "$3\r\n0-1\r\n"

	response := HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// XADD xread_key_multi_2 0-2 humidity 97
	message = []byte("*5\r\n$4\r\nXADD\r\n$17\r\nxread_key_multi_2\r\n$3\r\n0-2\r\n$8\r\nhumidity\r\n$2\r\n97\r\n")
	expected = "$3\r\n0-2\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// XREAD streams xread_key_multi xread_key_multi_2 0-0 0-1
	message = []byte("*6\r\n$5\r\nXREAD\r\n$7\r\nstreams\r\n$15\r\nxread_key_multi\r\n$17\r\nxread_key_multi_2\r\n$3\r\n0-0\r\n$3\r\n0-1\r\n")
	expected = "*2\r\n*2\r\n$15\r\nxread_key_multi\r\n*1\r\n*2\r\n$3\r\n0-1\r\n*2\r\n$11\r\ntemperature\r\n$2\r\n95\r\n*2\r\n$17\r\nxread_key_multi_2\r\n*1\r\n*2\r\n$3\r\n0-2\r\n*2\r\n$8\r\nhumidity\r\n$2\r\n97\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// XREAD streams xread_key_multi xread_key_multi_2 0-0 0-2
	message = []byte("*6\r\n$5\r\nXREAD\r\n$7\r\nstreams\r\n$15\r\nxread_key_multi\r\n$17\r\nxread_key_multi_2\r\n$3\r\n0-0\r\n$3\r\n0-2\r\n")
	expected = "*1\r\n*2\r\n$15\r\nxread_key_multi\r\n*1\r\n*2\r\n$3\r\n0-1\r\n*2\r\n$11\r\ntemperature\r\n$2\r\n95\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

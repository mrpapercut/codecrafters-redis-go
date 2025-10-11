package commands

import "testing"

func TestHandleEcho(t *testing.T) {
	message := []byte("*2\r\n$4\r\nECHO\r\n$15\r\nmy test message\r\n")
	expected := "$15\r\nmy test message\r\n"

	response, err := HandleCommand(message)
	if err != nil {
		t.Fatalf("error handling ECHO command: %v", err)
	}

	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

package commands

import "testing"

func TestHandlePing(t *testing.T) {
	// PING
	message := []byte("*1\r\n$4\r\nPING\r\n")
	expected := "+PONG\r\n"

	response := HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

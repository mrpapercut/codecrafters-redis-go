package commands

import "testing"

func TestHandlePing(t *testing.T) {
	message := []byte("*1\r\n$4\r\nPING\r\n")
	expected := "+PONG\r\n"

	response, err := HandleCommand(message)
	if err != nil {
		t.Fatalf("error handling PING command: %v", err)
	}

	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

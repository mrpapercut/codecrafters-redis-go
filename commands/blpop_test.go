package commands

import (
	"testing"
	"time"
)

func TestHandleBLPOPWithoutWaiting(t *testing.T) {
	// RPUSH blpop_key a 0
	message := []byte("*3\r\n$5\r\nRPUSH\r\n$9\r\nblpop_key\r\n$1\r\na\r\n")
	expected := ":1\r\n"

	response := HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}

	// BLPOP blpop_key 0
	message = []byte("*3\r\n$5\r\nBLPOP\r\n$9\r\nblpop_key\r\n$1\r\n0\r\n")
	expected = "*2\r\n$9\r\nblpop_key\r\n$1\r\na\r\n"

	response = HandleCommand(message)
	if response != expected {
		t.Fatalf("expected response to be '%s', got '%s' instead", expected, response)
	}
}

func TestHandleBLPOPWithWaiting(t *testing.T) {
	// BLPOP blpop_key_wait 1
	blpopMessage := []byte("*3\r\n$5\r\nBLPOP\r\n$14\r\nblpop_key_wait\r\n$1\r\n1\r\n")
	blpopExpected := "*2\r\n$14\r\nblpop_key_wait\r\n$1\r\na\r\n"

	blpopResponse := make(chan string, 1)

	go func() {
		blpopResponse <- HandleCommand(blpopMessage)
	}()

	time.Sleep(100 * time.Millisecond)

	// LPUSH blpop_key_wait
	rpushMessage := []byte("*3\r\n$5\r\nRPUSH\r\n$14\r\nblpop_key_wait\r\n$1\r\na\r\n")
	rpushExpected := ":1\r\n"

	rpushResponse := HandleCommand(rpushMessage)
	if rpushResponse != rpushExpected {
		t.Fatalf("expected response to be '%s', got '%s' instead", rpushExpected, rpushResponse)
	}

	select {
	case response := <-blpopResponse:
		if response != blpopExpected {
			t.Fatalf("expected response to be '%s', got '%s' instead", blpopExpected, response)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("BLPOP did not unblock")
	}
}

func TestHandleBLPOPWithTimeout(t *testing.T) {
	// BLPOP blpop_key_timeout 0.5
	blpopMessage := []byte("*3\r\n$5\r\nBLPOP\r\n$17\r\nblpop_key_timeout\r\n$1\r\n2\r\n")
	blpopExpected := "*-1\r\n"

	blpopResponse := make(chan string, 1)

	go func() {
		blpopResponse <- HandleCommand(blpopMessage)
	}()

	time.Sleep(1 * time.Second)

	// LPUSH blpop_key_wait
	rpushMessage := []byte("*3\r\n$5\r\nRPUSH\r\n$14\r\nblpop_key_wait\r\n$1\r\na\r\n")
	rpushExpected := ":1\r\n"

	rpushResponse := HandleCommand(rpushMessage)
	if rpushResponse != rpushExpected {
		t.Fatalf("expected response to be '%s', got '%s' instead", rpushExpected, rpushResponse)
	}

	select {
	case response := <-blpopResponse:
		if response != blpopExpected {
			t.Fatalf("expected response to be '%s', got '%s' instead", blpopExpected, response)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("BLPOP did not unblock")
	}
}

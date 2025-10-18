package redis

import (
	"fmt"
)

const KEY_GET internalOperation = "KEY_GET"

func (r *Redis) Get(key string) (string, error) {
	responseChan := make(chan internalResponse)

	r.requestChan <- internalRequest{
		operation:    KEY_GET,
		key:          key,
		responseChan: responseChan,
	}

	response := <-responseChan

	if response.err != nil {
		return "", response.err
	}

	return response.value.ToRESP(), nil
}

func (r *Redis) handleGet(req internalRequest) {
	if r.isExpired(req.key) {
		req.responseChan <- internalResponse{err: fmt.Errorf("key not found")}
		return
	}

	value, ok := r.storage[req.key]
	if !ok {
		req.responseChan <- internalResponse{err: fmt.Errorf("key not found")}
		return
	}

	if value.Type != KeyStorage {
		req.responseChan <- internalResponse{err: fmt.Errorf("operation against a key holding the wrong kind of value")}
		return
	}

	req.responseChan <- internalResponse{
		value: value.Key,
		err:   nil,
	}
}

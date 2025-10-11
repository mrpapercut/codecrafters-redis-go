package redis

import (
	"fmt"
)

func (r *Redis) Get(key string) (string, error) {
	if r.isExpired(key) {
		return "", fmt.Errorf("key not found")
	}

	value, ok := r.storage[key]
	if !ok {
		return "", fmt.Errorf("key not found")
	}

	if value.Type != KeyStorage {
		return "", fmt.Errorf("operation against a key holding the wrong kind of value")
	}

	return value.Key.ToRESP(), nil
}

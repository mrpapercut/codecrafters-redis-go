package redis

import (
	"fmt"
	"time"
)

func (r *Redis) Get(key string) (string, error) {
	expiry, ok := r.expirations[key]
	if ok {
		if expiry.Before(time.Now()) {
			// Key expired
			r.cleanupKey(key)

			return "", fmt.Errorf("error: key expired")
		}
	}

	value, ok := r.storage[key]
	if !ok {
		return "", fmt.Errorf("error: key not found")
	}

	return value.ToRESP(), nil
}

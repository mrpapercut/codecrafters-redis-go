package redis

import (
	"fmt"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

type Redis struct {
	storage map[string]*resp.RESPValue
}

var redisLock = &sync.Mutex{}
var redisInstance *Redis

func GetInstance() *Redis {
	if redisInstance == nil {
		redisLock.Lock()
		defer redisLock.Unlock()

		if redisInstance == nil {
			redisInstance = &Redis{
				storage: make(map[string]*resp.RESPValue),
			}
		}
	}

	return redisInstance
}

func (r *Redis) Set(key *resp.RESPValue, value *resp.RESPValue) error {
	r.storage[key.String] = value

	return nil
}

func (r *Redis) Get(key string) (string, error) {
	value, ok := r.storage[key]
	if !ok {
		return "", fmt.Errorf("error: key not found")
	}

	return value.ToRESP(), nil
}

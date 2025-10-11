package redis

import (
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

type Redis struct {
	storage     map[string]*resp.RESPValue
	expirations map[string]*time.Time
}

var redisLock = &sync.Mutex{}
var redisInstance *Redis

func GetInstance() *Redis {
	if redisInstance == nil {
		redisLock.Lock()
		defer redisLock.Unlock()

		if redisInstance == nil {
			redisInstance = &Redis{
				storage:     make(map[string]*resp.RESPValue),
				expirations: make(map[string]*time.Time),
			}
		}
	}

	return redisInstance
}

func (r *Redis) cleanupKey(key string) {
	delete(r.storage, key)
	delete(r.expirations, key)
}

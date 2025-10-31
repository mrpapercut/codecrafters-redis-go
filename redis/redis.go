package redis

import (
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

type CommandOption func(*commandContextOption)

type Redis struct {
	storage       map[string]*StorageField
	expirations   map[string]*time.Time
	listWaiters   map[string][]chan *resp.RESPValue
	streamWaiters map[string]map[string][]chan *resp.RESPValue
	mu            sync.Mutex
}

var redisLock = &sync.Mutex{}
var redisInstance *Redis

func GetInstance() *Redis {
	if redisInstance == nil {
		redisLock.Lock()
		defer redisLock.Unlock()

		if redisInstance == nil {
			redisInstance = &Redis{
				storage:       make(map[string]*StorageField),
				expirations:   make(map[string]*time.Time),
				listWaiters:   make(map[string][]chan *resp.RESPValue),
				streamWaiters: make(map[string]map[string][]chan *resp.RESPValue),
				mu:            sync.Mutex{},
			}
		}
	}

	return redisInstance
}

func (r *Redis) isExpired(key string) bool {
	expiry, ok := r.expirations[key]
	if ok {
		if expiry.Before(time.Now()) {
			r.cleanupKey(key)

			return true
		}
	}

	return false
}

func (r *Redis) cleanupKey(key string) {
	delete(r.storage, key)
	delete(r.expirations, key)
}

package redis

import (
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

type StorageType string

const (
	KeyStorage  StorageType = "key"
	ListStorage StorageType = "list"
)

type StorageField struct {
	Type StorageType
	Key  *resp.RESPValue
	List []*resp.RESPValue
}

type Redis struct {
	storage     map[string]*StorageField
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
				storage:     make(map[string]*StorageField),
				expirations: make(map[string]*time.Time),
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

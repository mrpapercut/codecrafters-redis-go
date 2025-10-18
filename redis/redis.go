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

type internalOperation string

type internalRequest struct {
	operation    internalOperation
	key          string
	value        *resp.RESPValue
	opts         []CommandOption
	responseChan chan internalResponse
}

type internalResponse struct {
	value *resp.RESPValue
	err   error
}

type Redis struct {
	storage     map[string]*StorageField
	expirations map[string]*time.Time
	requestChan chan internalRequest
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
				requestChan: make(chan internalRequest),
			}

			go redisInstance.runLoop()
		}
	}

	return redisInstance
}

func (r *Redis) runLoop() {
	for request := range r.requestChan {
		switch request.operation {
		case KEY_GET:
			r.handleGet(request)
		case KEY_SET:
			r.handleSet(request)
		case LIST_GET:
			r.handleGetList(request)
		case LIST_SET:
			r.handleSetList(request)
		case LIST_REMOVE:
			r.handleRemoveList(request)
		case LIST_APPEND:
			r.handleAppendList(request)
		case LIST_PREPEND:
			r.handlePrependList(request)
		}
	}
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

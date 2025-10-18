package redis

import (
	"fmt"
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
	waiters     map[string][]chan *resp.RESPValue
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
				waiters:     make(map[string][]chan *resp.RESPValue),
				requestChan: make(chan internalRequest),
			}

			go redisInstance.runLoop()
		}
	}

	return redisInstance
}

func (r *Redis) AddWaiter(key string, ch chan *resp.RESPValue) {
	r.waiters[key] = append(r.waiters[key], ch)
}

func (r *Redis) RemoveWaiter(key string, ch chan *resp.RESPValue) {
	list := r.waiters[key]

	for i, c := range list {
		if c == ch {
			r.waiters[key] = append(list[:i], list[i+1:]...)
			break
		}
	}

	if len(r.waiters[key]) == 0 {
		delete(r.waiters, key)
	}
}

func (r *Redis) NotifyWaiters(key string) {
	list := r.waiters[key]

	if len(list) > 0 {
		value, err := r.PopList(key)
		if err != nil {
			fmt.Printf("error popping list: %v\n", err)
			return
		}

		ch := list[0]
		r.waiters[key] = list[1:]

		go func() {
			ch <- value
		}()
	}
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
		case LIST_POP:
			r.handlePopList(request)
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

package redis

import (
	"fmt"
	"time"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

const KEY_GET internalOperation = "KEY_GET"
const KEY_SET internalOperation = "KEY_SET"

func (r *Redis) Get(key string) (*resp.RESPValue, error) {
	responseChan := make(chan internalResponse)

	r.requestChan <- internalRequest{
		operation:    KEY_GET,
		key:          key,
		responseChan: responseChan,
	}

	response := <-responseChan

	return response.value, response.err
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

	if !value.IsKey() {
		req.responseChan <- internalResponse{err: fmt.Errorf("operation against a key holding the wrong kind of value")}
		return
	}

	req.responseChan <- internalResponse{
		value: value.Key,
		err:   nil,
	}
}

type commandContextOption struct {
	expire *time.Time
}

func WithExpiration(expire *time.Time) CommandOption {
	return func(o *commandContextOption) {
		o.expire = expire
	}
}

func (r *Redis) Set(key *resp.RESPValue, value *resp.RESPValue, opts ...CommandOption) error {
	responseChan := make(chan internalResponse)

	r.requestChan <- internalRequest{
		operation:    KEY_SET,
		key:          key.String,
		value:        value,
		opts:         opts,
		responseChan: responseChan,
	}

	response := <-responseChan

	return response.err
}

func (r *Redis) handleSet(req internalRequest) {
	ctx := &commandContextOption{}

	for _, opt := range req.opts {
		opt(ctx)
	}

	r.storage[req.key] = &StorageField{
		Type: KeyStorage,
		Key:  req.value,
	}

	if ctx.expire != nil {
		r.expirations[req.key] = ctx.expire
	}

	req.responseChan <- internalResponse{
		err: nil,
	}
}

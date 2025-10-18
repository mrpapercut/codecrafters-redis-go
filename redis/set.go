package redis

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

const KEY_SET internalOperation = "KEY_SET"

type CommandOption func(*commandContextOption)

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

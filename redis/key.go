package redis

import (
	"fmt"
	"time"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

func (r *Redis) Get(key string) (*resp.RESPValue, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.isExpired(key) {
		return resp.NullBulkstring(), nil
	}

	value, ok := r.storage[key]
	if !ok {
		return resp.NullBulkstring(), nil
	}

	if !value.IsKey() {
		return nil, fmt.Errorf("operation against a key holding the wrong kind of value")
	}

	return value.Key, nil
}

type commandContextOption struct {
	expire *time.Time
}

func WithExpiration(expire *time.Time) CommandOption {
	return func(o *commandContextOption) {
		o.expire = expire
	}
}

func (r *Redis) Set(key string, value *resp.RESPValue, opts ...CommandOption) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	ctx := &commandContextOption{}

	for _, opt := range opts {
		opt(ctx)
	}

	r.storage[key] = &StorageField{
		Type: KeyStorage,
		Key:  value,
	}

	if ctx.expire != nil {
		r.expirations[key] = ctx.expire
	}

	return nil
}

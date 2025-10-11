package redis

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

type CommandSetOption func(*commandSetContextOption)

type commandSetContextOption struct {
	expire *time.Time
}

func WithExpiration(expire *time.Time) CommandSetOption {
	return func(o *commandSetContextOption) {
		o.expire = expire
	}
}

func (r *Redis) Set(key *resp.RESPValue, value *resp.RESPValue, opts ...CommandSetOption) error {
	ctx := &commandSetContextOption{}

	for _, opt := range opts {
		opt(ctx)
	}

	r.storage[key.String] = value

	if ctx.expire != nil {
		r.expirations[key.String] = ctx.expire
	}

	return nil
}

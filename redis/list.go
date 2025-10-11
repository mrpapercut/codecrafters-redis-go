package redis

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

func (r *Redis) PushList(key string, val *resp.RESPValue) (int, error) {
	value, ok := r.storage[key]

	if !ok || r.isExpired(key) {
		r.storage[key] = &StorageField{
			Type: ListStorage,
			List: []*resp.RESPValue{val},
		}
	} else {
		if value.Type != ListStorage {
			return 0, fmt.Errorf("operation against a key holding the wrong kind of value")
		}

		r.storage[key].List = append(r.storage[key].List, val)
	}

	return len(r.storage[key].List), nil
}

package redis

import "github.com/codecrafters-io/redis-starter-go/resp"

func (r *Redis) Type(key string) *resp.RESPValue {
	response := &resp.RESPValue{
		Type:   resp.SimpleString,
		String: "none",
	}

	value, ok := r.storage[key]
	if !ok {
		return response
	}

	switch value.Type {
	case KeyStorage:
		response.String = "string"
	case ListStorage:
		response.String = "list"
	case StreamStorage:
		response.String = "stream"
	}

	return response
}

package redis

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

func (r *Redis) GetList(key string) (*resp.RESPValue, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	response := &resp.RESPValue{
		Type: resp.Array,
	}

	value, ok := r.storage[key]
	if !ok {
		return response, nil
	}

	if !value.IsList() {
		return nil, fmt.Errorf("operation against a key holding the wrong kind of value")
	}

	response.Array = value.List

	return response, nil
}

func (r *Redis) SetList(key string, value *resp.RESPValue) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.storage[key]
	if ok && !existing.IsList() {
		return fmt.Errorf("operation against a key holding the wrong kind of value")
	}

	r.storage[key] = &StorageField{
		Type: ListStorage,
		List: value.Array,
	}

	return nil
}

func (r *Redis) RemoveList(key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.storage[key]
	if ok {
		r.cleanupKey(key)
	}

	return nil
}

func (r *Redis) AppendList(key string, value *resp.RESPValue) (*resp.RESPValue, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	response := &resp.RESPValue{
		Type: resp.Integer,
	}

	existing, ok := r.storage[key]

	if !ok || r.isExpired(key) {
		r.storage[key] = &StorageField{
			Type: ListStorage,
			List: []*resp.RESPValue{value},
		}
	} else {
		if !existing.IsList() {
			response.Integer = 0

			return response, fmt.Errorf("operation against a key holding the wrong kind of value")
		}

		r.storage[key].List = append(r.storage[key].List, value)
	}

	go r.notifyWaiters(ListWaiter, key)

	response.Integer = int64(len(r.storage[key].List))

	return response, nil
}

func (r *Redis) PrependList(key string, value *resp.RESPValue) (*resp.RESPValue, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	response := &resp.RESPValue{
		Type: resp.Integer,
	}

	existing, ok := r.storage[key]

	if !ok || r.isExpired(key) {
		r.storage[key] = &StorageField{
			Type: ListStorage,
			List: []*resp.RESPValue{value},
		}
	} else {
		if !existing.IsList() {
			response.Integer = 0

			return response, fmt.Errorf("operation against a key holding the wrong kind of value")
		}

		oldList := r.storage[key].List

		r.storage[key].List = []*resp.RESPValue{value}
		r.storage[key].List = append(r.storage[key].List, oldList...)
	}

	go r.notifyWaiters(ListWaiter, key)

	response.Integer = int64(len(r.storage[key].List))

	return response, nil
}

func (r *Redis) PopList(key string) (*resp.RESPValue, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.storage[key]
	if !ok {
		return resp.NullBulkstring(), nil
	}

	list := &resp.RESPValue{
		Type:  resp.Array,
		Array: existing.List,
	}

	if len(list.Array) == 0 {
		return resp.NullBulkstring(), nil
	}

	popped := list.Array[0]
	list.Array = list.Array[1:]

	if len(list.Array) == 0 {
		r.cleanupKey(key)
	} else {
		r.storage[key] = &StorageField{
			Type: ListStorage,
			List: list.Array,
		}
	}

	return popped, nil
}

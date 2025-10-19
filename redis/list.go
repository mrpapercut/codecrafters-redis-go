package redis

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

const (
	LIST_GET     internalOperation = "LIST_GET"
	LIST_SET     internalOperation = "LIST_SET"
	LIST_REMOVE  internalOperation = "LIST_REMOVE"
	LIST_APPEND  internalOperation = "LIST_APPEND"
	LIST_PREPEND internalOperation = "LIST_PREPEND"
	LIST_POP     internalOperation = "LIST_POP"
)

func (r *Redis) GetList(key string) (*resp.RESPValue, error) {
	responseChan := make(chan internalResponse)

	r.requestChan <- internalRequest{
		operation:    LIST_GET,
		key:          key,
		responseChan: responseChan,
	}

	response := <-responseChan

	return response.value, response.err
}

func (r *Redis) handleGetList(req internalRequest) {
	value, ok := r.storage[req.key]
	if !ok {
		req.responseChan <- internalResponse{
			value: &resp.RESPValue{
				Type: resp.Array,
			},
			err: nil,
		}
		return
	}

	if !value.IsList() {
		req.responseChan <- internalResponse{err: fmt.Errorf("operation against a key holding the wrong kind of value")}
		return
	}

	req.responseChan <- internalResponse{
		value: &resp.RESPValue{
			Type:  resp.Array,
			Array: value.List,
		},
		err: nil,
	}
}

func (r *Redis) SetList(key string, value *resp.RESPValue) error {
	responseChan := make(chan internalResponse)

	r.requestChan <- internalRequest{
		operation:    LIST_SET,
		key:          key,
		value:        value,
		responseChan: responseChan,
	}

	response := <-responseChan

	return response.err
}

func (r *Redis) handleSetList(req internalRequest) {
	value, ok := r.storage[req.key]
	if ok && !value.IsList() {
		req.responseChan <- internalResponse{err: fmt.Errorf("operation against a key holding the wrong kind of value")}
		return
	}

	r.storage[req.key] = &StorageField{
		Type: ListStorage,
		List: req.value.Array,
	}

	req.responseChan <- internalResponse{}
}

func (r *Redis) RemoveList(key string) error {
	responseChan := make(chan internalResponse)

	r.requestChan <- internalRequest{
		operation:    LIST_REMOVE,
		key:          key,
		responseChan: responseChan,
	}

	response := <-responseChan

	return response.err
}

func (r *Redis) handleRemoveList(req internalRequest) {
	_, ok := r.storage[req.key]
	if ok {
		r.cleanupKey(req.key)
	}

	req.responseChan <- internalResponse{}
}

func (r *Redis) AppendList(key string, value *resp.RESPValue) (*resp.RESPValue, error) {
	responseChan := make(chan internalResponse)

	r.requestChan <- internalRequest{
		operation:    LIST_APPEND,
		key:          key,
		value:        value,
		responseChan: responseChan,
	}

	response := <-responseChan

	r.NotifyWaiters(key)

	return response.value, response.err
}

func (r *Redis) handleAppendList(req internalRequest) {
	response := &resp.RESPValue{
		Type: resp.Integer,
	}

	value, ok := r.storage[req.key]

	if !ok || r.isExpired(req.key) {
		r.storage[req.key] = &StorageField{
			Type: ListStorage,
			List: []*resp.RESPValue{req.value},
		}
	} else {
		if !value.IsList() {
			response.Integer = 0
			req.responseChan <- internalResponse{value: response, err: fmt.Errorf("operation against a key holding the wrong kind of value")}

			return
		}

		r.storage[req.key].List = append(r.storage[req.key].List, req.value)
	}

	response.Integer = int64(len(r.storage[req.key].List))

	req.responseChan <- internalResponse{value: response}
}

func (r *Redis) PrependList(key string, value *resp.RESPValue) (*resp.RESPValue, error) {
	responseChan := make(chan internalResponse)

	r.requestChan <- internalRequest{
		operation:    LIST_PREPEND,
		key:          key,
		value:        value,
		responseChan: responseChan,
	}

	response := <-responseChan

	r.NotifyWaiters(key)

	return response.value, response.err
}

func (r *Redis) handlePrependList(req internalRequest) {
	response := &resp.RESPValue{
		Type: resp.Integer,
	}

	value, ok := r.storage[req.key]

	if !ok || r.isExpired(req.key) {
		r.storage[req.key] = &StorageField{
			Type: ListStorage,
			List: []*resp.RESPValue{req.value},
		}
	} else {
		if !value.IsList() {
			response.Integer = 0
			req.responseChan <- internalResponse{value: response, err: fmt.Errorf("operation against a key holding the wrong kind of value")}

			return
		}

		oldList := r.storage[req.key].List

		r.storage[req.key].List = []*resp.RESPValue{req.value}
		r.storage[req.key].List = append(r.storage[req.key].List, oldList...)
	}

	response.Integer = int64(len(r.storage[req.key].List))

	req.responseChan <- internalResponse{value: response}
}

func (r *Redis) PopList(key string) (*resp.RESPValue, error) {
	responseChan := make(chan internalResponse)

	r.requestChan <- internalRequest{
		operation:    LIST_POP,
		key:          key,
		responseChan: responseChan,
	}

	response := <-responseChan

	return response.value, response.err
}

func (r *Redis) handlePopList(req internalRequest) {
	value, ok := r.storage[req.key]
	if !ok {
		req.responseChan <- internalResponse{value: resp.NullBulkstring()}
		return
	}

	list := &resp.RESPValue{
		Type:  resp.Array,
		Array: value.List,
	}

	if len(list.Array) == 0 {
		req.responseChan <- internalResponse{value: resp.NullBulkstring()}
		return
	}

	popped := list.Array[0]
	list.Array = list.Array[1:]

	if len(list.Array) == 0 {
		r.cleanupKey(req.key)
	} else {
		r.storage[req.key] = &StorageField{
			Type: ListStorage,
			List: list.Array,
		}
	}

	req.responseChan <- internalResponse{value: popped}
}

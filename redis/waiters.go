package redis

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

type WaiterType string

const (
	ListWaiter   WaiterType = "list"
	StreamWaiter WaiterType = "stream"
)

func (r *Redis) AddWaiter(waiterType WaiterType, key string, ch chan *resp.RESPValue) {
	if _, ok := r.waiters[waiterType]; !ok {
		r.waiters[waiterType] = make(map[string][]chan *resp.RESPValue)
	}

	r.waiters[waiterType][key] = append(r.waiters[waiterType][key], ch)
}

func (r *Redis) RemoveWaiter(waiterType WaiterType, key string, ch chan *resp.RESPValue) {
	list := r.waiters[waiterType][key]

	for i, c := range list {
		if c == ch {
			r.waiters[waiterType][key] = append(list[:i], list[i+1:]...)
			break
		}
	}

	if len(r.waiters[waiterType][key]) == 0 {
		delete(r.waiters[waiterType], key)
	}
}

func (r *Redis) notifyWaiters(waiterType WaiterType, key string) {
	list := r.waiters[waiterType][key]

	if len(list) > 0 {
		value, err := r.PopList(key)
		if err != nil {
			fmt.Printf("error popping list: %v\n", err)
			return
		}

		ch := list[0]
		r.waiters[waiterType][key] = list[1:]

		go func() {
			ch <- value
		}()
	}
}

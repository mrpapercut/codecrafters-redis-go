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

func (r *Redis) AddListWaiter(key string, ch chan *resp.RESPValue) {
	r.listWaiters[key] = append(r.listWaiters[key], ch)
}

func (r *Redis) RemoveListWaiter(key string, ch chan *resp.RESPValue) {
	list := r.listWaiters[key]

	for i, c := range list {
		if c == ch {
			r.listWaiters[key] = append(list[:i], list[i+1:]...)
			break
		}
	}

	if len(r.listWaiters[key]) == 0 {
		delete(r.listWaiters, key)
	}
}

func (r *Redis) notifyListWaiters(key string) {
	list := r.listWaiters[key]

	if len(list) > 0 {
		value, err := r.PopList(key)
		if err != nil {
			fmt.Printf("error popping list: %v\n", err)
			return
		}

		ch := list[0]
		r.listWaiters[key] = list[1:]

		go func() {
			ch <- value
		}()
	}
}

func (r *Redis) AddStreamWaiter(key string, id string, ch chan *resp.RESPValue) {
	if _, ok := r.streamWaiters[key]; !ok {
		r.streamWaiters[key] = make(map[string][]chan *resp.RESPValue)
	}

	if id == "$" {
		// Get max id in existing stream, or 0-0 if stream does not exist
		stream, ok := r.storage[key]
		if !ok || !r.storage[key].IsStream() {
			id = "0-0"
		} else {
			id = fmt.Sprintf("%d-%d", stream.Stream.LastEntry.Time, stream.Stream.LastEntry.Sequence)
		}
	}

	r.streamWaiters[key][id] = append(r.streamWaiters[key][id], ch)
}

func (r *Redis) RemoveStreamWaiter(key string, id string, ch chan *resp.RESPValue) {
	list := r.streamWaiters[key][id]

	for i, c := range list {
		if c == ch {
			r.streamWaiters[key][id] = append(list[:i], list[i+1:]...)
			break
		}
	}

	if len(r.streamWaiters[key][id]) == 0 {
		delete(r.streamWaiters[key], id)
	}

	if len(r.streamWaiters[key]) == 0 {
		delete(r.streamWaiters, key)
	}
}

func (r *Redis) notifyStreamWaiters(key string) {
	list := r.streamWaiters[key]

	if len(list) > 0 {
		for id, chans := range list {
			entries, err := r.GetStreamEntries(key, id)
			if err != nil {
				fmt.Printf("error getting stream entries: %v\n", err)
				return
			}

			if len(entries.Array) > 1 {
				ch := chans[0]
				r.streamWaiters[key][id] = chans[1:]

				go func() {
					ch <- entries
				}()

				break
			}
		}
	}
}

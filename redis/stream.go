package redis

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

/**
 * A stream is handled as a map[key][]StorageField{Stream: map[entry id ms]map[entry id seq][]*resp.RESPValue}
 * where sequence ids are the index of the final slice.
 */

const STREAM_APPEND internalOperation = "STREAM_APPEND"

func (r *Redis) AppendStream(key string, id string, value *resp.RESPValue) (*resp.RESPValue, error) {
	responseChan := make(chan internalResponse)

	r.requestChan <- internalRequest{
		operation:    STREAM_APPEND,
		key:          key,
		id:           id,
		value:        value,
		responseChan: responseChan,
	}

	response := <-responseChan

	return response.value, response.err
}

func (r *Redis) handleAppendStream(req internalRequest) {
	stream, ok := r.storage[req.key]
	if !ok {
		stream = &StorageField{
			Type:   StreamStorage,
			Stream: &StreamField{Entries: make(map[int64]map[int64][]*resp.RESPValue)},
		}

		r.storage[req.key] = stream
	}

	if !r.storage[req.key].IsStream() {
		req.responseChan <- internalResponse{err: fmt.Errorf("operation against a key holding the wrong kind of value")}
		return
	}

	entryIDms, entryIDseq, err := r.parseStreamEntryID(req.id)
	if err != nil {
		req.responseChan <- internalResponse{err: err}
		return
	}

	if entryIDms == 0 && entryIDseq == 0 {
		req.responseChan <- internalResponse{err: fmt.Errorf("The ID specified in XADD must be greater than 0-0")}
		return
	}

	// Check if ms part of id is latest
	if entryIDms >= 0 && r.storage[req.key].Stream.LastEntry != nil && entryIDms < r.storage[req.key].Stream.LastEntry.Time {
		req.responseChan <- internalResponse{err: fmt.Errorf("The ID specified in XADD is equal or smaller than the target stream top item")}
		return
	}

	// Check if seq part of id is latest
	if entryIDseq >= 0 {
		_, ok = r.storage[req.key].Stream.Entries[entryIDms]
		if ok && entryIDseq <= r.storage[req.key].Stream.LastEntry.Sequence {
			req.responseChan <- internalResponse{err: fmt.Errorf("The ID specified in XADD is equal or smaller than the target stream top item")}
			return
		}

		if !ok {
			r.storage[req.key].Stream.Entries[entryIDms] = make(map[int64][]*resp.RESPValue)
		}
	}

	// All good, create the stream in storage
	r.storage[req.key].Stream.Entries[entryIDms][entryIDseq] = make([]*resp.RESPValue, 0)
	r.storage[req.key].Stream.Entries[entryIDms][entryIDseq] = append(r.storage[req.key].Stream.Entries[entryIDms][entryIDseq], &resp.RESPValue{
		Type: resp.Map,
		Map:  req.value.Map,
	})

	// Store latest entry in stream for easier lookup
	r.storage[req.key].Stream.LastEntry = &StreamEntryID{
		Time:     entryIDms,
		Sequence: entryIDseq,
	}

	response := &resp.RESPValue{
		Type:   resp.BulkString,
		String: fmt.Sprintf("%d-%d", entryIDms, entryIDseq),
	}

	req.responseChan <- internalResponse{value: response}
}

func (r *Redis) parseStreamEntryID(id string) (int64, int64, error) {
	if id == "*" {
		return time.Now().UnixMilli(), 0, nil
	}

	parts := strings.Split(id, "-")
	if len(parts) != 2 {
		return -1, -1, fmt.Errorf("Invalid stream ID specified as stream command argument")
	}

	ms, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return -1, -1, fmt.Errorf("Invalid stream ID specified as stream command argument")
	}

	if parts[1] == "*" {
		return ms, -1, nil
	}

	seq, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return -1, -1, fmt.Errorf("Invalid stream ID specified as stream command argument")
	}

	return ms, seq, nil
}

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

	validMS, validSeq, err := r.validateStreamEntryID(req.key, entryIDms, entryIDseq)
	if err != nil {
		req.responseChan <- internalResponse{err: err}
		return
	}

	// All good, create the stream in storage
	_, ok = r.storage[req.key].Stream.Entries[validMS]
	if !ok {
		r.storage[req.key].Stream.Entries[validMS] = make(map[int64][]*resp.RESPValue)
	}

	r.storage[req.key].Stream.Entries[validMS][validSeq] = make([]*resp.RESPValue, 0)
	r.storage[req.key].Stream.Entries[validMS][validSeq] = append(r.storage[req.key].Stream.Entries[validMS][validSeq], &resp.RESPValue{
		Type: resp.Map,
		Map:  req.value.Map,
	})

	// Store latest entry in stream for easier lookup
	r.storage[req.key].Stream.LastEntry = &StreamEntryID{
		Time:     validMS,
		Sequence: validSeq,
	}

	response := &resp.RESPValue{
		Type:   resp.BulkString,
		String: fmt.Sprintf("%d-%d", validMS, validSeq),
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

func (r *Redis) validateStreamEntryID(key string, entryIDms int64, entryIDseq int64) (int64, int64, error) {
	if entryIDms == 0 && entryIDseq == 0 {
		return -1, -1, fmt.Errorf("The ID specified in XADD must be greater than 0-0")
	}

	if r.storage[key].Stream.LastEntry != nil && entryIDms < r.storage[key].Stream.LastEntry.Time {
		return -1, -1, fmt.Errorf("The ID specified in XADD is equal or smaller than the target stream top item")
	}

	if entryIDseq < 0 {
		if r.storage[key].Stream.LastEntry != nil && r.storage[key].Stream.LastEntry.Time == entryIDms {
			entryIDseq = r.storage[key].Stream.LastEntry.Sequence + 1
		} else if entryIDms == 0 {
			entryIDseq = 1
		} else {
			entryIDseq = 0
		}
	}

	if r.storage[key].Stream.LastEntry != nil && entryIDseq <= r.storage[key].Stream.LastEntry.Sequence {
		return -1, -1, fmt.Errorf("The ID specified in XADD is equal or smaller than the target stream top item")
	}

	return entryIDms, entryIDseq, nil
}

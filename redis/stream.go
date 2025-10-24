package redis

//lint:file-ignore ST1005 errors are in Redis format

import (
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

/**
 * A stream is handled as a map[key][]StorageField{Stream: map[entry id ms]map[entry id seq][]*resp.RESPValue}
 * where sequence ids are the index of the final slice.
 */

const STREAM_RANGE_GET internalOperation = "STREAM_RANGE_GET"
const STREAM_APPEND internalOperation = "STREAM_APPEND"

type XRangeStartEnd struct {
	StartMS  int64
	StartSeq int64
	EndMS    int64
	EndSeq   int64
}

func (r *Redis) GetStreamEntries(key string, id string) (*resp.RESPValue, error) {
	stream, ok := r.storage[key]
	if !ok {
		// return empty array or something
		return nil, fmt.Errorf("stream not found")
	}

	if !r.storage[key].IsStream() {
		// req.responseChan <- internalResponse{err: fmt.Errorf("operation against a key holding the wrong kind of value")}
		return nil, fmt.Errorf("operation against a key holding the wrong kind of value")
	}

	idMS, idSeq, err := r.parseStreamEntryID(id)
	if err != nil {
		return nil, fmt.Errorf("error: invalid stream entry ID: %s", id)
	}

	slog.Info("GetStream", "stream", stream, "idMS", idMS, "idSeq", idSeq)

	response := &resp.RESPValue{
		Type:  resp.Array,
		Array: make([]*resp.RESPValue, 0),
	}

	response.Array = append(response.Array, &resp.RESPValue{
		Type:   resp.BulkString,
		String: key,
	})

	for idx, sequences := range stream.Stream.Entries {
		if idx < idMS {
			continue
		}

		streamEntries := &resp.RESPValue{
			Type:  resp.Array,
			Array: make([]*resp.RESPValue, 0),
		}

		for seq, entries := range sequences {
			if idx == idMS && seq <= idSeq {
				continue
			}

			if len(entries) == 0 {
				continue
			}

			sequence := &resp.RESPValue{
				Type:  resp.Array,
				Array: make([]*resp.RESPValue, 0),
			}

			slog.Info("GetStream", "appending", fmt.Sprintf("%d-%d", idx, seq))

			sequence.Array = append(sequence.Array, &resp.RESPValue{
				Type:   resp.BulkString,
				String: fmt.Sprintf("%d-%d", idx, seq),
			})

			for _, entry := range entries {
				sequence.Array = append(sequence.Array, r.getStreamValuesAsSlice(entry))
			}

			streamEntries.Array = append(streamEntries.Array, sequence)
		}

		response.Array = append(response.Array, streamEntries)
	}

	return response, nil
}

func (r *Redis) GetStreamsByRange(key string, xrangeRange *XRangeStartEnd) (*resp.RESPValue, error) {
	// TODO fix async
	stream, ok := r.storage[key]
	if !ok {
		// return empty array or something
		return nil, fmt.Errorf("stream not found")
	}

	if !r.storage[key].IsStream() {
		// req.responseChan <- internalResponse{err: fmt.Errorf("operation against a key holding the wrong kind of value")}
		return nil, fmt.Errorf("operation against a key holding the wrong kind of value")
	}

	response := &resp.RESPValue{
		Type:  resp.Array,
		Array: make([]*resp.RESPValue, 0),
	}

	for idx, sequences := range stream.Stream.Entries {
		if idx < xrangeRange.StartMS || idx > xrangeRange.EndMS {
			continue
		}

		for seq, entries := range sequences {
			if idx == xrangeRange.StartMS && (seq < xrangeRange.StartSeq || seq > xrangeRange.EndSeq) {
				continue
			}

			if len(entries) == 0 {
				continue
			}

			sequence := &resp.RESPValue{
				Type:  resp.Array,
				Array: make([]*resp.RESPValue, 0),
			}

			sequence.Array = append(sequence.Array, &resp.RESPValue{
				Type:   resp.BulkString,
				String: fmt.Sprintf("%d-%d", idx, seq),
			})

			for _, entry := range entries {
				sequence.Array = append(sequence.Array, r.getStreamValuesAsSlice(entry))
			}

			response.Array = append(response.Array, sequence)
		}
	}

	r.sortStreams(response)

	return response, nil
}

func (r *Redis) getStreamValuesAsSlice(entry *resp.RESPValue) *resp.RESPValue {
	arr := &resp.RESPValue{
		Type:  resp.Array,
		Array: make([]*resp.RESPValue, 0),
	}

	for k, v := range entry.Map {
		arr.Array = append(arr.Array, k, v)
	}

	return arr
}

func (r *Redis) sortStreams(streams *resp.RESPValue) {
	slices.SortFunc(streams.Array, func(a *resp.RESPValue, b *resp.RESPValue) int {
		partsA := strings.Split(a.Array[0].String, "-")
		partsB := strings.Split(b.Array[0].String, "-")

		msA, _ := strconv.Atoi(partsA[0])
		msB, _ := strconv.Atoi(partsB[0])
		seqA, _ := strconv.Atoi(partsA[1])
		seqB, _ := strconv.Atoi(partsB[1])

		if msA == msB {
			return seqA - seqB
		}

		return msA - msB
	})
}

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
	_, ok := r.storage[req.key]
	if !ok {
		stream := &StorageField{
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

	if r.storage[key].Stream.LastEntry != nil && r.storage[key].Stream.LastEntry.Time == entryIDms && entryIDseq <= r.storage[key].Stream.LastEntry.Sequence {
		return -1, -1, fmt.Errorf("The ID specified in XADD is equal or smaller than the target stream top item")
	}

	return entryIDms, entryIDseq, nil
}

package redis

//lint:file-ignore ST1005 errors are in Redis format

import (
	"fmt"
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

type XRangeStartEnd struct {
	StartMS  int64
	StartSeq int64
	EndMS    int64
	EndSeq   int64
}

func (r *Redis) GetStreamEntries(key string, id string) (*resp.RESPValue, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	stream, ok := r.storage[key]
	if !ok {
		return nil, fmt.Errorf("stream not found")
	}

	if !r.storage[key].IsStream() {
		return nil, fmt.Errorf("operation against a key holding the wrong kind of value")
	}

	idMS, idSeq, err := r.parseStreamEntryID(id)
	if err != nil {
		return nil, fmt.Errorf("error: invalid stream entry ID: %s", id)
	}

	response := &resp.RESPValue{
		Type: resp.Array,
		Array: []*resp.RESPValue{{
			Type:   resp.BulkString,
			String: key,
		}},
	}

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

			sequence := r.getStreamSequence(idx, seq, entries)

			streamEntries.Array = append(streamEntries.Array, sequence)
		}

		response.Array = append(response.Array, streamEntries)
	}

	return response, nil
}

func (r *Redis) GetStreamsByRange(key string, xrangeRange *XRangeStartEnd) (*resp.RESPValue, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	stream, ok := r.storage[key]
	if !ok {
		return nil, fmt.Errorf("stream not found")
	}

	if !r.storage[key].IsStream() {
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

			sequence := r.getStreamSequence(idx, seq, entries)

			response.Array = append(response.Array, sequence)
		}
	}

	r.sortStreams(response)

	return response, nil
}

func (r *Redis) getStreamSequence(idx int64, seq int64, entries []*resp.RESPValue) *resp.RESPValue {
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

	return sequence
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
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.storage[key]
	if !ok {
		stream := &StorageField{
			Type:   StreamStorage,
			Stream: &StreamField{Entries: make(map[int64]map[int64][]*resp.RESPValue)},
		}

		r.storage[key] = stream
	}

	if !r.storage[key].IsStream() {
		return nil, fmt.Errorf("operation against a key holding the wrong kind of value")
	}

	entryIDms, entryIDseq, err := r.parseStreamEntryID(id)
	if err != nil {
		return nil, err
	}

	validMS, validSeq, err := r.validateStreamEntryID(key, entryIDms, entryIDseq)
	if err != nil {
		return nil, err
	}

	// All good, create the stream in storage
	_, ok = r.storage[key].Stream.Entries[validMS]
	if !ok {
		r.storage[key].Stream.Entries[validMS] = make(map[int64][]*resp.RESPValue)
	}

	r.storage[key].Stream.Entries[validMS][validSeq] = make([]*resp.RESPValue, 0)
	r.storage[key].Stream.Entries[validMS][validSeq] = append(r.storage[key].Stream.Entries[validMS][validSeq], &resp.RESPValue{
		Type: resp.Map,
		Map:  value.Map,
	})

	// Store latest entry in stream for easier lookup
	r.storage[key].Stream.LastEntry = &StreamEntryID{
		Time:     validMS,
		Sequence: validSeq,
	}

	response := &resp.RESPValue{
		Type:   resp.BulkString,
		String: fmt.Sprintf("%d-%d", validMS, validSeq),
	}

	return response, nil
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

package commands

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

const XREAD SupportedCommand = "xread"

func HandleXREAD(cmd *resp.RESPValue) string {
	if len(cmd.Array) < 4 {
		return resp.GenericError("wrong number of arguments for 'xread' command")
	}

	hasBlock := false
	var blockValue int64 = -1

	streams := make([]string, 0)
	ids := make([]string, 0)

	argIdx := 2
	if cmd.Array[1].String == "block" {
		hasBlock = true

		val, err := strconv.ParseInt(cmd.Array[2].String, 10, 64)
		if err != nil {
			return resp.GenericError("timeout is not an integer or out of range")
		}

		blockValue = val

		argIdx = 3
	}

	// fmt.Printf("hasBlock? %t (block value: %d)\n", hasBlock, blockValue)

	// fmt.Printf("rest args: %#v (length: %d)\n", cmd.Array[argIdx:], len(cmd.Array[argIdx:]))

	if len(cmd.Array[argIdx:])%2 != 0 {
		return resp.GenericError("Unbalanced 'xread' list of streams: for each stream key an ID, '+', or '$' must be specified.")
	}

	streamCount := (len(cmd.Array[argIdx:]) / 2) + argIdx
	for ; argIdx < streamCount; argIdx++ {
		streams = append(streams, cmd.Array[argIdx].String)
	}

	idCount := len(cmd.Array[argIdx:]) + argIdx
	for ; argIdx < idCount; argIdx++ {
		ids = append(ids, cmd.Array[argIdx].String)
	}

	slog.Info("HandleXREAD", "hasBlock?", hasBlock, "blockValue", blockValue, "streams", streams, "ids", ids)

	allResponses := &resp.RESPValue{
		Type:  resp.Array,
		Array: make([]*resp.RESPValue, 0),
	}

	for i, key := range streams {
		streamResponse, err := redisInstance.GetStreamEntries(key, ids[i])
		if err != nil {
			return resp.GenericError(fmt.Sprintf("error getting stream: %v", err))
		}
		allResponses.Array = append(allResponses.Array, streamResponse)
	}

	// fmt.Printf("streams: %s\n", strings.Join(streams, ", "))
	// fmt.Printf("ids: %s", strings.Join(ids, ", "))

	return allResponses.ToRESP()
}

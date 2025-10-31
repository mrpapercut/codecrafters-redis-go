package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

const XREAD SupportedCommand = "xread"

func HandleXREAD(cmd *resp.RESPValue) string {
	if len(cmd.Array) < 4 {
		return resp.GenericError("wrong number of arguments for 'xread' command")
	}

	hasBlock := false
	var timeoutArg float64 = -1

	streams := make([]string, 0)
	ids := make([]string, 0)

	argIdx := 2
	if strings.ToLower(cmd.Array[1].String) == "block" {
		hasBlock = true

		timeoutArg, err := strconv.ParseFloat(cmd.Array[2].String, 64)
		if err != nil {
			return resp.GenericError("timeout is not an integer or out of range")
		}
		if timeoutArg == 0 {
			timeoutArg = 300
		}

		argIdx = 4
	}

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

	if hasBlock {
		return handleXREADBlock(timeoutArg, streams, ids)
	}

	allResponses := &resp.RESPValue{
		Type:  resp.Array,
		Array: make([]*resp.RESPValue, 0),
	}

	for i, key := range streams {
		streamResponse, err := redisInstance.GetStreamEntries(key, ids[i])
		if err != nil {
			return resp.GenericError(fmt.Sprintf("error getting stream: %v", err))
		}

		if len(streamResponse.Array[1].Array) > 0 {
			allResponses.Array = append(allResponses.Array, streamResponse)
		}
	}

	return allResponses.ToRESP()
}

func handleXREADBlock(timeoutArg float64, streams []string, ids []string) string {
	timeout := time.Duration(timeoutArg * float64(time.Second))

	ch := make(chan *resp.RESPValue)

	for i, stream := range streams {
		redisInstance.AddStreamWaiter(stream, ids[i], ch)
	}

	select {
	case res := <-ch:
		for i, stream := range streams {
			redisInstance.RemoveStreamWaiter(stream, ids[i], ch)
		}

		response := &resp.RESPValue{
			Type:  resp.Array,
			Array: []*resp.RESPValue{res},
		}

		return response.ToRESP()
	case <-time.After(timeout):
		nullArray := resp.NullArray()

		return nullArray.ToRESP()
	}
}

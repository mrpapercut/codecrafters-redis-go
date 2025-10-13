package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

const RPUSH SupportedCommand = "rpush"

func HandleRPUSH(cmd *resp.RESPValue) string {
	lastLength := 0

	for _, val := range cmd.Array[2:] {
		listLength, err := redisInstance.PushList(cmd.Array[1].String, val)
		if err != nil {
			return resp.GenericError(fmt.Sprintf("error handling rpush: %v", err))
		}

		lastLength = listLength
	}

	resp := &resp.RESPValue{
		Type:    resp.Integer,
		Integer: int64(lastLength),
	}

	return resp.ToRESP()
}

package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

const LPUSH SupportedCommand = "lpush"

func HandleLPUSH(cmd *resp.RESPValue) string {
	var lastLength int64 = 0

	for _, val := range cmd.Array[2:] {
		listLength, err := redisInstance.PrependList(cmd.Array[1].String, val)
		if err != nil {
			return resp.GenericError(fmt.Sprintf("error handling rpush: %v", err))
		}

		lastLength = listLength.Integer
	}

	resp := &resp.RESPValue{
		Type:    resp.Integer,
		Integer: lastLength,
	}

	return resp.ToRESP()
}

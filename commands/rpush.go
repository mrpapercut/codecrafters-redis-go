package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

const RPUSH SupportedCommand = "rpush"

func HandleRPUSH(cmd *resp.RESPValue) (string, error) {
	lastLength := 0

	for i := range cmd.Array[2:] {
		listLength, err := redisInstance.PushList(cmd.Array[1].String, cmd.Array[i])
		if err != nil {
			return "", fmt.Errorf("error handling rpush: %v", err)
		}

		lastLength = listLength
	}

	resp := &resp.RESPValue{
		Type:    resp.Integer,
		Integer: int64(lastLength),
	}

	return resp.ToRESP(), nil
}

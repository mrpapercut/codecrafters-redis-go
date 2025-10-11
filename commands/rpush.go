package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

const RPUSH SupportedCommand = "rpush"

func HandleRPUSH(cmd *resp.RESPValue) (string, error) {
	listLength, err := redisInstance.PushList(cmd.Array[1].String, cmd.Array[2])
	if err != nil {
		return "", fmt.Errorf("error handling rpush: %v", err)
	}

	resp := &resp.RESPValue{
		Type:    resp.Integer,
		Integer: int64(listLength),
	}

	return resp.ToRESP(), nil
}

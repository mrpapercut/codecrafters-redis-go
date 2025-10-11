package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

const SET SupportedCommand = "set"

func HandleSET(cmd *resp.RESPValue) (string, error) {
	err := redisInstance.Set(cmd.Array[1], cmd.Array[2])
	if err != nil {
		return "", fmt.Errorf("error handling SET")
	}

	return RESPONSE_OK, nil
}

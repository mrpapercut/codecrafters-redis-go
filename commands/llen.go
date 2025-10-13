package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

const LLEN SupportedCommand = "llen"

func HandleLLEN(cmd *resp.RESPValue) string {
	list, err := redisInstance.GetList(cmd.Array[1].String)
	if err != nil {
		return resp.GenericError(fmt.Sprintf("error handling LLEN: %v", err))
	}

	response := &resp.RESPValue{
		Type:    resp.Integer,
		Integer: int64(len(list.Array)),
	}

	return response.ToRESP()
}

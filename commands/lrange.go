package commands

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

const LRANGE SupportedCommand = "lrange"

func HandleLRANGE(cmd *resp.RESPValue) string {
	listKey := cmd.Array[1].String
	start, err := strconv.Atoi(cmd.Array[2].String)
	if err != nil {
		return resp.SyntaxError("invalid start offset")
	}

	stop, err := strconv.Atoi(cmd.Array[3].String)
	if err != nil {
		return resp.SyntaxError("invalid stop offset")
	}

	list, err := redisInstance.GetList(listKey)
	if err != nil {
		return resp.GenericError(fmt.Sprintf("error retrieving list '%s'", listKey))
	}

	if stop >= len(list.Array) {
		stop = len(list.Array) - 1
	}

	response := &resp.RESPValue{
		Type: resp.Array,
	}

	if len(list.Array) == 0 || start >= len(list.Array) || start > stop {
		return response.ToRESP()
	}

	for i := start; i <= stop; i++ {
		response.Array = append(response.Array, list.Array[i])
	}

	return response.ToRESP()
}

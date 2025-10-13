package commands

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

const LPOP SupportedCommand = "lpop"

func HandleLPOP(cmd *resp.RESPValue) string {
	if len(cmd.Array) == 3 {
		return handleLPOPMultiple(cmd)
	}

	key := cmd.Array[1].String

	list, err := redisInstance.GetList(key)
	if err != nil {
		return resp.WrongTypeError(err.Error())
	}

	if len(list.Array) == 0 {
		response := &resp.RESPValue{
			Type:   resp.BulkString,
			IsNull: true,
		}

		return response.ToRESP()
	}

	popped := list.Array[0]

	list.Array = list.Array[1:]
	if len(list.Array) == 0 {
		redisInstance.RemoveList(key)
	} else {
		redisInstance.SetList(key, list.Array)
	}

	return popped.ToRESP()
}

func handleLPOPMultiple(cmd *resp.RESPValue) string {
	key := cmd.Array[1].String
	count, err := strconv.Atoi(cmd.Array[2].String)
	if err != nil {
		return resp.SyntaxError("LPOP argument not an integer")
	}

	list, err := redisInstance.GetList(key)
	if err != nil {
		return resp.WrongTypeError(err.Error())
	}

	if len(list.Array) == 0 {
		response := &resp.RESPValue{
			Type:   resp.Array,
			IsNull: true,
		}

		return response.ToRESP()
	}

	popped := make([]*resp.RESPValue, 0)
	for i := range count {
		popped = append(popped, list.Array[i])
	}

	list.Array = list.Array[count:]
	if len(list.Array) == 0 {
		redisInstance.RemoveList(key)
	} else {
		redisInstance.SetList(key, list.Array)
	}

	response := &resp.RESPValue{
		Type:  resp.Array,
		Array: popped,
	}

	return response.ToRESP()
}

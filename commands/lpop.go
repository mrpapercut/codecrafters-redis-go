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

	popped, err := redisInstance.PopList(key)
	if err != nil {
		return resp.GenericError(err.Error())
	}

	return popped.ToRESP()
}

func handleLPOPMultiple(cmd *resp.RESPValue) string {
	key := cmd.Array[1].String

	count, err := strconv.Atoi(cmd.Array[2].String)
	if err != nil {
		return resp.SyntaxError("LPOP argument not an integer")
	}

	_, err = redisInstance.GetList(key)
	if err != nil {
		return resp.WrongTypeError(err.Error())
	}

	popped := make([]*resp.RESPValue, 0)
	for range count {
		p, _ := redisInstance.PopList(key)
		popped = append(popped, p)
	}

	response := &resp.RESPValue{
		Type:  resp.Array,
		Array: popped,
	}

	return response.ToRESP()
}

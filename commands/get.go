package commands

import (
	"github.com/codecrafters-io/redis-starter-go/resp"
)

const GET SupportedCommand = "get"

func HandleGET(cmd *resp.RESPValue) (string, error) {
	res, err := redisInstance.Get(cmd.Array[1].String)
	if err != nil {
		nullObj := &resp.RESPValue{
			Type:   resp.BulkString,
			IsNull: true,
		}
		return nullObj.ToRESP(), nil
	}

	return res, nil
}

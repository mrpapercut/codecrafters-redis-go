package commands

import (
	"github.com/codecrafters-io/redis-starter-go/resp"
)

const GET SupportedCommand = "get"

func HandleGET(cmd *resp.RESPValue) string {
	res, err := redisInstance.Get(cmd.Array[1].String)
	if err != nil {
		return resp.NullBulkstring().ToRESP()
	}

	return res.ToRESP()
}

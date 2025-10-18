package commands

import "github.com/codecrafters-io/redis-starter-go/resp"

const TYPE SupportedCommand = "type"

func HandleTYPE(cmd *resp.RESPValue) string {
	key := cmd.Array[1].String

	response := redisInstance.Type(key)

	return response.ToRESP()
}

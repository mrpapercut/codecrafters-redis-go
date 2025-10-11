package commands

import "github.com/codecrafters-io/redis-starter-go/resp"

const PING SupportedCommand = "ping"

func HandlePING(cmd *resp.RESPValue) (string, error) {
	return "+PONG\r\n", nil
}

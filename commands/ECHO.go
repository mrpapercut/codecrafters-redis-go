package commands

import (
	"github.com/codecrafters-io/redis-starter-go/resp"
)

const ECHO SupportedCommand = "echo"

func HandleECHO(cmd *resp.RESPValue) string {
	return cmd.Array[1].ToRESP()
}

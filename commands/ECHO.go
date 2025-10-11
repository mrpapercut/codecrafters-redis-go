package commands

import (
	"github.com/codecrafters-io/redis-starter-go/resp"
)

const ECHO SupportedCommand = "ECHO"

func HandleECHO(cmd *resp.RESPValue) (string, error) {
	return cmd.Array[1].ToRESP(), nil
}

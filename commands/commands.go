package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

var parser = resp.GetParser()

func HandleCommand(rawcmd []byte) (string, error) {
	parsed := parser.ParseCommand(string(rawcmd))

	switch parsed.Type {
	case resp.Array:
		if len(parsed.Array) == 0 {
			return "", fmt.Errorf("missing command '%s'", rawcmd)
		}

		cmd := parsed.Array[0].String

		switch cmd {
		case "PING":
			return "+PONG\r\n", nil
		default:
			return "", fmt.Errorf("unsupported command '%s'", rawcmd)
		}
	default:
		return "", fmt.Errorf("unsupported command '%s'", rawcmd)
	}
}

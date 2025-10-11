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
		cmd := parsed.Array[0].String

		if cmd == "" {
			return "", fmt.Errorf("error: no command provided")
		}
	case resp.SimpleString:
		cmd := parsed.String

		switch cmd {
		case "PING":
			return "+PONG\r\n", nil
		}
	}

	return "", fmt.Errorf("unsupported command '%s'", rawcmd)
}

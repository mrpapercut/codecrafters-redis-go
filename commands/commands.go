package commands

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

type SupportedCommand string

var parser = resp.GetParser()

func HandleCommand(rawcmd []byte) (string, error) {
	parsed := parser.ParseCommand(string(rawcmd))

	switch parsed.Type {
	case resp.Array:
		if len(parsed.Array) == 0 {
			return "", fmt.Errorf("missing command '%s'", rawcmd)
		}

		cmd := strings.ToLower(parsed.Array[0].String)

		switch SupportedCommand(cmd) {
		case PING:
			return HandlePING(parsed)
		case ECHO:
			return HandleECHO(parsed)
		default:
			return "", fmt.Errorf("unsupported command '%s'", rawcmd)
		}
	default:
		return "", fmt.Errorf("unsupported command '%s'", rawcmd)
	}
}

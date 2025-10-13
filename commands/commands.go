package commands

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/redis"
	"github.com/codecrafters-io/redis-starter-go/resp"
)

type SupportedCommand string

const RESPONSE_OK = "+OK\r\n"

var parser = resp.GetParser()
var redisInstance = redis.GetInstance()

func HandleCommand(rawcmd []byte) string {
	parsed := parser.ParseCommand(string(rawcmd))

	switch parsed.Type {
	case resp.Array:
		if len(parsed.Array) == 0 {
			return resp.SyntaxError(fmt.Sprintf("missing command '%s'", rawcmd))
		}

		cmd := strings.ToLower(parsed.Array[0].String)

		supportedCommands := map[SupportedCommand]func(*resp.RESPValue) string{
			PING:   HandlePING,
			ECHO:   HandleECHO,
			SET:    HandleSET,
			GET:    HandleGET,
			RPUSH:  HandleRPUSH,
			LRANGE: HandleLRANGE,
			LLEN:   HandleLLEN,
		}

		fn, ok := supportedCommands[SupportedCommand(cmd)]
		if !ok {
			return resp.SyntaxError(fmt.Sprintf("unsupported command '%s'", cmd))
		}

		return fn(parsed)
	default:
		return resp.SyntaxError(fmt.Sprintf("invalid syntax '%s'", rawcmd))
	}
}

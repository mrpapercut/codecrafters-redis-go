package commands

import (
	"github.com/codecrafters-io/redis-starter-go/resp"
)

const XADD SupportedCommand = "xadd"

func HandleXADD(cmd *resp.RESPValue) string {
	if len(cmd.Array) < 5 && len(cmd.Array)%2 == 0 {
		return resp.GenericError("wrong number of arguments for 'xadd' command")
	}

	key := cmd.Array[1].String
	id := cmd.Array[2].String

	keyValPairs := &resp.RESPValue{
		Type: resp.Map,
		Map:  make(map[*resp.RESPValue]*resp.RESPValue),
	}

	for i := 3; i+1 < len(cmd.Array); i += 2 {
		k := &resp.RESPValue{
			Type:   resp.BulkString,
			String: cmd.Array[i].String,
		}

		v := &resp.RESPValue{
			Type:   resp.BulkString,
			String: cmd.Array[i+1].String,
		}

		keyValPairs.Map[k] = v
	}

	storedID, err := redisInstance.AppendStream(key, id, keyValPairs)
	if err != nil {
		return resp.GenericError(err.Error())
	}

	return storedID.ToRESP()
}

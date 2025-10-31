package commands

import (
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

const BLPOP SupportedCommand = "blpop"

func HandleBLPOP(cmd *resp.RESPValue) string {
	key := cmd.Array[1].String

	timeoutArg, err := strconv.ParseFloat(cmd.Array[2].String, 64)
	if err != nil {
		return resp.SyntaxError("BLPOP argument not an integer or double")
	}
	if timeoutArg == 0 {
		timeoutArg = 300
	}

	list, err := redisInstance.GetList(key)
	if err != nil {
		return resp.WrongTypeError(err.Error())
	}

	if len(list.Array) > 0 {
		popped, err := redisInstance.PopList(key)
		if err != nil {
			return resp.GenericError(err.Error())
		}

		response := &resp.RESPValue{
			Type: resp.Array,
			Array: []*resp.RESPValue{{
				Type:   resp.BulkString,
				String: key,
			}, popped},
		}

		return response.ToRESP()
	}

	timeout := time.Duration(timeoutArg * float64(time.Second))

	ch := make(chan *resp.RESPValue)

	redisInstance.AddListWaiter(key, ch)

	select {
	case res := <-ch:
		redisInstance.RemoveListWaiter(key, ch)

		response := &resp.RESPValue{
			Type: resp.Array,
			Array: []*resp.RESPValue{{
				Type:   resp.BulkString,
				String: key,
			}, res},
		}

		return response.ToRESP()
	case <-time.After(timeout):
		nullArray := resp.NullArray()

		return nullArray.ToRESP()
	}
}

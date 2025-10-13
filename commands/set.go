package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/redis"
	"github.com/codecrafters-io/redis-starter-go/resp"
)

const SET SupportedCommand = "set"

func HandleSET(cmd *resp.RESPValue) string {
	opts := make([]redis.CommandSetOption, 0)

	for i := 0; i < len(cmd.Array); i++ {
		option := strings.ToLower(cmd.Array[i].String)

		var value int
		var err error

		if option == "ex" || option == "px" {
			if len(cmd.Array) < i+1 {
				return resp.SyntaxError(fmt.Sprintf("error: missing argument for '%s'", option))
			}

			value, err = strconv.Atoi(cmd.Array[i+1].String)
			if err != nil {
				return resp.SyntaxError(fmt.Sprintf("error: invalid argument '%s'", cmd.Array[i+1].String))
			}

			var expiry time.Time
			if option == "ex" {
				expiry = time.Now().Add(time.Duration(value) * time.Second)
			}

			if option == "px" {
				expiry = time.Now().Add(time.Duration(value) * time.Millisecond)
			}

			opts = append(opts, redis.WithExpiration(&expiry))
		}
	}

	err := redisInstance.Set(cmd.Array[1], cmd.Array[2], opts...)
	if err != nil {
		return resp.GenericError("error handling SET")
	}

	response := &resp.RESPValue{
		Type:   resp.SimpleString,
		String: "OK",
	}

	return response.ToRESP()
}

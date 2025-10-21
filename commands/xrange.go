package commands

//lint:file-ignore ST1005 errors are in Redis format

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/redis"
	"github.com/codecrafters-io/redis-starter-go/resp"
)

const XRANGE SupportedCommand = "xrange"

func HandleXRANGE(cmd *resp.RESPValue) string {
	if len(cmd.Array) != 4 {
		return resp.GenericError("wrong number of arguments for 'xrange' command")
	}

	key := cmd.Array[1].String

	xrangeRange, err := parseXRANGEStartEnd(cmd.Array[2].String, cmd.Array[3].String)
	if err != nil {
		return resp.GenericError(err.Error())
	}

	streams, err := redisInstance.GetStreamsByRange(key, xrangeRange)
	if err != nil {
		return resp.GenericError(err.Error())
	}

	return streams.ToRESP()
}

func parseXRANGEStartEnd(start string, end string) (*redis.XRangeStartEnd, error) {
	parsedRange := &redis.XRangeStartEnd{}

	maxInt := int64(1<<63 - 1)

	if start == "-" {
		parsedRange.StartMS = 0
		parsedRange.StartSeq = 0
	} else if strings.Contains(start, "-") {
		parts := strings.Split(start, "-")
		i, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Invalid stream ID specified as stream command argument")
		}

		s, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Invalid stream ID specified as stream command argument")
		}

		parsedRange.StartMS = i
		parsedRange.StartSeq = s
	} else {
		i, err := strconv.ParseInt(start, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Invalid stream ID specified as stream command argument")
		}

		parsedRange.StartMS = i
		parsedRange.StartSeq = 0
	}

	if end == "+" {
		parsedRange.EndMS = maxInt
		parsedRange.EndSeq = maxInt
	} else if strings.Contains(end, "-") {
		parts := strings.Split(end, "-")
		i, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Invalid stream ID specified as stream command argument")
		}

		s, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Invalid stream ID specified as stream command argument")
		}

		parsedRange.EndMS = i
		parsedRange.EndSeq = s
	} else {
		i, err := strconv.ParseInt(end, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Invalid stream ID specified as stream command argument")
		}

		parsedRange.EndMS = i
		parsedRange.EndSeq = maxInt
	}

	return parsedRange, nil
}

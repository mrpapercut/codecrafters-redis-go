package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

const XRANGE SupportedCommand = "xrange"

func HandleXRANGE(cmd *resp.RESPValue) string {
	if len(cmd.Array) != 4 {
		return resp.GenericError("wrong number of arguments for 'xrange' command")
	}

	key := cmd.Array[1].String

	startIndex, startSeq, endIndex, endSeq, err := parseXRANGEStartEnd(cmd.Array[2].String, cmd.Array[3].String)
	if err != nil {
		return resp.GenericError(err.Error())
	}

	fmt.Printf("XRANGE key: %s, startIndex: %d, startSeq: %d, endIndex: %d, endSeq: %d\n", key, startIndex, startSeq, endIndex, endSeq)

	return ""
}

func parseXRANGEStartEnd(start string, end string) (int64, int64, int64, int64, error) {
	var (
		startIndex int64
		startSeq   int64
		endIndex   int64
		endSeq     int64
	)

	if start == "-" {
		startIndex = 0
		startSeq = 0
	} else if strings.Contains(start, "-") {
		parts := strings.Split(start, "-")
		i, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return 0, 0, 0, 0, fmt.Errorf("Invalid stream ID specified as stream command argument")
		}

		s, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return 0, 0, 0, 0, fmt.Errorf("Invalid stream ID specified as stream command argument")
		}

		startIndex = i
		startSeq = s
	} else {
		i, err := strconv.ParseInt(start, 10, 64)
		if err != nil {
			return 0, 0, 0, 0, fmt.Errorf("Invalid stream ID specified as stream command argument")
		}

		startIndex = i
		startSeq = 0
	}

	if end == "+" {
		endIndex = 1<<63 - 1
		endSeq = 1<<63 - 1
	} else if strings.Contains(end, "-") {
		parts := strings.Split(end, "-")
		i, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return 0, 0, 0, 0, fmt.Errorf("Invalid stream ID specified as stream command argument")
		}

		s, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return 0, 0, 0, 0, fmt.Errorf("Invalid stream ID specified as stream command argument")
		}

		endIndex = i
		endSeq = s
	} else {
		i, err := strconv.ParseInt(end, 10, 64)
		if err != nil {
			return 0, 0, 0, 0, fmt.Errorf("Invalid stream ID specified as stream command argument")
		}

		endIndex = i
		endSeq = 1<<63 - 1
	}

	return startIndex, startSeq, endIndex, endSeq, nil
}

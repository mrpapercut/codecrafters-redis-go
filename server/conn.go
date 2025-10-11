package server

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/codecrafters-io/redis-starter-go/commands"
)

type Conn struct {
}

func Handle(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		buf := make([]byte, 4096)

		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Printf("error reading: %v", err)
			break
		}

		data := buf[:n]

		response, err := commands.HandleCommand(data)
		if err != nil {
			conn.Write([]byte("-ERR " + err.Error() + "\r\n"))
			break
		}

		conn.Write([]byte(response))
	}
}

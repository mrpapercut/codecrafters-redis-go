package server

import (
	"fmt"
	"log"
	"net"
	"os"
)

type Server struct {
	l net.Listener
}

func GetInstance() *Server {
	return &Server{}
}

func (s *Server) StartListening(address string) {
	var err error

	s.l, err = net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("failed to bind to port 6379: %v\n", err)
		os.Exit(1)
	}

	for {
		conn, err := s.l.Accept()
		if err != nil {
			fmt.Printf("error accepting connection: %v\n", err)
			os.Exit(1)
		}

		go Handle(conn)
	}
}

func (s *Server) Close() {
	if s.l != nil {
		err := s.l.Close()
		if err != nil {
			log.Fatalf("error shutting down server: %v", err)
		}

		fmt.Println("server shut down")
	}
}

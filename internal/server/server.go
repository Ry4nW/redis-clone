package server

import (
	"fmt"
	"log"
	"net"

	"redis-clone/internal/command"
)

type Server struct {
	addr     string
	handlers *command.Handlers
}

func New(addr string, handlers *command.Handlers) *Server {
	return &Server{addr: addr, handlers: handlers}
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatal("Error listening:", err)
	}
	defer listener.Close()

	for {
		fmt.Println("")
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting conn:", err)
			continue
		}
		fmt.Println("goroutine started")
		go s.handleConnection(conn)
	}
}

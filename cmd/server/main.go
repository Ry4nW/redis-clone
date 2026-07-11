package main

import (
	"flag"
	"fmt"
	"redis-clone/internal/command"
	"redis-clone/internal/server"
)

func main() {
	port := flag.Int("port", 8091, "port to listen on")
	flag.Parse()

	addr := fmt.Sprintf(":%d", *port)
	handlers := command.NewHandlers()
	s := server.New(addr, handlers)
	s.Start()
}

package main

import (
	"flag"
	"fmt"
	"log"
	"redis-clone/internal/command"
	"redis-clone/internal/server"
)

func main() {
	port := flag.Int("port", 8091, "port to listen on")
	flag.Parse()

	addr := fmt.Sprintf(":%d", *port)

	aof, err := command.New("appendonly.aof", command.Always)
	if err != nil {
		log.Fatalf("failed to open AOF: %v", err)
	}
	defer aof.Close()

	handlers := command.NewHandlersWithAOF(aof)
	if err := aof.Load(handlers); err != nil {
		log.Fatalf("failed to replay AOF: %v", err)
	}

	s := server.New(addr, handlers)
	s.Start()
}

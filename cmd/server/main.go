package main

import "redis-clone/internal/server"

func main() {
	s := server.New(":8091")
	s.Start()
}

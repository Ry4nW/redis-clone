package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func errPrint(err error) {
	if err != nil {
		log.Fatal("Error listening:", err)
	}
}

func main() {

	listener, err := net.Listen("tcp", ":8091")
	if err != nil {
		log.Fatal("Error listening:", err)
	}

	defer listener.Close()

	for {

		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting conn:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	exitMsg := "write 'exit' to exit\n"
	_, err := conn.Write([]byte(exitMsg))
	if err != nil {
		log.Printf("Server write error: %v", err)
	}

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err == io.EOF {
			return
		}

		ackMsg := strings.TrimSpace(message)

		fmt.Println(ackMsg)
		if ackMsg == "EXIT" {
			return
		}

		response := fmt.Sprintf("ACK: %s\n", ackMsg)
		_, err = conn.Write([]byte(response))
		if err != nil {
			log.Printf("Server write error: %v", err)
		}
	}
}

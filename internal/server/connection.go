package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

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
		if ackMsg == "exit" {
			return
		}

		response := fmt.Sprintf("ACK: %s\n", ackMsg)
		_, err = conn.Write([]byte(response))
		if err != nil {
			log.Printf("Server write error: %v", err)
		}
	}
}

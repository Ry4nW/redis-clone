package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"redis-clone/internal/command"
	"redis-clone/internal/resp"
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
		request, err := resp.Parse(reader)

		fmt.Println(ackMsg)
		if err != nil {
			if errors.Is(err, io.EOF) {
				// normal client disconnection
				return
			}

			if errors.Is(err, resp.ErrMalformedRESP) {
				// TODO: encode and send redis-specific erorr response
				// send Redis protocol error
				continue
			}
			log.Printf("parse error: %v", err)
			return
		}

		responseStr, err := command.Dispatch(request)
		if err != nil {
			// convert to Redis error response
			log.Printf("command error: %v", err)

		}

		_, err = conn.Write([]byte(responseStr))
		if err != nil {
			log.Printf("Server write error: %v", err)
		}
	}
}

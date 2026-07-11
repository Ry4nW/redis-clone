package server

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"

	"redis-clone/internal/resp"
)

func Write(conn net.Conn, resp string) {
	_, err := conn.Write([]byte(resp))
	if err != nil {
		log.Printf("Server write error: %v", err)
	} else {
		log.Printf("value written: %v", resp)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	exitMsg := "write 'exit' to exit\n"
	_, err := conn.Write([]byte(exitMsg))
	if err != nil {
		log.Printf("Server write error: %v", err)
	}

	reader := bufio.NewReader(conn)

	for {
		request, err := resp.ParseSimple(reader)

		log.Printf("request made: %v", request)

		if err != nil {
			if errors.Is(err, io.EOF) {
				// normal client disconnection
				return
			}

			if errors.Is(err, resp.ErrMalformedRESP) || errors.Is(err, resp.ErrEmptyRequest) {
				Write(conn, resp.Encode(resp.NewError(err.Error())))
				continue
			}
			log.Printf("parse error: %v", err)
			return
		}

		responseStr, err := s.handlers.Dispatch(request)
		if err != nil {
			log.Printf("command error: %v", err)
		}
		Write(conn, responseStr)
	}
}

package server

import (
	"fmt"
	"net"

	"github.com/palSagnik/httpfromtcp/internal/response"
)

type Server struct {
	Address  string
	Listener net.Listener

	closed bool
}

func Serve(port int) (*Server, error) {
	address := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	server := &Server{Address: address, Listener: l, closed: false}
	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	if err := s.Listener.Close(); err != nil {
		return err
	}
	s.closed = true
	return nil
}

func (s *Server) handle(conn net.Conn) {
	headers := response.GetDefaultHeaders(0)

	response.WriteStatusLine(conn, response.StatusCodeOK)
	response.WriteHeaders(conn, headers)
}

func (s *Server) listen() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			return
		}

		if s.closed {
			return
		}

		go s.handle(conn)
	}
}

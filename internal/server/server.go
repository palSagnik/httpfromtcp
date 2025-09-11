package server

import (
	"fmt"
	"net"
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
	out := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World! ")
	conn.Write(out)
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

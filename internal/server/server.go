package server

import (
	"bytes"
	"fmt"
	"io"
	"net"

	"github.com/palSagnik/httpfromtcp/internal/request"
	"github.com/palSagnik/httpfromtcp/internal/response"
)

type Server struct {
	Address  string
	Listener net.Listener
	HttpHandler	Handler

	closed bool
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError 


func Serve(port int, handler Handler) (*Server, error) {
	address := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	server := &Server{Address: address, Listener: l, HttpHandler: handler, closed: false}
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
	defer conn.Close()
	headers := response.GetDefaultHeaders(0)

	request, err := request.RequestFromReader(conn)
	if err != nil {
		response.WriteStatusLine(conn, response.StatusCodeBadRequest)
		response.WriteHeaders(conn, headers)
		return
	}

	writer := bytes.NewBuffer([]byte{})
	handlerError := s.HttpHandler(writer, request)
	if handlerError != nil {
		bodyBytes := []byte(handlerError.Message)
		headers.Replace("Content-Length", fmt.Sprintf("%d", len(bodyBytes)))

		response.WriteStatusLine(conn, handlerError.StatusCode)
		response.WriteHeaders(conn, headers)
		conn.Write(bodyBytes)
		return
	}
	
	body := writer.Bytes()
	headers.Replace("Content-Length", fmt.Sprintf("%d", len(body)))
	response.WriteStatusLine(conn, response.StatusCodeOK)
	response.WriteHeaders(conn, headers)
	conn.Write(body)
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

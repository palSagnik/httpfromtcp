package server

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/palSagnik/httpfromtcp/internal/request"
	"github.com/palSagnik/httpfromtcp/internal/response"
)

type Server struct {
	Address  string
	Listener net.Listener
	Handler	Handler

	closed atomic.Bool
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w *response.Writer, req *request.Request) 


func Serve(port int, handler Handler) (*Server, error) {
	address := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	server := &Server{Address: address, Listener: l, Handler: handler}
	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.Listener != nil {
		return s.Listener.Close()
	}
	return nil
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	responseWriter := response.NewWriter(conn)

	request, err := request.RequestFromReader(conn)
	if err != nil {
		handlerErr := &HandlerError{
			StatusCode: response.StatusCodeBadRequest,
			Message:    err.Error(),
		}
		handlerErr.Write(responseWriter)
		return
	}

	s.Handler(responseWriter, request)
}

func (s *Server) listen() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			return
		}

		if s.closed.Load() {
			return
		}

		go s.handle(conn)
	}
}

func (he HandlerError) Write(res *response.Writer) {
	res.WriteStatusLine(he.StatusCode)
	messageBytes := []byte(he.Message)
	headers := response.GetDefaultHeaders(len(messageBytes))
	res.WriteHeaders(headers)
	res.WriteBody(messageBytes)
}
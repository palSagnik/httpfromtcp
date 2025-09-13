package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/palSagnik/httpfromtcp/internal/request"
	"github.com/palSagnik/httpfromtcp/internal/response"
	"github.com/palSagnik/httpfromtcp/internal/server"
)

const PORT = 42069

func main() {
	server, err := server.Serve(PORT, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", PORT)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler (w io.Writer, req *request.Request) *server.HandlerError {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandlerError{
			StatusCode: response.StatusCodeBadRequest,
			Message: "Your problem is not my problem\n",
		}

	case "/myproblem":
		return &server.HandlerError{
			StatusCode: response.StatusCodeInternalServerError,
			Message: "Woopsie, my bad\n",
		}

	default:
		w.Write([]byte("All good, frfr\n"))
	}
	
	return nil
}

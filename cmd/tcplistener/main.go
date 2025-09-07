package main

import (
	"fmt"
	"log"
	"net"

	"github.com/palSagnik/httpfromtcp/internal/request"
)

const port = ":42069"
func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error listening: %v", err)
		}
		fmt.Println("Connection Accepted")

		request, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("Error from RequestFromReader: %v", err)
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", request.RequestLine.Method)
		fmt.Printf("- Target: %s\n", request.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", request.RequestLine.HttpVersion)

		fmt.Println("Headers:")
		for key, value := range request.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}

		fmt.Println("Body:")
		fmt.Println(string(request.Body))
	}
}

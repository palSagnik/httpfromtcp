package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
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

		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Println(line)
		}
	}
}

func getLinesChannel(conn io.ReadCloser) <-chan string {
	out := make(chan string)

	go func ()  {
		defer conn.Close()
		defer close(out)

		var line string
		for {
			data := make([]byte, 8)
			n, err := conn.Read(data)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				break
			}

			data = data[:n]
			if i := bytes.IndexByte(data, '\n'); i != -1 {
				line += string(data[:i])
				data = data[i + 1:]
				out <- line
				line = ""
			}
			line += string(data)
		}

		if len(line) > 0 {
			out <- line
		}

	}()
	return out
}
package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

const filePath = "messages.txt"
func main() {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("could not open %s: %s\n", filePath, err)
	}
	defer file.Close()

	fmt.Printf("Reading data from %s\n", filePath)
	fmt.Println("=====================================")

	lines := getLinesChannel(file)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(file io.ReadCloser) <-chan string {
	out := make(chan string)

	go func ()  {
		defer file.Close()
		defer close(out)

		var line string
		for {
			data := make([]byte, 8)
			n, err := file.Read(data)
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
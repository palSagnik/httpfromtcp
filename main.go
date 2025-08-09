package main

import (
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

	for {
		buf := make([]byte, 8)
		_, err := file.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Printf("error: %s\n", err.Error())
			break
		}
		fmt.Printf("read: %s\n", buf)
	}
}
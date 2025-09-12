package response

import (
	"fmt"
	"io"

	"github.com/palSagnik/httpfromtcp/internal/headers"
)

func GetDefaultHeaders(contentlen int) headers.Headers {
	headers := headers.NewHeaders()
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")
	headers.Set("Content-Length", fmt.Sprintf("%d", contentlen))

	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return err
}
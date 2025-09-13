package response

import (
	"io"

	"github.com/palSagnik/httpfromtcp/internal/headers"
)

type Writer struct {
	Writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{Writer: w}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	return writeStatusLine(w.Writer, statusCode)
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	return writeHeaders(w.Writer, headers)
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	return w.Writer.Write(p)
}
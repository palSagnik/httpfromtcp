package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type ParserState int

const (
	STATE_INITIALISED ParserState = iota
	STATE_DONE
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	state      ParserState
}

func newRequest() *Request {
	return &Request{
		state: STATE_INITIALISED,
	}
}

const BUFFER_SIZE = 8
const SEPARATOR = "\r\n"

var ErrorMalformedStartLine = fmt.Errorf("malformed start-line")
var ErrorUnsupportedHttpVersion = fmt.Errorf("unrecognised http version")

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	buf := make([]byte, 2048)
	bufLen := 0

	for !request.done() {
		// Read from the reader and store in the buffer
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			if err == io.EOF {
				request.state = STATE_DONE
				break
			}
			return nil, err
		}

		bufLen += n
		readN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}

	return request, nil
}

func parseRequestLine(data []byte) (int, *RequestLine, error) {
	idx := bytes.Index(data, []byte(SEPARATOR))
	if idx == -1 {
		return 0, nil, nil
	}

	reqLine := data[:idx]
	read := idx + len(SEPARATOR)

	parts := bytes.Split(reqLine, []byte(" "))
	fmt.Println(parts)
	if len(parts) != 3 {
		return read, nil, ErrorMalformedStartLine
	}

	method := string(parts[0])
	target := string(parts[1])
	httpVersion := string(parts[2])

	// method
	if strings.ToUpper(method) != method {
		return read, nil, errors.New("method must be all capital letters")
	}

	// version
	httpParts := strings.Split(httpVersion, "/")
	if len(httpParts) != 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		if httpParts[1] != "1.1" {
			return read, nil, ErrorUnsupportedHttpVersion
		}
		return read, nil, ErrorMalformedStartLine
	}

	return read, &RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   httpParts[1],
	}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
	outer:
	for {
		switch r.state {
		case STATE_INITIALISED:
			n, reqLine, err := parseRequestLine(data[read:])
			if err != nil {
				return 0, err
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = *reqLine
			read += n

			r.state = STATE_DONE
		
		case STATE_DONE:
			break outer	
		}
	}

	return read, nil
}

func (r *Request) done() bool {
	return r.state == STATE_DONE
}


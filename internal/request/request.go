package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	idx := bytes.Index(buf, []byte(crlf))
	if idx == -1 {
		return nil, errors.New("no crlf found in request line")
	}

	reqLine := string(buf[:idx])
	requestLine, err := parseRequestLine(reqLine)
	if err != nil {
		return nil, err
	}

	return &Request{RequestLine: *requestLine}, nil
}

func parseRequestLine(reqLine string) (*RequestLine, error) {
	parts := strings.Split(reqLine, " ")
	if len(parts) != 3 {
		return nil, errors.New("poorly formatted request-line")
	}

	method := string(parts[0])
	target := string(parts[1])
	httpVersion := string(parts[2])

	// method
	if strings.ToUpper(method) != method {
		return nil, errors.New("method must be all capital letters")
	}

	// version
	versionParts := strings.Split(httpVersion, "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", reqLine)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}
	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		Method: method,
		RequestTarget: target,
		HttpVersion: version,
	}, nil
}

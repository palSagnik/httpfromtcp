package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

const SEPARATOR = "\r\n"

var ErrorMalformedHeader = fmt.Errorf("malformed header data")

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(SEPARATOR))
	if idx == -1 {
		return 0, false, nil
	}
	
	// CRLF found at the start => end of headers
	// consume CRLF
	if idx == 0 {
		return 2, true, nil
	}

	parts := bytes.SplitN(data[:idx], []byte(":"), 2)

	// remove whitespaces from start and end
	if parts[0][len(parts) - 1] == byte(' ') {
		return 0, false, ErrorMalformedHeader
	}

	value := string(bytes.TrimSpace(parts[1]))
	key := string(bytes.TrimSpace(parts[0]))
	
	h.Set(key, value)
	
	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	h[key] = value
}
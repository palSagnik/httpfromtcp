package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

const SEPARATOR = "\r\n"

var ErrorMalformedHeader = fmt.Errorf("malformed header data")
var ErrorInvalidFieldName = fmt.Errorf("invalid field name")

func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false
	for {
		idx := bytes.Index(data[read:], []byte(SEPARATOR))
		if idx == -1 {
			break 
		}
		
		// CRLF found at the start => end of headers
		// consume CRLF
		if idx == 0 {
			done = true
			read += len(SEPARATOR)
			break
		}

		key, value, err := parseHeader(data[read:read+idx])
		if err != nil {
			return 0, false, err
		}
		
		read += idx + len(SEPARATOR)
		h.Set(key, value)
	}

	return read, done, nil
}

func parseHeader(data []byte) (string, string, error) {
	parts := bytes.SplitN(data, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", ErrorMalformedHeader
	}

	// remove whitespaces from start and end
	if parts[0][len(parts) - 1] == byte(' ') {
		return "", "", ErrorMalformedHeader
	}

	value := string(bytes.TrimSpace(parts[1]))
	name := string(bytes.TrimSpace(parts[0]))

	if !validFieldName(name) {
		return "", "", ErrorInvalidFieldName
	}

	return name, value, nil
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	v, ok := h[key]
	if ok {
		value = strings.Join([]string{v, value}, ", ")
	}
	h[key] = value
}

func validFieldName(key string) bool {
	for _, c := range key {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && !checkSpecial(c) {
			fmt.Println(key)
			return false
		}
	}
	return true
}

func checkSpecial(c rune) bool {
	allowed := []rune{'!', '#', '$', '%', '&', '*', '+', '-', '.', '^', '_', '`', '|', '~', '\''}
	
	for _, special := range allowed {
		if special == c {
			return true
		}
	}
	return false
}
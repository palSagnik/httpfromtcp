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

	if !validFieldName(key) {
		return 0, false, ErrorInvalidFieldName
	}

	h.Set(key, value)

	return idx + 2, false, nil
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
package requestline

import (
	"bytes"
	"errors"
)

// TODO: need to add tests
type PARSER_STATE string

const (
	INVALID_REQUEST = "invalid request"
)

var SEPARATOR = []byte("\r\n")

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func ParseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := b[:idx]
	read := idx + len(SEPARATOR)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, errors.New(INVALID_REQUEST)
	}

	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, errors.New(INVALID_REQUEST)
	}

	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(parts[2]),
	}

	return rl, read, nil
}

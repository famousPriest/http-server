package request

import (
	"bytes"
	"errors"
	"io"
)

const (
	INVALID_REQUEST                   = "invalid request"
	INVALID_METHOD_NAME               = "invalid method name"
	INVALID_HTTP_VERSION              = "invalid HTTP version"
	MISSING_SEPARATOR                 = "missing separator"
	REQUIRED_ELEMENTS_IN_REQUEST_LINE = "request line must contain three elements: method name, request target and HTTP version"

	DEFAULT_HTTP_VERSION = "HTTP/1.1"

	BUFFER_SIZE             = 8
	StateInit   parserState = "init"
	StateDone   parserState = "done"
)

var SEPARATOR = []byte("\r\n")

type parserState string

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	State       parserState
}

func newRequest() *Request {
	return &Request{
		State: StateInit,
	}
}

func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n

	return n, nil
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		switch r.State {
		case StateInit:
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				return 0, nil
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n

			r.State = StateDone

		case StateDone:
			break outer
		}
	}
	return 0, nil
}

func (r *Request) isDone() bool {
	return r.State == StateDone
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
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

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	buf := make([]byte, BUFFER_SIZE)
	bufLen := 0
	for !request.isDone() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			return nil, err
		}

		bufLen += n
		readN, err := request.parse(buf[:bufLen+n])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}

	return request, nil
}

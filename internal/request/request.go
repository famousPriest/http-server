package request

import (
	"bytes"
	"errors"
	"i-hate-js/internal/body"
	"i-hate-js/internal/headers"
	"io"
	"strconv"
)

// TODO: refactor to make it more modular and testable
type PARSER_STATE string

const (
	INVALID_REQUEST                   = "invalid request"
	INVALID_HEADER_FORMAT             = "invalid header format"
	INVALID_METHOD_NAME               = "invalid method name"
	INVALID_HTTP_VERSION              = "invalid HTTP version"
	MISSING_SEPARATOR                 = "missing separator"
	REQUIRED_ELEMENTS_IN_REQUEST_LINE = "request line must contain three elements: method name, request target and HTTP version"

	DEFAULT_HTTP_VERSION = "HTTP/1.1"
	CONTENT_LENGTH       = "Content-Length"

	BUFFER_SIZE                     = 1024
	STATE_INIT         PARSER_STATE = "init"
	STATE_DONE         PARSER_STATE = "done"
	STATE_PARSING_BODY PARSER_STATE = "parsing body"
)

var SEPARATOR = []byte("\r\n")

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
	RequestLine    RequestLine
	RequestHeaders headers.Headers
	State          PARSER_STATE
	Body           body.Body
}

func newRequest() *Request {
	return &Request{
		State: STATE_INIT,
	}
}

func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	endIndex = min(endIndex, len(cr.data))
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n

	return n, nil
}

// TODO: need to refactor error handling here, currently always returns 0 as a number of read bytes
// TODO: need to consider cases of body not being fully read yet
func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		switch r.State {
		case STATE_INIT:
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				return 0, nil
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n

			h, n, err := parseHeaders(data[read:])
			if err != nil {
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.RequestHeaders = h
			r.State = STATE_PARSING_BODY

		case STATE_PARSING_BODY:
			contentLengthStr := r.RequestHeaders.GetValue(CONTENT_LENGTH)

			if contentLengthStr == "" {
				r.State = STATE_DONE
				break
			}

			contentLength, err := strconv.Atoi(contentLengthStr)
			if err != nil {
				return 0, err
			}

			b, err := parseBody(data[read:], contentLength)
			if err != nil {
				return 0, err
			}

			r.Body = b
			r.State = STATE_DONE

		case STATE_DONE:
			break outer
		}
	}
	return read, nil
}

func (r *Request) isDone() bool {
	return r.State == STATE_DONE
}

func parseBody(b []byte, length int) (body.Body, error) {
	body := body.Body{}
	body.Parse(b, length)
	return body, nil
}

func parseHeaders(b []byte) (headers.Headers, int, error) {
	headers := make(headers.Headers)
	totalRead := 0

	for {
		n, done, err := headers.Parse(b[totalRead:])
		if err != nil {
			return nil, 0, errors.New(INVALID_HEADER_FORMAT)
		}
		totalRead += n

		if done {
			break
		}
	}

	return headers, totalRead, nil
}

// TODO: should move parser for request lines to a different file. This file should be used for handling requests only
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
		bufLen += n

		if bufLen > 0 {
			readN, parseErr := request.parse(buf[:bufLen])
			if parseErr != nil {
				return nil, parseErr
			}

			copy(buf, buf[readN:bufLen])
			bufLen -= readN
		}

		if err != nil {
			if err == io.EOF && request.isDone() {
				return request, nil
			}
			return nil, err
		}
	}

	return request, nil
}

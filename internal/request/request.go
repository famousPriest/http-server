package request

import (
	"i-hate-js/helpers"
	"i-hate-js/internal/body"
	"i-hate-js/internal/headers"
	reqLine "i-hate-js/internal/request-line"
	"io"
	"strconv"
)

type PARSER_STATE string

const (
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

type Request struct {
	RequestLine    reqLine.RequestLine
	RequestHeaders headers.Headers
	State          PARSER_STATE
	Body           body.Body
}

func (r *Request) isDone() bool {
	return r.State == STATE_DONE
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
			rl, n, err := reqLine.ParseRequestLine(data[read:])
			if err != nil {
				return 0, nil
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n

			h, n, err := helpers.ParseHeaders(data[read:])
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

			b, err := helpers.ParseBody(data[read:], contentLength)
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

func newRequest() *Request {
	return &Request{
		State: STATE_INIT,
	}
}

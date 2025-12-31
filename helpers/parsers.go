package helpers

import (
	"errors"
	"i-hate-js/internal/body"
	"i-hate-js/internal/headers"
)

const (
	INVALID_HEADER_FORMAT = "invalid header format"
)

func ParseBody(b []byte, length int) (body.Body, error) {
	body := body.Body{}
	body.Parse(b, length)
	return body, nil
}

func ParseHeaders(b []byte) (headers.Headers, int, error) {
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

package headers

import (
	"bytes"
	"errors"
)

type Headers map[string]string

func createMockHeader() Headers {
	return Headers{
		"Content-Type":    "application/json",
		"User-Agent":      "GoMockClient/1.0",
		"Authorization":   "Bearer mock_token_123",
		"Accept":          "*/*",
		"X-Custom-Header": "mock-value",
	}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	crlfIdx := bytes.Index(data, []byte("\r\n"))

	switch crlfIdx {
	case -1:
		return 0, false, nil
	case 0:
		return 0, true, nil
	}

	data = data[:crlfIdx]

	parts := bytes.SplitN(data, []byte(":"), 2)
	keyBytes := parts[0]
	valueBytes := parts[1]

	if bytes.Contains(keyBytes, []byte(" ")) || len(keyBytes) == 0 {
		return crlfIdx + 2, false, errors.New("invalid host name")
	}

	keyString := string(keyBytes)
	valueString := string(bytes.TrimSpace(valueBytes))

	h[keyString] = string(valueString)

	return len(data), false, nil
}

package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"
)

const (
	INVALID_HOST_NAME       = "invalid host name"
	INVALID_CHAR_IN_REQUEST = "invalid characters in request line"
)

var INVALID_CHARS = []string{
	" ", "!", "#", "$", "'", "*", "+", "-", "=", ".", "_", "`", "|", "~", "@",
}

type Headers map[string]string

func newHeaders() Headers {
	return Headers{
		"Content-Type":    "application/json",
		"User-Agent":      "GoMockClient/1.0",
		"Authorization":   "Bearer mock_token_123",
		"Accept":          "*/*",
		"X-Custom-Header": "mock-value",
		"Host":            "localhost:42069\r\n\r\n",
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
		return crlfIdx + 2, false, errors.New(INVALID_HOST_NAME)
	}

	for i, s := range INVALID_CHARS {
		if bytes.Contains(keyBytes, []byte(s)) {
			return crlfIdx + 2, false, errors.New(INVALID_CHAR_IN_REQUEST)
		}
		i++
	}

	key := strings.ToLower(string(keyBytes))
	value := strings.Trim(string(valueBytes), " \r\n")

	isFound := false
	for k, v := range h {
		if strings.ToLower(k) == key {
			h[k] = strings.Trim(string(v), " \r\n") + ", " + value
			isFound = true
			break
		}
	}

	if !isFound {
		h[key] = value
	}

	for k, v := range h {
		fmt.Printf("KEY: %s VALUE: %s\n", k, v)
	}
	return len(data), false, nil
}

func main() {
	headers := newHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)

	if err != nil {
		log.Fatalln("oopsey")
	}
	fmt.Print(n)
	fmt.Print(done)
	fmt.Print(err)
}

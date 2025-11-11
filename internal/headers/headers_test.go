package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newHeaders() Headers {
	return Headers{
		"Content-Type":    "application/json",
		"User-Agent":      "GoMockClient/1.0",
		"Authorization":   "Bearer mock_token_123",
		"Accept":          "*/*",
		"X-Custom-Header": "mock-value",
	}
}

func TestHeaders(t *testing.T) {
	headers := newHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = newHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

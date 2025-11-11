package request

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequestLineParse(t *testing.T) {
	raw := "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"
	reader := &chunkReader{
		data:            raw,
		numBytesPerRead: 3,
	}

	req, err := RequestFromReader(reader)
	if err != nil {
		log.Fatalln("oops")
	}

	fmt.Println("r: ", req)
	require.NoError(t, err)
	require.NotNil(t, req)
}

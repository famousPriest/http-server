package main

import (
	"fmt"
	"i-hate-js/internal/request"
	"log"
	"net"
)

func main() {

	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalln("oopsey")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("oopsey")
		}

		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalln("oopsey")
		}

		fmt.Printf("request line:\n")
		fmt.Printf("- method: %s\n", r.RequestLine.Method)
		fmt.Printf("- target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- version: %s\n", r.RequestLine.HttpVersion)
	}
}

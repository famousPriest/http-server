package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	address, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		log.Fatalln("oopsey on resolve")
	}

	conn, err := net.DialUDP("udp", nil, address)
	if err != nil {
		log.Fatalln("oopsey on dial")
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("type smth: ")
	for {
		fmt.Println(">")

		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln("oopsey on ReadString")
			continue
		}

		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Fatalln("oopsey on Write")
			continue
		}
	}

}

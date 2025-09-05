package main

import (
	"fmt"
	"net"
	"bufio"
	"os"
	"log"

)

func main() {
	raddr, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		log.Fatal("error resolving udp addr at :42069")
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	defer conn.Close()
	if err != nil {
		log.Fatal("error dialling up udp addr at :42069")
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf(">")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("error reading string")
		}
		
		_, err = conn.Write([]byte(input))
		if err != nil {
			log.Fatal("error writing to connection")
		}
	}
}

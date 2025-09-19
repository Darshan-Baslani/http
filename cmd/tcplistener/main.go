package main

import (
	"fmt"
	"log"
	"net"

	"http/internal/request"
)


func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error listening to port :42069")
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error aceepting request")
		}
		req, err := request.RequestFromReader(conn)
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)
	}

}

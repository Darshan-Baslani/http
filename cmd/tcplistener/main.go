package main

import (
	"bytes"
	"fmt"
	"log"
	"io"
	"net"
)

func getLinesChannel(file io.ReadCloser) <-chan string {
	data := make(chan string)
	go func() {
		str := ""
		for {
			stream := make([]byte, 8)
			readSize, err := file.Read(stream)
			if err == io.EOF {
				close(data)
				break
			}
			stream = stream[:readSize]
			if i := bytes.IndexByte(stream, byte('\n')); i != -1 {
				str += string(stream[:i])
				// fmt.Printf("read: %s\n", str)
				stream = stream[i+1:]
				data <- str
				str = ""
			}
			str += string(stream)

			// fmt.Printf("read: %s\n", string(stream))
		}
	}()
	return data
}

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
		for data := range getLinesChannel(conn)	{
			fmt.Printf("read: %s\n", data)
		}
	}

}

package request

import (
	"io"
	"strings"
	"fmt"

	"http/utils"
)

type parserState int

const (
	initialized parserState = iota
	done
)

type Request struct {
	RequestLine RequestLine
	state       parserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type bufferStatus struct {
	bytesRead int
	bytesParsed int
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := new(Request)
	req.state = initialized
	buffStatus := new(bufferStatus)
	bufferSize := 8
	for {
		buffer := make([]byte, bufferSize)
		readSize, err := reader.Read(buffer)
		if err == io.EOF {
			req.state = done
			break
		}
		if err != nil {
			return nil, fmt.Errorf("cannot read the reader")
		}
		buffStatus.bytesRead += readSize

		buffer = buffer[:readSize]
		bytesParsed, err := req.parse(&buffer)
		if err != nil {
			return nil, fmt.Errorf("%s", err)
		}
		buffStatus.bytesParsed += bytesParsed
		buffer = append(buffer, make([]byte, bufferSize)...)
	}

	return req, nil
}

func (r *Request) parse(data *[]byte) (int, error) {
	n, err := parseRequestLine(r, string(data))
	if err != nil {
		return -1, fmt.Errorf("%s", err)
	}

	return n, nil
}

func parseRequestLine(req *Request, msg string) (int, error){
	fmt.Printf("raw msg: %q\n", msg)
	idx := strings.Index(msg, "\r\n")
	if idx != -1 {
		msg = msg[:idx]
	} else {
		return 0, nil
	}

	parts := strings.Split(msg, " ")
	if len(parts) != 3 {
		return -1, fmt.Errorf("invalid number of parts in request line: %d, string received: %s", len(parts), msg)
	}
	method, reqtarget, httpV := parts[0], parts[1], parts[2]

	if (!utils.IsAllUpper(method)){
		return -1, fmt.Errorf("method needs to be all upper case")
	}

	httpV = strings.TrimPrefix(httpV, "HTTP/")
	if (httpV != "1.1") {
		return -1, fmt.Errorf("http version not supported")
	}

	req.RequestLine = RequestLine{httpV, reqtarget, method}

	return idx-1, nil
}

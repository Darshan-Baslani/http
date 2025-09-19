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
	bufferSize := 16384

	buffer := make([]byte, bufferSize)
	for {
		readSize, err := reader.Read(buffer[buffStatus.bytesRead:])
		if err == io.EOF {
			req.state = done
			break
		}
		if err != nil {
			return nil, fmt.Errorf("cannot read the reader")
		}
		fmt.Printf("after read : %s\n", buffer)
		buffStatus.bytesRead += readSize

		bytesParsed, err := req.parse(buffer[:buffStatus.bytesRead])
		if err != nil {
			return nil, fmt.Errorf("%s", err)
		}

		if bytesParsed == 0 && len(buffer) < buffStatus.bytesRead+readSize{
			buffer = append(buffer, make([]byte, readSize)...)
		} else if bytesParsed > 0{
			buffStatus.bytesParsed += bytesParsed
			req.state = done
			break
		}
		fmt.Printf("buffer len: %d\n", len(buffer))
	}

	return req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	fmt.Printf("data: %s\n", data)
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

	return len(msg), nil
}

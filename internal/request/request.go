package request

import (
	"io"
	"strings"
	"fmt"

	"http/utils"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	msg, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("cannot read the reader")
	}

	reqline, err := parseRequestLine(string(msg))
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	return &Request{*reqline}, nil
}

func parseRequestLine(msg string) (*RequestLine, error){
	if idx := strings.Index(msg, "\r\n"); idx != -1 {
		msg = msg[:idx]
	} else {
		return nil, fmt.Errorf("empty string")
	}

	parts := strings.Split(msg, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid number of parts in request line: %d", len(parts))
	}
	method, reqtarget, httpV := parts[0], parts[1], parts[2]

	if (!utils.IsAllUpper(method)){
		return nil, fmt.Errorf("method needs to be all upper case")
	}

	httpV = strings.TrimPrefix(httpV, "HTTP/")
	if (httpV != "1.1") {
		return nil, fmt.Errorf("http version not supported")
	}

	return &RequestLine{httpV, reqtarget, method}, nil
}

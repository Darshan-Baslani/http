package headers

import (
	"bytes"
	"fmt"
	"strings"

	"http/utils"
)

const (
	crlf = "\r\n"
)

type Headers map[string]string

func NewHeaders() Headers {
	h := make(Headers)
	return h
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	done := false
	readSize := 0
	for {
		crlfPos := bytes.IndexAny(data, crlf)
		if crlfPos == -1 {
			break
		}
		// double crlf condition
		if crlfPos == 0 {
			done = true
			readSize += 1
			break
		}
		readSize += crlfPos + 1

		key, value, err := parseHeader(data[:crlfPos])
		if err != nil {
			readSize = 0
			return readSize, done, fmt.Errorf("%s", err)
		}
		h[key] = value
		data = data[crlfPos+len(crlf):]
	}

	return readSize, done, nil
}

func parseHeader(data []byte) (string, string, error) {
	colon := bytes.IndexAny(data, ":")
	if data[colon-1] == ' ' {
		return "", "", fmt.Errorf("unstructured request")
	}

	key := strings.ToLower(string(bytes.TrimSpace(data[:colon])))
	if !utils.IsValidFieldName(key) {
		return "", "", fmt.Errorf("field name contains invalid chars")
	}

	value := string(bytes.TrimSpace(data[colon+1:]))

	return key, value, nil
}

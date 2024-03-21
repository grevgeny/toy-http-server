package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Request struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    string
}

func readRequest(conn net.Conn) (*Request, error) {
	reader := bufio.NewReader(conn)
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading request line: %w", err)
	}

	parts := strings.Split(strings.TrimSpace(requestLine), " ")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid request line")
	}

	request := &Request{
		Method:  parts[0],
		Path:    parts[1],
		Headers: make(map[string]string),
	}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("error reading header: %w", err)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}

		keyVal := strings.SplitN(line, ": ", 2)
		if len(keyVal) == 2 {
			request.Headers[keyVal[0]] = keyVal[1]
		}
	}

	return request, nil
}

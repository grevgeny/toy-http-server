package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

type Request struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    []byte
}

func ParseRequest(conn net.Conn) (*Request, error) {
	reader := bufio.NewReader(conn)

	// Read the request line
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading request line: %w", err)
	}

	parts := strings.Fields(strings.TrimSpace(requestLine))
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid request line")
	}

	// Extract headers
	headers, err := extractHeaders(reader)
	if err != nil {
		return nil, fmt.Errorf("error extacting headers: %w", err)
	}

	request := &Request{
		Method:  parts[0],
		Path:    parts[1],
		Headers: headers,
	}

	// Extract body for POST and PUT methods
	if request.Method == "POST" {
		body, err := extractBody(reader, headers["Content-Length"])
		if err != nil {
			return nil, fmt.Errorf("error extracting body: %w", err)
		}

		request.Body = body
	}

	return request, nil
}

func extractHeaders(reader *bufio.Reader) (map[string]string, error) {
	headers := make(map[string]string)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("error reading header: %w", err)
		}

		if line := strings.TrimSpace(line); line == "" {
			break
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header line: %s", line)
		}

		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		headers[key] = value
	}

	return headers, nil
}

func extractBody(reader *bufio.Reader, lenHeader string) ([]byte, error) {
	contentLength, err := strconv.Atoi(lenHeader)
	if err != nil {
		return nil, fmt.Errorf("error parsing Content-Length: %w", err)
	}

	if contentLength < 0 {
		return nil, fmt.Errorf("invalid Content-Length: %d", contentLength)
	}

	bodyBytes := make([]byte, contentLength)
	_, err = io.ReadFull(reader, bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}

	return bodyBytes, nil
}

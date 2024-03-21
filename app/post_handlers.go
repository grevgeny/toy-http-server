package main

import (
	"net"
	"os"
	"strings"
)

func handlePost(conn net.Conn, request *Request) {
	switch {

	case strings.HasPrefix(request.Path, "/files/"):
		handleFilePost(conn, request)
	}
}

func handleFilePost(conn net.Conn, request *Request) {
	filename := strings.TrimPrefix(request.Path, "/files/")
	filePath := directory + "/" + filename

	err := os.WriteFile(filePath, request.Body, 0644)
	exitOnError(err, "Error writing file")

	conn.Write([]byte("HTTP/1.1 201 Created\r\n\r\n"))
}

package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func handleRoot(conn net.Conn) {
	writeResponseOK(conn, "", "")
}

func handleEcho(conn net.Conn, path string) {
	writeResponseOK(conn, strings.TrimPrefix(path, "/echo/"), "text/plain")
}

func handleUserAgent(conn net.Conn, userAgent string) {
	writeResponseOK(conn, userAgent, "text/plain")
}

func handleFile(conn net.Conn, path string) {
	filename := strings.TrimPrefix(path, "/files/")
	file, err := os.ReadFile(directory + "/" + filename)

	if err != nil {
		writeResponseNotFound(conn)
		return
	}

	response := string(file)
	writeResponseOK(conn, response, "application/octet-stream")
}

func writeResponseOK(conn net.Conn, response string, content_type string) {
	conn.Write([]byte("HTTP/1.1 200 OK\r\n"))

	if response == "" {
		conn.Write([]byte("\r\n"))
		return
	}

	conn.Write([]byte("Content-Type: " + content_type + "\r\n"))
	conn.Write([]byte(fmt.Sprint("Content-Length: ", len(response), "\r\n")))
	conn.Write([]byte("\r\n"))
	conn.Write([]byte(response))
}

func writeResponseNotFound(conn net.Conn) {
	conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
}

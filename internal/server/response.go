package server

import (
	"fmt"
	"net"
)

func WriteResponseOK(conn net.Conn, response string, content_type string, encoding string) {
	conn.Write([]byte("HTTP/1.1 200 OK\r\n"))

	if response == "" {
		conn.Write([]byte("\r\n"))
		return
	}

	conn.Write([]byte("Content-Type: " + content_type + "\r\n"))
	if encoding != "" {
		switch encoding {
		case "gzip":
			conn.Write([]byte("Content-Encoding: " + encoding + "\r\n"))
		}
	}
	conn.Write([]byte(fmt.Sprint("Content-Length: ", len(response), "\r\n")))
	conn.Write([]byte("\r\n"))
	conn.Write([]byte(response))
}

func WriteResponseCreated(conn net.Conn) {
	conn.Write([]byte("HTTP/1.1 201 Created\r\n\r\n"))
}

func WriteResponseBad(conn net.Conn) {
	conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
}

func WriteResponseNotFound(conn net.Conn) {
	conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
}

func WriteResponseNowAllowed(conn net.Conn) {
	conn.Write([]byte("HTTP/1.1 405 Method Not Allowed\r\n\r\n"))
}

func WriteResponseError(conn net.Conn) {
	conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
}

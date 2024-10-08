package server

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net"
	"strings"
)

func WriteResponseOK(conn net.Conn, response string, content_type string, encoding string) {
	conn.Write([]byte("HTTP/1.1 200 OK\r\n"))

	if response == "" {
		conn.Write([]byte("\r\n"))
		return
	}

	conn.Write([]byte("Content-Type: " + content_type + "\r\n"))

	if encoding != "" {
		for _, enc := range strings.Split(encoding, ", ") {
			switch enc {
			case "gzip":
				conn.Write([]byte("Content-Encoding: " + enc + "\r\n"))
				response = compressResponse(response)
				break
			}
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

func compressResponse(response string) string {
	var buffer bytes.Buffer
	gz := gzip.NewWriter(&buffer)
	gz.Write([]byte(response))
	gz.Close()
	return buffer.String()
}

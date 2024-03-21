package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func handleConnection(conn net.Conn) {
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	exitOnError(err, "Error reading from connection")

	parsed_req := strings.Split(string(buffer), "\r\n")
	path := strings.Split(parsed_req[0], " ")[1]

	if path == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if strings.HasPrefix(path, "/echo/") {
		content, _ := strings.CutPrefix(path, "/echo/")

		conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
		conn.Write([]byte("Content-Type: text/plain\r\n"))
		conn.Write([]byte(fmt.Sprint("Content-Length: ", len(content), "\r\n")))
		conn.Write([]byte("\r\n"))
		conn.Write([]byte(content))

	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}

}

func exitOnError(err error, message string) {
	if err == nil {
		return
	}

	fmt.Printf("%s: %s", message, err.Error())
	os.Exit(1)
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	exitOnError(err, "Failed to bind to port 4221")

	defer l.Close()

	for {
		conn, err := l.Accept()
		exitOnError(err, "Error accepting connection")

		go handleConnection(conn)
	}

}

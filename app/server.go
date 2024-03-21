package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var directory string

func main() {
	parseFlags()

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	exitOnError(err, "Failed to bind to port 4221")

	defer l.Close()

	for {
		conn, err := l.Accept()
		exitOnError(err, "Error accepting connection")

		go handleConnection(conn)
	}

}

func parseFlags() {
	flag.StringVar(&directory, "directory", "", "directory containing files")
	flag.Parse()
}

func exitOnError(err error, message string) {
	if err == nil {
		return
	}

	fmt.Printf("%s: %s", message, err.Error())
	os.Exit(1)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	request, err := readRequest(conn)
	if err != nil {
		exitOnError(err, "Error reading request")
	}

	switch {
	case request.Path == "/":
		handleRoot(conn)

	case strings.HasPrefix(request.Path, "/echo/"):
		handleEcho(conn, request.Path)

	case strings.HasPrefix(request.Path, "/user-agent"):
		handleUserAgent(conn, request.Headers["User-Agent"])

	case strings.HasPrefix(request.Path, "/files/"):
		handleFile(conn, request.Path)

	default:
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}

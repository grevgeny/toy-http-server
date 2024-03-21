package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var directory string

func parseFlags() {
	flag.StringVar(&directory, "directory", "", "directory containing source files")
	flag.Parse()
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	parsed_request := readRequest(conn)
	path := strings.Split(parsed_request[0], " ")[1]

	switch {
	case path == "/":
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

	case strings.HasPrefix(path, "/echo/"):
		handleEcho(conn, path)

	case strings.HasPrefix(path, "/user-agent"):
		handleUserAgent(conn, parsed_request)

	case strings.HasPrefix(path, "/files/"):
		handleFile(conn, path)

	default:
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}

func readRequest(conn net.Conn) []string {
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	exitOnError(err, "Error reading from connection")

	return strings.Split(string(buffer), "\r\n")
}

func handleEcho(conn net.Conn, path string) {
	content, _ := strings.CutPrefix(path, "/echo/")

	conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
	conn.Write([]byte("Content-Type: text/plain\r\n"))
	conn.Write([]byte(fmt.Sprint("Content-Length: ", len(content), "\r\n")))
	conn.Write([]byte("\r\n"))
	conn.Write([]byte(content))
}

func handleUserAgent(conn net.Conn, parsed_req []string) {
	for _, line := range parsed_req {
		if !strings.HasPrefix(line, "User-Agent") {
			continue
		}

		userAgent := strings.TrimPrefix(line, "User-Agent: ")

		conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
		conn.Write([]byte("Content-Type: text/plain\r\n"))
		conn.Write([]byte(fmt.Sprint("Content-Length: ", len(userAgent), "\r\n")))
		conn.Write([]byte("\r\n"))
		conn.Write([]byte(userAgent))

		break
	}
}

func handleFile(conn net.Conn, path string) {
	filename := strings.TrimPrefix(path, "/files/")
	file, err := os.ReadFile(directory + "/" + filename)

	if err != nil {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		return
	}

	conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
	conn.Write([]byte("Content-Type: application/octet-stream\r\n"))
	conn.Write([]byte(fmt.Sprint("Content-Length: ", len(file), "\r\n")))
	conn.Write([]byte("\r\n"))
	conn.Write(file)
}

func exitOnError(err error, message string) {
	if err == nil {
		return
	}

	fmt.Printf("%s: %s", message, err.Error())
	os.Exit(1)
}

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

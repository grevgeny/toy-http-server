package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var directory string

func readRequest(conn net.Conn) []string {
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	exitOnError(err, "Error reading from connection")

	return strings.Split(string(buffer), "\r\n")
}

func handleEcho(conn net.Conn, path string) {
	response, _ := strings.CutPrefix(path, "/echo/")
	writeResponseOK(conn, response, "text/plain")
}

func handleUserAgent(conn net.Conn, parsed_req []string) {
	for _, line := range parsed_req {
		if !strings.HasPrefix(line, "User-Agent") {
			continue
		}

		response := strings.TrimPrefix(line, "User-Agent: ")
		writeResponseOK(conn, response, "text/plain")

		return
	}
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
	conn.Write([]byte("Content-Type: " + content_type + "\r\n"))
	conn.Write([]byte(fmt.Sprint("Content-Length: ", len(response), "\r\n")))
	conn.Write([]byte("\r\n"))
	conn.Write([]byte(response))
}

func writeResponseNotFound(conn net.Conn) {
	conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
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

	request_components := readRequest(conn)
	path := strings.Split(request_components[0], " ")[1]

	switch {
	case path == "/":
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

	case strings.HasPrefix(path, "/echo/"):
		handleEcho(conn, path)

	case strings.HasPrefix(path, "/user-agent"):
		handleUserAgent(conn, request_components)

	case strings.HasPrefix(path, "/files/"):
		handleFile(conn, path)

	default:
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}

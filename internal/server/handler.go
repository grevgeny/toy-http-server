package server

import (
	"net"
	"os"
	"path/filepath"
	"strings"
)

type Handler struct {
	directory string
}

func NewHandler(directory string) (*Handler, error) {
	return &Handler{directory: directory}, nil
}

func (h *Handler) ServeHTTP(conn net.Conn, req *Request) {
	switch req.Method {
	case "GET":
		h.handleGet(conn, req)
	case "POST":
		h.handlePost(conn, req)
	default:
		WriteResponseNowAllowed(conn)
	}
}

func (h *Handler) handleGet(conn net.Conn, req *Request) {
	switch {
	case req.Path == "/":
		WriteResponseOK(conn, "", "", "")

	case strings.HasPrefix(req.Path, "/echo/"):
		echoString := strings.TrimPrefix(req.Path, "/echo/")
		enc := req.Headers["Accept-Encoding"]
		WriteResponseOK(conn, echoString, "text/plain", enc)

	case strings.HasPrefix(req.Path, "/user-agent"):
		userAgent := req.Headers["User-Agent"]
		WriteResponseOK(conn, userAgent, "text/plain", "")

	case strings.HasPrefix(req.Path, "/files/"):
		filePath := filepath.Join(h.directory, strings.TrimPrefix(req.Path, "/files/"))

		content, err := os.ReadFile(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				WriteResponseNotFound(conn)
			} else {
				WriteResponseError(conn)
			}
			return
		}

		WriteResponseOK(conn, string(content), "application/octet-stream", "")

	default:
		WriteResponseNotFound(conn)
	}

}

func (h *Handler) handlePost(conn net.Conn, req *Request) {
	switch {

	case strings.HasPrefix(req.Path, "/files/"):
		filePath := filepath.Join(h.directory, strings.TrimPrefix(req.Path, "/files/"))

		err := os.WriteFile(filePath, req.Body, 0644)
		if err != nil {
			WriteResponseError(conn)
			return
		}

		WriteResponseCreated(conn)

	default:
		WriteResponseNotFound(conn)
	}
}

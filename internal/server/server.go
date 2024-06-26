package server

import (
	"fmt"
	"log"
	"net"
)

type Config struct {
	Port      int
	Directory string
}

type Server struct {
	config   Config
	listener net.Listener
	handler  *Handler
}

func New(config Config) (*Server, error) {
	handler, err := NewHandler(config.Directory)
	if err != nil {
		return nil, fmt.Errorf("failed to create handler: %w", err)
	}

	return &Server{
		config:  config,
		handler: handler,
	}, nil
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("0.0.0.0:%d", s.config.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}
	s.listener = listener

	log.Printf("Server listening on %s", addr)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return fmt.Errorf("error accepting connection: %w", err)
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) Shutdown() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	req, err := ParseRequest(conn)
	if err != nil {
		log.Printf("Error parsing request: %v", err)
		WriteResponseBad(conn)
		return
	}

	s.handler.ServeHTTP(conn, req)
}

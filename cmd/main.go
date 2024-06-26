package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/grevgeny/toy-http-server/internal/server"
)

func main() {
	// Parse command-line flag
	var (
		directory string
		port      int
	)
	flag.StringVar(&directory, "directory", "", "directory containing files")
	flag.IntVar(&port, "port", 4221, "port to listen on")
	flag.Parse()

	// Create and configure the server
	srv, err := server.New(server.Config{
		Port:      port,
		Directory: directory,
	})
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting server on port %d", port)
		if err := srv.Start(); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := srv.Shutdown(); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server stopped")
}

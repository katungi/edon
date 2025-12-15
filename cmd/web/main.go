package main

import (
	"log"
	"os"

	"github.com/katungi/edon/internal/server"
)

func main() {
	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		staticDir = "cmd/web/static"
	}
	srv, err := server.New()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	defer srv.Close()

	if err := srv.Start(staticDir); err != nil {
		log.Fatal(err)
	}
}

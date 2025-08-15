package main

import (
	"log"

	"github.com/make-bin/server-tpl/pkg/server"
)

func main() {
	srv := server.NewServer()
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

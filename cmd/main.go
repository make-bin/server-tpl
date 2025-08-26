package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/make-bin/server-tpl/pkg/server"
	"github.com/make-bin/server-tpl/pkg/utils/config"
	"github.com/make-bin/server-tpl/pkg/utils/logger"
)

func main() {
	// Initialize configuration
	cfg := config.New()

	// Initialize logger
	logger.Init(cfg.Log.Level)

	// Create server instance
	srv := server.New(cfg)

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
}

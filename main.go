package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

//go:embed frontend/dist
var frontendAssets embed.FS

func main() {
	// Initialize database
	if err := initApp(); err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}

	// Get embedded frontend filesystem (strip "frontend/dist" prefix)
	frontendFS, err := fs.Sub(frontendAssets, "frontend/dist")
	if err != nil {
		log.Fatalf("Failed to load frontend assets: %v", err)
	}

	server := newServer(frontendFS)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		fmt.Println("\nShutting down...")
		closeDB()
		os.Exit(0)
	}()

	fmt.Printf("🚀 工作安排管理系统 started at http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, server); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

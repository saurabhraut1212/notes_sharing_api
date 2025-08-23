package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/saurabhraut1212/notes_sharing_api/internal/config"
	"github.com/saurabhraut1212/notes_sharing_api/internal/db"
	"github.com/saurabhraut1212/notes_sharing_api/internal/router"
)

func main() {

	cfg := config.Load()

	client, err := db.New(cfg.MongoURI)
	if err != nil {
		log.Fatal(err)
	}

	app := router.Setup(client, cfg)

	// Channel to listen for OS signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Run server in goroutine
	go func() {
		log.Printf("🚀 Server is running successfully on http://localhost:%s", cfg.Port)
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatal("Failed to start server: ", err)
		}
	}()

	// Block until signal is received
	<-done
	log.Println("⏳ Shutting down server...")

	// Create a context with timeout to gracefully shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("✅ Server stopped gracefully")
}

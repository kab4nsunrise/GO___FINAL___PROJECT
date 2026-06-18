package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"todo/pkg/db"
	"todo/pkg/server"
)

func main() {
	srv := server.New()
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := db.Close(); err != nil {
		log.Printf("Error closing DB: %v", err)
	}
	log.Println("Server stopped")
}

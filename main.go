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

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	if err := db.Close(); err != nil {
		log.Printf("Error closing DB: %v", err)
	}
	log.Println("Server stopped")
}

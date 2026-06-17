package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"todo/pkg/api"
	"todo/pkg/db"
)

const (
	defaultPort   = "7540"
	defaultDBFile = "scheduler.db"
	webDir        = "./web"
)

type Server struct {
	port string
}

func New() *Server {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}
	return &Server{port: port}
}

func (s *Server) Start() error {

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = defaultDBFile
	}

	absDBPath, err := filepath.Abs(dbFile)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for DB: %w", err)
	}
	fmt.Printf("Database file: %s\n", absDBPath)

	if err := db.Init(dbFile); err != nil {
		return fmt.Errorf("database init: %w", err)
	}

	api.Init()

	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	addr := ":" + s.port
	fmt.Printf("Server starting on http://localhost%s\n", addr)
	return http.ListenAndServe(addr, nil)
}

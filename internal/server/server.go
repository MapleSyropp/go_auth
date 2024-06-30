package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MapleSyropp/go_auth/internal/database"
)

type Server struct {
	port int
	db   *database.Postgres
}

func NewServer(port int) *http.Server {
	NewServer := &Server{
		port: port,
		db:   database.CreateDatabase(),
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.NewRouter(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return server
}

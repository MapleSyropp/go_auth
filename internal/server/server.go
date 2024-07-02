package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MapleSyropp/go_auth/internal/database"
)

type Server struct {
	port int
	db   *database.Postgres
}

func NewServer(port int) *http.Server {
	database, err := database.CreateDatabase()
	if err != nil {
		log.Fatal("failed to create database", err)
	}
	NewServer := &Server{
		port: port,
		db:   database,
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

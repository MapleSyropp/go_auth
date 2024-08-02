package main

import (
	"fmt"

	"github.com/MapleSyropp/go_auth/internal/server"
)

func main() {
	newServer := server.NewServer(8080)
	err := newServer.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("newServer failed at start: %s", err))
	}
}

package main

import (
	"fmt"

	"github.com/MapleSyropp/go_auth/internal/server"
)

func main() {
	server := server.NewServer(8080)
	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("server failed at start: %s", err))
	}
}

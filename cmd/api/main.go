package main

import (
	"net/http"
	"github.com/MapleSyropp/go_auth/internal"
)

func main()  {
	router := routes.NewRouter()
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}


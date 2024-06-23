package routes

import (
	"fmt"
	"net/http"
)

func NewRouter() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/", helloWorld)
	router.HandleFunc("/register", register)
	router.HandleFunc("/login", login)
	router.HandleFunc("/getUser", getUser)
	router.HandleFunc("/test", test)
	return router
}

func helloWorld(w http.ResponseWriter, req *http.Request) {
	fmt.Println("inainaina")
}

func register(w http.ResponseWriter, req *http.Request) {
	fmt.Println("register")
}

func login(w http.ResponseWriter, req *http.Request) {
	fmt.Println("inainaina")
}

func getUser(w http.ResponseWriter, req *http.Request) {
	fmt.Println("inainaina")
}

func test(w http.ResponseWriter, req *http.Request) {
	fmt.Println("yea")
}

package routes

import (
	"fmt"
	"net/http"
	"github.com/rs/cors"
	"github.com/MapleSyropp/go_auth/cmd/web"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/register", register)
	mux.HandleFunc("/login-form", loginForm)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/getUser", getUser)

	router := cors.Default().Handler(mux)
	return router
}

func index(w http.ResponseWriter, req *http.Request) {
	web.Main().Render(req.Context(), w)
}

func register(w http.ResponseWriter, req *http.Request) {
}

func loginForm(w http.ResponseWriter, req *http.Request) {
	web.Login().Render(req.Context(), w)
}

func login(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	// Handle login logic here
	fmt.Fprintf(w, "Login attempted with username: %s and password: %s", username, password)
}

func getUser(w http.ResponseWriter, req *http.Request) {
}

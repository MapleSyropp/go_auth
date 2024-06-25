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
	mux.HandleFunc("/register-form", registerForm)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/login-form", loginForm)

	router := cors.Default().Handler(mux)
	return router
}

func index(w http.ResponseWriter, req *http.Request) {
	web.Main().Render(req.Context(), w)
}

func registerForm(w http.ResponseWriter, req *http.Request) {
	web.Register().Render(req.Context(), w)
}

func register(w http.ResponseWriter, req *http.Request) {
	web.Login().Render(req.Context(), w)
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
	fmt.Print(w, "Login attempted with username: %s and password: %s", username, password)
	web.Success().Render(req.Context(), w)
}

package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MapleSyropp/go_auth/cmd/web"
	"github.com/MapleSyropp/go_auth/internal/database"
	"github.com/MapleSyropp/go_auth/internal/models"
	"github.com/rs/cors"
)

func (s *Server) NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.indexHandler)
	mux.HandleFunc("/register", s.registerHandler)
	mux.HandleFunc("/register-form", s.registerFormHandler)
	mux.HandleFunc("/login", s.loginHandler)
	mux.HandleFunc("/login-form", s.loginFormHandler)
	mux.HandleFunc("/repeat", s.repeatHandler)

	router := cors.Default().Handler(mux)
	return router
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	web.Main().Render(r.Context(), w)
}

func (s *Server) repeatHandler(w http.ResponseWriter, r *http.Request) {
	web.Repeat().Render(r.Context(), w)
}

func (s *Server) registerFormHandler(w http.ResponseWriter, r *http.Request) {
	web.Register().Render(r.Context(), w)
}

func (s *Server) registerHandler(w http.ResponseWriter, r *http.Request) {
	user := new(models.UserReq)
	username := r.FormValue("username")
	password := r.FormValue("password")
	user.Username = username
	user.Password = password
	database.SaveUser(user, s.db)
	web.Login().Render(r.Context(), w)
}

func (s *Server) loginFormHandler(w http.ResponseWriter, r *http.Request) {
	web.Login().Render(r.Context(), w)
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// Handle login logic here
	fmt.Print(w, "Login attempted with username: %s and password: %s", username, password)
	web.Success().Render(r.Context(), w)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

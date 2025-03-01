package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"

	"github.com/MapleSyropp/go_auth/cmd/web"
	"github.com/MapleSyropp/go_auth/internal/database"
	"github.com/MapleSyropp/go_auth/internal/models"
)

func (s *Server) NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/client/", s.RegisterClientRoutes(mux))

	mux.HandleFunc("/", web.IndexHandler)
	mux.HandleFunc("/register", HTTPHandleFunc(s.registerHandler))
	mux.HandleFunc("/login", HTTPHandleFunc(s.loginHandler))
	mux.HandleFunc("/repeat", web.RepeatHandler)
	mux.HandleFunc("/register-form", web.RegisterFormHandler)
	mux.HandleFunc("/login-form", web.LoginFormHandler)
	mux.HandleFunc("/logout", HTTPHandleFunc(s.logoutHandler))
	mux.HandleFunc("/redirect", HTTPHandleFunc(s.redirectHandler))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:4201"},
		AllowedMethods:   []string{http.MethodGet, http.MethodOptions, http.MethodPost},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	router := c.Handler(mux)
	return router
}

func (s *Server) registerHandler(w http.ResponseWriter, r *http.Request) *models.ApiError {
	user := new(models.UserReq)

	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		return &models.ApiError{Message: "Could not read req", Code: http.StatusBadRequest}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		return &models.ApiError{Error: err, Message: "user could not be created", Code: http.StatusInternalServerError}
	}

	user.Username = username
	user.Password = string(hashedPassword)
	err = database.SaveUser(user, s.db)
	if err != nil {
		return &models.ApiError{Error: err, Message: "user could not be created", Code: http.StatusInternalServerError}
	}

	web.LoginFormHandler(w, r)
	return nil
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) *models.ApiError {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return WriteJSON(w, http.StatusMethodNotAllowed, models.ApiError{Message: fmt.Sprintf("method not allowed: %s", r.Method)})
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	user, err := database.GetUser(username, s.db)
	if err != nil {
		web.FailedHandler(w, r)
		return &models.ApiError{Error: err, Message: "", Code: http.StatusInternalServerError}
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		web.FailedHandler(w, r)
		// http.Error(w, "Invalid user", http.StatusUnauthorized)
		// return &models.ApiError{Error: err, Message: "", Code: http.StatusUnauthorized}
		return &models.ApiError{Error: err, Message: bcrypt.ErrMismatchedHashAndPassword.Error(), Code: http.StatusUnauthorized}
	}

	token, fail := buildJWT(user)
	if fail != nil {
		return &models.ApiError{Error: err, Message: bcrypt.ErrMismatchedHashAndPassword.Error(), Code: http.StatusUnauthorized}
		// return WriteJSON(w, http.StatusInternalServerError, models.ApiError{Error: bcrypt.ErrMismatchedHashAndPassword.Error()})
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(72 * time.Hour),
		HttpOnly: true,
	})

	web.SuccessHandler(w, r)
	return nil
}

func (s *Server) logoutHandler(w http.ResponseWriter, r *http.Request) *models.ApiError {
	http.SetCookie(w, &http.Cookie{
		Name:    "jwt",
		Expires: time.Unix(0, 0),
	})
	http.Redirect(w, r, "/repeat", http.StatusSeeOther)
	return nil
}

func (s *Server) redirectHandler(w http.ResponseWriter, r *http.Request) *models.ApiError {
	cookie, err := r.Cookie("redirect_uri")
	if err != nil {
		http.Error(w, "Redirect URI not found", http.StatusBadRequest)
		return nil
	}
	redirectURI, _ := url.QueryUnescape(cookie.Value)
	// redirectURI := r.URL.Query().Get("redirect_uri")

	// clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "redirect_uri",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	if r.Header.Get("HX-Request") == "true" {
		// Use HX-Redirect for HTMX requests
		w.Header().Set("HX-Redirect", redirectURI)
	} else {
		// Redirect the user back to the simple app
		http.Redirect(w, r, redirectURI, http.StatusSeeOther)
	}
	return nil
}

type apiHandler func(w http.ResponseWriter, r *http.Request) *models.ApiError

func HTTPHandleFunc(f apiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
			logger.Error(fmt.Sprintf("Error %s: %s", strconv.Itoa(err.Code), err.Message))
			http.Error(w, err.Message, err.Code)
			// WriteJSON(w, http.StatusBadRequest, models.ApiError{Error: err.Error()})
		}
	}
}

func buildJWT(user *models.User) (string, *models.ApiError) {
	privateKeyString, err := os.ReadFile("private_key.pem")
	if err != nil {
		return "", &models.ApiError{Error: err, Message: err.Error(), Code: http.StatusInternalServerError}
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyString)
	if err != nil {
		return "", &models.ApiError{Error: err, Message: err.Error(), Code: http.StatusInternalServerError}
	}

	claims := &jwt.MapClaims{
		"userID": user.ID,
		"exp":    time.Now().Add(72 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, *claims)

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", &models.ApiError{Error: err, Message: err.Error(), Code: http.StatusInternalServerError}
	}
	return signedToken, nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) *models.ApiError {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return &models.ApiError{Error: json.NewEncoder(w).Encode(v), Message: "", Code: status}
}

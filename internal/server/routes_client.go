package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/MapleSyropp/go_auth/internal/database"
	"github.com/MapleSyropp/go_auth/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) RegisterClientRoutes(mux *http.ServeMux) http.Handler {
	mux.HandleFunc("POST /client/register", HTTPHandleFunc(s.clientRegisterHandler))
	mux.HandleFunc("POST /client/login", HTTPHandleFunc(s.clientLoginHandler))
	mux.HandleFunc("/client/logout", HTTPHandleFunc(s.clientLogoutHandler))
	mux.HandleFunc("/client/redirect", HTTPHandleFunc(s.clientRedirectHandler))
	return mux
}

func (s *Server) clientRegisterHandler(w http.ResponseWriter, r *http.Request) *models.ApiError {
	user := new(models.UserReq)
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return &models.ApiError{Error: err, Message: "Could not read req", Code: http.StatusInternalServerError}
	}
	username := user.Username
	password := user.Password

	if username == "" || password == "" {
		return &models.ApiError{Message: "Could not read req", Code: http.StatusInternalServerError}
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

	return nil
}

func (s *Server) clientLoginHandler(w http.ResponseWriter, r *http.Request) *models.ApiError {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return WriteJSON(w, http.StatusMethodNotAllowed, models.ApiError{Message: fmt.Sprintf("method not allowed: %s", r.Method)})
	}

	newUser := new(models.UserReq)
	err := json.NewDecoder(r.Body).Decode(&newUser)
	user, err := database.GetUser(newUser.Username, s.db)
	if err != nil {
		return &models.ApiError{Error: err, Message: "User Not Found", Code: http.StatusNotFound}
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newUser.Password)) != nil {
		return &models.ApiError{Error: err, Message: bcrypt.ErrMismatchedHashAndPassword.Error(), Code: http.StatusUnauthorized}
	}

	token, fail := buildJWT(user)
	if fail != nil {
		return &models.ApiError{Error: err, Message: bcrypt.ErrMismatchedHashAndPassword.Error(), Code: http.StatusUnauthorized}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(72 * time.Hour),
		HttpOnly: true,
	})

	return nil
}

func (s *Server) clientLogoutHandler(w http.ResponseWriter, r *http.Request) *models.ApiError {
	http.SetCookie(w, &http.Cookie{
		Name:    "jwt",
		Expires: time.Unix(0, 0),
	})
	http.Redirect(w, r, "/repeat", http.StatusSeeOther)
	return nil
}

func (s *Server) clientRedirectHandler(w http.ResponseWriter, r *http.Request) *models.ApiError {
	cookie, err := r.Cookie("redirect_uri")
	if err != nil {
		return &models.ApiError{Error: err, Message: "Redirect URI not found", Code: http.StatusBadRequest}
	}
	redirectURI, _ := url.QueryUnescape(cookie.Value)

	http.SetCookie(w, &http.Cookie{
		Name:   "redirect_uri",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.Redirect(w, r, redirectURI, http.StatusSeeOther)
	return nil
}

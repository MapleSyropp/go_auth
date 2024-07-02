package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"

	"github.com/MapleSyropp/go_auth/cmd/web"
	"github.com/MapleSyropp/go_auth/internal/database"
	"github.com/MapleSyropp/go_auth/internal/models"
)

func (s *Server) NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", web.IndexHandler)
	mux.HandleFunc("/register", HTTPHandleFunc(s.registerHandler))
	mux.HandleFunc("/login", HTTPHandleFunc(s.loginHandler))
	mux.HandleFunc("/repeat", web.RepeatHandler)
	mux.HandleFunc("/register-form", web.RegisterFormHandler)
	mux.HandleFunc("/getData", jwtAuth(s.getData, s.db))
	mux.HandleFunc("/login-form", web.LoginFormHandler)

	router := cors.Default().Handler(mux)
	return router
}

func (s *Server) registerHandler(w http.ResponseWriter, r *http.Request) error {
	user := new(models.UserReq)
	username := r.FormValue("username")
	password := r.FormValue("password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		return err
	}
	user.Username = username
	user.Password = string(hashedPassword)
	err = database.SaveUser(user, s.db)
	savedUser, err := database.GetUser(username, s.db)
	if err != nil {
		return err
	}
	token, err := buildJWT(savedUser)
	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, models.ApiError{Error: bcrypt.ErrMismatchedHashAndPassword.Error()})
	}
	fmt.Println(token)
	web.LoginFormHandler(w, r)
	return nil
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return WriteJSON(w, http.StatusMethodNotAllowed, models.ApiError{Error: fmt.Sprintf("method not allowed: %s", r.Method)})
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	user, err := database.GetUser(username, s.db)
	if err != nil {
		return err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		web.FailedHandler(w, r)
		// http.Error(w, "Invalid user", http.StatusUnauthorized)
		return err
	}
	web.SuccessHandler(w, r)
	return nil
}

func (s *Server) getData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("llego")
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func HTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, models.ApiError{Error: err.Error()})
		}
	}
}

func jwtAuth(next http.HandlerFunc, postgres *database.Postgres) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "Authorization header is required ina", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")

		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		tokenString := bearerToken[1]

		if err := validateJWT(tokenString); err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Validate claims (userID eq token userID, exp time)

		next(w, r)
	}
}

func buildJWT(user *models.User) (string, error) {
	privateKeyString, err := os.ReadFile("private_key.pem")
	if err != nil {
		return "", err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyString))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"exp": 1500,
		"ID":  user.ID,
	})

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func validateJWT(tokenString string) error {
	publicKeyString, err := os.ReadFile("public_key.pem")
	if err != nil {
		return err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyString))
	if err != nil {
		return err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
	} else {
		return err
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

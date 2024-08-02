package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	mux.HandleFunc("/", web.IndexHandler)
	mux.HandleFunc("/register", HTTPHandleFunc(s.registerHandler))
	mux.HandleFunc("/login", HTTPHandleFunc(s.loginHandler))
	mux.HandleFunc("/repeat", web.RepeatHandler)
	mux.HandleFunc("/register-form", web.RegisterFormHandler)
	mux.HandleFunc("/getData", jwtAuth(s.getData))
	mux.HandleFunc("/login-form", web.LoginFormHandler)
	mux.HandleFunc("/logout", HTTPHandleFunc(s.logoutHandler))

	router := cors.Default().Handler(mux)
	return router
}

func (s *Server) registerHandler(w http.ResponseWriter, r *http.Request) *models.ApiError {
	user := new(models.UserReq)
	username := r.FormValue("username")
	password := r.FormValue("password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		return &models.ApiError{Error: err, Message: "user could not be created", Code: 500}
	}

	user.Username = username
	user.Password = string(hashedPassword)
	err = database.SaveUser(user, s.db)
	if err != nil {
		return &models.ApiError{Error: err, Message: "user could not be created", Code: 500}
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

func (s *Server) getData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("llego")
}

type apiHandler func(w http.ResponseWriter, r *http.Request) *models.ApiError

func HTTPHandleFunc(f apiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
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

func jwtAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		if err != nil {
			http.Error(w, "jwt cookie is required", http.StatusUnauthorized)
			return
		}
		tokenString := cookie.Value

		claims, err2 := validateJWT(tokenString)
		if err2 != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		fmt.Println(claims)
		fmt.Println(claims["userID"].(float64))
		fmt.Println(claims["exp"])
		exp, ok := claims["exp"].(float64)
		if !ok || time.Unix(int64(exp), 0).Before(time.Now()) {
			http.Error(w, "token expired", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func validateJWT(tokenString string) (jwt.MapClaims, *models.ApiError) {
	publicKeyString, err := os.ReadFile("public_key.pem")
	if err != nil {
		return nil, &models.ApiError{Error: err, Message: err.Error(), Code: http.StatusInternalServerError}
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyString)
	if err != nil {
		return nil, &models.ApiError{Error: err, Message: err.Error(), Code: http.StatusInternalServerError}
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil || !token.Valid {
		return nil, &models.ApiError{Error: err, Message: "Unauthorized", Code: http.StatusUnauthorized}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, &models.ApiError{Message: "Invalid token", Code: http.StatusUnauthorized}
	}

	return claims, nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) *models.ApiError {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return &models.ApiError{Error: json.NewEncoder(w).Encode(v), Message: "", Code: status}
}

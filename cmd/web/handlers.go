package web

import (
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	Main().Render(r.Context(), w)
}

func RepeatHandler(w http.ResponseWriter, r *http.Request) {
	Repeat().Render(r.Context(), w)
}

func RegisterFormHandler(w http.ResponseWriter, r *http.Request) {
	Register().Render(r.Context(), w)
}

func LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	Login().Render(r.Context(), w)
}

func SuccessHandler(w http.ResponseWriter, r *http.Request) {
	Success().Render(r.Context(), w)
}

func FailedHandler(w http.ResponseWriter, r *http.Request) {
	Failed().Render(r.Context(), w)
}

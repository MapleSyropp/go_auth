package web

import (
	"net/http"
	"net/url"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	redirectURI := r.URL.Query().Get("redirect_uri")

	if redirectURI != "" {
		// Set the redirect_uri in a cookie
		http.SetCookie(w, &http.Cookie{
			Name:  "redirect_uri",
			Value: url.QueryEscape(redirectURI),
			Path:  "/",
			// Optionally add more attributes, like HttpOnly, Secure, etc.
		})
	}
	Main().Render(r.Context(), w)
}

func RepeatHandler(w http.ResponseWriter, r *http.Request) {
	Repeat().Render(r.Context(), w)
}

func RegisterFormHandler(w http.ResponseWriter, r *http.Request) {
	Register().Render(r.Context(), w)
}

func LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
	Login().Render(r.Context(), w)
}

func SuccessHandler(w http.ResponseWriter, r *http.Request) {
	Success().Render(r.Context(), w)
}

func FailedHandler(w http.ResponseWriter, r *http.Request) {
	Failed().Render(r.Context(), w)
}

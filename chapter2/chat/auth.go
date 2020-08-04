package main

import (
	"net/http"
)

// The authHandler type not only implements the ServeHTTP method,
// but also stores (wraps) http.Handler in the next field
type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		// not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		// some other error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// success - call the next handler
	h.next.ServeHTTP(w, r)
}

// MustAuth helper create authHanlder that wraps any other http.Handler
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

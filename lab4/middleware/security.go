package middleware

import (
	"net/http"

	"github.com/gorilla/csrf"
)

// CSRF middleware for protecting POST forms
var csrfMiddleware = csrf.Protect([]byte("32-byte-long-auth-key"))

func CSRFMiddleware(next http.Handler) http.Handler {
	return csrfMiddleware(next)
}

// SecurityHeadersMiddleware sets HTTP security headers
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	})
}

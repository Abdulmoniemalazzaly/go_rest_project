package middlewares

import (
	"net/http"
	"slices"
)

var allowedOrigins = []string{
	"https://my-origin-url.com",
	"https://localhost:8080",
}

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if slices.Contains(allowedOrigins, origin) {
			w.Header().Set("Access-Control-Allowed-Origin", origin)
		} else {
			http.Error(w, "Not allowed by CORS", http.StatusForbidden)
		}

		w.Header().Set("Access-Control-Allowe-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Exposed-Headers", "Authorization")
		w.Header().Set("Access-Control-Allowe-Methods", "GET, POST, PUT, PATCH, DELETE")
		w.Header().Set("Access-Control-Allowe-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}

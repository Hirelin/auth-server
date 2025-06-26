package cors

import (
	"net/http"
	"slices"
	"strings"
)

var allowedOrigins = []string{}

/*
Cors.SetAllowedOrigins(origins string)
sets the allowed origins for CORS requests

@param origins: space separated string
@return: void
*/
func SetAllowedOrigins(origins string) {
	for origin := range strings.SplitSeq(origins, " ") {
		allowedOrigins = append(allowedOrigins, strings.TrimSpace(origin))
	}
}

func CorsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if slices.Contains(allowedOrigins, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

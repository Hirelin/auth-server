package adapter

import "net/http"

func AuthMailAddon() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO:addon logic to send mail on login
			next.ServeHTTP(w, r)
		})
	}
}

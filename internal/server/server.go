package server

import (
	"net/http"
	"time"

	"hirelin-auth/internal/routes"
)

func GetMux() *http.ServeMux {
	mux := http.NewServeMux()
	return mux
}

func SetFileServer(mux *http.ServeMux, path string, except string) {
	fileServer := http.FileServer(http.Dir(path))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) >= 5 && r.URL.Path[:len(except)] == except {
			http.NotFound(w, r)
			return
		}
		routes.FileServerMiddleware(http.StripPrefix("/", fileServer)).ServeHTTP(w, r)
	})
}

func CreateServer(addr string, mux *http.ServeMux) *http.Server {
	server := &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return server
}

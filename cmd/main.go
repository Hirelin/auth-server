package main

import (
	"encoding/json"
	"log"
	"net/http"

	"hirelin-auth/cmd/adapter"
	"hirelin-auth/cmd/handlers"
	"hirelin-auth/cmd/middleware"
	"hirelin-auth/cmd/oauth"
	"hirelin-auth/internal/logger"
	"hirelin-auth/internal/routes"
	"hirelin-auth/internal/server"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "pong"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
	}
}

func bindRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/ping", routes.GET(Ping))

	mux.HandleFunc("/api/auth/session", routes.GET(middleware.ProtectedMiddleware(handlers.GetSessionUser)))
	mux.HandleFunc("/api/auth/logout", routes.POST(middleware.GlobalMiddleWare(handlers.LogoutAPI)))
}

func main() {
	logger.LoadEnv()

	host := logger.GetEnv("HOST")
	port := logger.GetEnv("PORT")

	// Database
	db, err := server.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	server.SetDB(db)

	// Server
	mux := server.GetMux()

	// OAuth
	oauth.SelectProviders(adapter.OauthSqlcAdapter(), oauth.Providers.GOOGLE)
	oauth.WithOAuth(mux, routes.GET)

	// routes
	routes.BindRoutes(mux, bindRoutes)
	s := server.CreateServer(host+":"+port, mux)

	// Start server
	logger.Logger.Println(logger.Colors["green"] + "Server started on " + logger.Colors["magenta"] + host + ":" + port + logger.Colors["reset"])
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

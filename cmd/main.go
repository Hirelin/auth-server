package main

import (
	"log"
	"net/http"

	"hirelin-auth/cmd/adapter"
	"hirelin-auth/cmd/handlers"
	"hirelin-auth/cmd/middleware"
	"hirelin-auth/cmd/oauth"
	"hirelin-auth/internal/cors"
	"hirelin-auth/internal/logger"
	"hirelin-auth/internal/routes"
	"hirelin-auth/internal/server"
)

func bindRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/ping", routes.GET(handlers.Ping))

	mux.HandleFunc("/api/auth/session", routes.GET(middleware.ProtectedMiddleware(handlers.GetSessionUser)))
	mux.HandleFunc("/api/auth/logout", routes.POST(middleware.GlobalMiddleWare(handlers.LogoutAPI)))
}

func main() {
	logger.LoadEnv()

	host := logger.GetEnv("HOST")
	port := logger.GetEnv("PORT")
	allowedOrigins := logger.GetEnv("ALLOWED_ORIGINS")

	// Database
	db, err := server.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	server.SetDB(db)

	// Server
	mux := server.GetMux()
	cors.SetAllowedOrigins(allowedOrigins)

	// OAuth
	oauth.SelectProviders(adapter.OauthSqlcAdapter(), oauth.Providers.GOOGLE)
	oauth.WithOAuth(mux, routes.GET)
	oauth.WithAddons(adapter.AuthMailAddon())

	// routes
	routes.BindRoutes(mux, bindRoutes)
	s := server.CreateServer(host+":"+port, mux)

	// Start server
	logger.Logger.Println(logger.Colors["green"] + "Server started on " + logger.Colors["magenta"] + host + ":" + port + logger.Colors["reset"])
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

package server

import (
	"net/http"
	"wbrost-go/internal/handler"
	"wbrost-go/internal/middleware"
)

func SetupRoutes(authHandler *handler.AuthHandler) http.Handler {
	mux := http.NewServeMux()

	// Auth routes
	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.HandleFunc("/api/auth/signup", authHandler.Signup)
	// Protected routes (require valid token)
	mux.HandleFunc("/api/auth/me", authHandler.GetCurrentUser)
	// Apply CORS middleware
	handler := middleware.CORS(mux)

	return handler
}

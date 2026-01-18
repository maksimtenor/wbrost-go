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

	// Этот маршрут для проверки статуса API ключей
	mux.HandleFunc("/api/profile/apikeys/status", authHandler.GetApiKeysStatus)

	// Добавьте в SetupRoutes:
	mux.HandleFunc("/api/profile/update", authHandler.UpdateProfile)

	mux.HandleFunc("/api/reports", handler.GetReportsHandler)
	mux.HandleFunc("/api/reports/request", handler.RequestReportHandler)

	// Apply CORS middleware
	handler := middleware.CORS(mux)

	return handler
}

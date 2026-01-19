package server

import (
	"net/http"
	"wbrost-go/internal/handler"
	"wbrost-go/internal/middleware"
)

func SetupRoutes(authHandler *handler.AuthHandler, wbStatsHandler *handler.WBStatsHandler) http.Handler {
	mux := http.NewServeMux()

	// Auth routes
	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.HandleFunc("/api/auth/signup", authHandler.Signup)

	// Protected routes (require valid token)
	mux.HandleFunc("/api/auth/me", authHandler.GetCurrentUser)
	mux.HandleFunc("/api/profile/apikeys/status", authHandler.GetApiKeysStatus)
	mux.HandleFunc("/api/profile/update", authHandler.UpdateProfile)

	// WB Reports routes - ОБНОВЛЕНО!
	mux.HandleFunc("/api/wb/stats", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			wbStatsHandler.GetWBStats(w, r)
		case http.MethodPost:
			wbStatsHandler.CreateWBReport(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Apply CORS middleware
	handler := middleware.CORS(mux)

	return handler
}

package server

import (
	"net/http"
	"wbrost-go/internal/handler"
)

func SetupRoutes(
	authHandler *handler.AuthHandler,
	wbStatsHandler *handler.WBStatsHandler,
	wbArticlesHandler *handler.WBArticlesHandler,
) http.Handler {
	mux := http.NewServeMux()

	// Авторизационные Роуты
	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.HandleFunc("/api/auth/signup", authHandler.Signup)

	// Защищенные маршруты (требуется действительный токен)
	mux.HandleFunc("/api/auth/me", authHandler.GetCurrentUser)
	mux.HandleFunc("/api/profile/apikeys/status", authHandler.GetApiKeysStatus)
	mux.HandleFunc("/api/profile/update", authHandler.UpdateProfile)

	// Обновление пользователя из админки (заблокировать, удалить, выдать права или забрать PRO)
	mux.HandleFunc("/api/user/update", authHandler.UpdateUserParams)

	// WB-Отчеты Роуты
	mux.HandleFunc("/api/wb/stats", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			wbStatsHandler.GetWBReports(w, r)
		case http.MethodPost:
			wbStatsHandler.CreateWBReport(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Статистика Роуты
	mux.HandleFunc("/api/stat/details", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			wbStatsHandler.GetStatDetail(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Карточки товаров Роуты
	mux.HandleFunc("/api/articles", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			wbArticlesHandler.GetArticles(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/articles/request", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			wbArticlesHandler.CreateArticlesRequest(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/articles/cost-price", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			wbArticlesHandler.UpdateCostPrice(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/dashboard/stats", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			wbStatsHandler.GetDashboardStats(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/site/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			authHandler.GetUsersList(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Apply CORS middleware
	//handler := middleware.CORS(mux)

	return mux
}

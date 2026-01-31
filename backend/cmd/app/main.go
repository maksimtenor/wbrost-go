package main

import (
	"log"
	"net/http"
	"wbrost-go/internal/config"
	"wbrost-go/internal/handler"
	"wbrost-go/internal/middleware"
	"wbrost-go/internal/repository/article"
	"wbrost-go/internal/repository/database/postgres"
	"wbrost-go/internal/repository/stat"
	"wbrost-go/internal/repository/user"
	"wbrost-go/internal/server"
	"wbrost-go/internal/service/auth"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	log.Printf("Starting server with config:")
	log.Printf("  DB: %s@%s:%s/%s", cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)
	log.Printf("  Server port: %s", cfg.ServerPort)
	log.Printf("  Allowed origins: %v", cfg.AllowedOrigins)

	// Формируем строку подключения к БД
	connectionString := "host=" + cfg.DBHost +
		" port=" + cfg.DBPort +
		" user=" + cfg.DBUser +
		" password=" + cfg.DBPassword +
		" dbname=" + cfg.DBName +
		" sslmode=disable"

	// Инициализируем БД
	db, err := postgres.NewPostgresDB(connectionString)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Инициализируем репозитории
	userRepo := user.NewUserRepository(db)
	wbStatsGetRepo := stat.NewWBStatsGetRepository(db)
	statsRepo := stat.NewStatRepository(db)
	analyticsRepo := stat.NewAnalyticsRepository(db, userRepo)
	dashboardRepo := stat.NewDashboardRepository(db, userRepo)
	articlesGetRepo := article.NewWBArticlesGetRepository(db)
	articleRepo := article.NewWBArticlesRepository(db)

	// Инициализируем сервис
	authService := auth.NewAuthService(userRepo)

	// Создаем обработчики
	authHandler := handler.NewAuthHandler(authService, userRepo, cfg.JWTSecret)
	wbStatsHandler := handler.NewWBStatsHandler(userRepo, wbStatsGetRepo, statsRepo, analyticsRepo, dashboardRepo, cfg.JWTSecret)
	wbArticlesHandler := handler.NewWBArticlesHandler(userRepo, articlesGetRepo, articleRepo, cfg.JWTSecret)

	// Настраиваем маршруты
	httpHandler := server.SetupRoutes(authHandler, wbStatsHandler, wbArticlesHandler)
	// Обертываем в CORS middleware
	handlerWithCORS := middleware.CORS(cfg)(httpHandler)

	serverAddr := ":" + cfg.ServerPort
	log.Printf("Server starting on %s", serverAddr)

	if err := http.ListenAndServe(serverAddr, handlerWithCORS); err != nil {
		log.Fatal("Server failed:", err)
	}
}

package main

import (
	"log"
	"net/http"
	"wbrost-go/internal/config"
	"wbrost-go/internal/handler"
	"wbrost-go/internal/repository"
	"wbrost-go/internal/server"
	"wbrost-go/internal/service"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Формируем строку подключения к БД
	connectionString := "host=" + cfg.DBHost +
		" port=" + cfg.DBPort +
		" user=" + cfg.DBUser +
		" password=" + cfg.DBPassword +
		" dbname=" + cfg.DBName +
		" sslmode=disable"

	// Инициализируем БД
	db, err := repository.NewPostgresDB(connectionString)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Инициализируем репозиторий и сервис
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)

	// Создаем обработчики
	authHandler := handler.NewAuthHandler(authService, userRepo, cfg.JWTSecret)

	// Настраиваем маршруты
	httpHandler := server.SetupRoutes(authHandler)

	serverAddr := ":" + cfg.ServerPort
	log.Printf("Server starting on %s", serverAddr)

	// Запускаем сервер
	if err := http.ListenAndServe(serverAddr, httpHandler); err != nil {
		log.Fatal("Server failed:", err)
	}
}

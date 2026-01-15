package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	_ "os"
	"time"

	"wbrost-go/internal/config"
	"wbrost-go/internal/handler"
	"wbrost-go/internal/repository"
	"wbrost-go/internal/service"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Загружаем .env из backend/.env
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Warning: .env file not found, trying ./backend/.env")
		err = godotenv.Load("./backend/.env")
		if err != nil {
			log.Println("Warning: No .env file found, using environment variables")
		}
	}

	// Загружаем конфиг
	cfg := config.Load()

	// Подключаемся к БД
	db, err := initDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Проверяем соединение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal("Database ping failed:", err)
	}

	log.Println("✅ Connected to PostgreSQL successfully!")

	// Инициализируем репозитории
	statRepo := repository.NewStatRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Инициализируем сервисы
	siteService := service.NewSiteService(statRepo, userRepo)

	// Инициализируем хендлеры
	siteHandler := handler.NewSiteHandler(siteService)

	// Создаём роутер
	r := mux.NewRouter()

	// CORS middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// API маршруты
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/dashboard", siteHandler.Index).Methods("GET", "OPTIONS")
	api.HandleFunc("/info", siteHandler.Info).Methods("GET")
	api.HandleFunc("/privacy", siteHandler.Privacy).Methods("GET")
	api.HandleFunc("/terms", siteHandler.Terms).Methods("GET")
	api.HandleFunc("/donation", siteHandler.Donation).Methods("GET")

	// Статика для фронтенда
	fs := http.FileServer(http.Dir("../frontend/dist"))
	r.PathPrefix("/").Handler(fs)

	// Запускаем сервер
	port := cfg.ServerPort
	if port == "" {
		port = ":8080"
	}

	log.Printf("🚀 Server starting on http://localhost%s", port)
	log.Printf("📁 Serving static files from ../frontend/dist")

	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func initDatabase(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	log.Printf("📦 Connecting to PostgreSQL: %s@%s:%s/%s",
		cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)

	return sql.Open("postgres", dsn)
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Загружаем .env файл
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Параметры подключения к PostgreSQL (используем ваши переменные)
	dbHost := getEnv("DB_HOST", "")
	dbPort := getEnv("DB_PORT", "")
	dbUser := getEnv("DB_USER", "")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "")

	// Строка подключения PostgreSQL
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	fmt.Printf("Connecting to PostgreSQL: %s@%s:%s/%s\n", dbUser, dbHost, dbPort, dbName)
	fmt.Printf("Connection string: postgres://%s:****@%s:%s/%s\n", dbUser, dbHost, dbPort, dbName)

	// Подключаемся к БД
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Проверяем соединение (дольше таймаут)
	db.SetMaxOpenConns(1)
	if err := db.Ping(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	fmt.Println("✅ Connected to PostgreSQL successfully!")

	// Создаем драйвер для миграций
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("Failed to create migration driver:", err)
	}

	// Создаем объект миграции
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal("Failed to create migration instance:", err)
	}

	// Выполняем миграцию вверх
	fmt.Println("Running migrations...")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Migration failed:", err)
	}

	fmt.Println("✅ Migrations completed successfully!")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

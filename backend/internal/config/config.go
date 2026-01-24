package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	ServerPort     string
	JWTSecret      string
	Worker         WorkerConfig
	AllowedOrigins []string
}

type WorkerConfig struct {
	Interval         int // Интервал в секундах для статистики
	ArticlesInterval int // Интервал в секундах для карточек товаров
}

func Load() *Config {
	// Читаем переменные окружения
	dbPort := getEnv("DB_PORT", "")
	serverPort := getEnv("SERVER_PORT", "")

	// Определяем origins в зависимости от окружения
	allowedOrigins := []string{
		"http://localhost:3000",
		"http://localhost:3001",
		"http://localhost:8080",
		"http://localhost:8081",
	}

	// Если в окружении заданы доп origins
	if extraOrigins := os.Getenv("ALLOWED_ORIGINS"); extraOrigins != "" {
		allowedOrigins = append(allowedOrigins, extraOrigins)
	}

	return &Config{
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         dbPort,
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "123123123"),
		DBName:         getEnv("DB_NAME", "wbrost_go"),
		ServerPort:     serverPort,
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key"),
		AllowedOrigins: allowedOrigins,
		Worker: WorkerConfig{
			Interval:         getEnvAsInt("WORKER_INTERVAL", 60),
			ArticlesInterval: getEnvAsInt("WORKER_ARTICLES_INTERVAL", 60),
		},
	}
}

// GetDBConnectionString возвращает строку подключения к PostgreSQL
func (c *Config) GetDBConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

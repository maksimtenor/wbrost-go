package config

import (
	"os"
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
	Interval         int // Интервал опроса на новые события в секундах для статистики
	ArticlesInterval int // Интервал опроса на новые события в секундах для карточек товаров
}

func Load() *Config {
	dbPort := getEnv("DB_PORT", "")
	serverPort := getEnv("SERVER_PORT", "")

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

package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL драйвер
)

// PostgresDB обертка для подключения к БД
type PostgresDB struct {
	*sql.DB
}

// NewPostgresDB создает новое подключение к PostgreSQL
func NewPostgresDB(connectionString string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Устанавливаем настройки пула соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return &PostgresDB{db}, nil
}

// Close закрывает соединение с БД
func (db *PostgresDB) Close() error {
	if db.DB != nil {
		return db.DB.Close()
	}
	return nil
}

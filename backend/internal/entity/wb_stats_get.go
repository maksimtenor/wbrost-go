package entity

import (
	"database/sql"
	"time"
)

// WBStatsGet - соответствует таблице wb_stats_get
type WBStatsGet struct {
	ID        int            `json:"id" db:"id"`
	UserID    int            `json:"id_user" db:"id_user"`
	Status    sql.NullInt64  `json:"status" db:"status"`
	DateFrom  string         `json:"date_from" db:"date_from"`
	DateTo    string         `json:"date_to" db:"date_to"`
	Created   time.Time      `json:"created" db:"created"`
	Updated   time.Time      `json:"updated" db:"updated"`
	LastError sql.NullString `json:"last_error" db:"last_error"`
}

// Константы статусов
const (
	StatusWait    = 0 // В обработке
	StatusSuccess = 1 // Успешно
	StatusError   = 2 // Ошибка
)

// Константы для ошибок
const (
	TooManyRequests = "too many requests"
)

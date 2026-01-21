package entity

import (
	"database/sql"
	"time"
)

// WBArticlesGet - соответствует таблице wb_articles_get
type WBArticlesGet struct {
	ID        int            `json:"id" db:"id"`
	UserID    int            `json:"id_user" db:"id_user"`
	Status    sql.NullInt64  `json:"status" db:"status"`
	Created   time.Time      `json:"created" db:"created"`
	Updated   time.Time      `json:"updated" db:"updated"`
	LastError sql.NullString `json:"last_error" db:"last_error"`
}

// Константы статусов (аналогично статистике)
const (
	ArticlesStatusWait    = 0 // В обработке
	ArticlesStatusSuccess = 1 // Успешно
	ArticlesStatusError   = 2 // Ошибка
)

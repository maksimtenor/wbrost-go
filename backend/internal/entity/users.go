package entity

import (
	"database/sql"
	"time"
)

type Users struct {
	ID              int            `json:"id" db:"id_user"`
	Taxes           int            `json:"taxes" db:"taxes"`
	Username        string         `json:"username" db:"username"`
	PasswordHash    string         `json:"-" db:"password"`
	Email           sql.NullString `json:"email" db:"email"`
	Admin           int            `json:"admin" db:"admin"`
	Block           int            `json:"block" db:"block"`
	Pro             int            `json:"pro" db:"pro"`
	Name            sql.NullString `json:"name" db:"name"`
	Phone           sql.NullString `json:"phone" db:"phone"`
	WbKey           sql.NullString `json:"wb_key,omitempty" db:"wb_key"`
	OzonKey         sql.NullString `json:"ozon_key,omitempty" db:"ozon_key"`
	U2782212Wbrosus int            `json:"u2782212_wbrosus" db:"u2782212_wbrosus"`
	OzonStatus      int            `json:"ozon_status" db:"ozon_status"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
	Del             int            `json:"del" db:"del"`
	LastLogin       time.Time      `json:"last_login,omitempty" db:"last_login"`
}

// Константы для типов аккаунтов
const (
	UserAdmin      = 1
	UserSuperAdmin = 2
	WithoutAdmin   = 0
)

// Константы для статусов PRO
const (
	UserPro   = 1
	UserTrial = 0
)

// Константы для лимитов статистики
const (
	StatLimitFrom = "-3m"
	StatLimitTo   = "0d"
)

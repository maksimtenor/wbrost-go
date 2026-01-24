package entity

import (
	"database/sql"
	"time"
)

type User struct {
	ID              int            `json:"id" db:"id_user"`
	Taxes           int            `json:"taxes" db:"taxes"`
	Username        string         `json:"username" db:"username"`
	Password        string         `json:"-" db:"password"`
	Email           string         `json:"email" db:"email"`
	Admin           int            `json:"admin" db:"admin"`
	Block           int            `json:"block" db:"block"`
	Pro             int            `json:"pro" db:"pro"`
	Name            string         `json:"name" db:"name"`
	Phone           sql.NullString `json:"phone" db:"phone"`
	WbKey           *string        `json:"wb_key,omitempty" db:"wb_key"`
	OzonKey         *string        `json:"ozon_key,omitempty" db:"ozon_key"`
	U2782212Wbrosus int            `json:"u2782212_wbrosus" db:"u2782212_wbrosus"`
	OzonStatus      int            `json:"ozon_status" db:"ozon_status"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
	Del             int            `json:"del" db:"del"`
	LastLogin       *time.Time     `json:"last_login,omitempty" db:"last_login"`
}

// Константы для типов аккаунтов
const (
	UserAdmin      = 2
	UserSuperAdmin = 1
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

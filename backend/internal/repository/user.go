package repository

import (
	"database/sql"
	"fmt"
	"time"
)

// User - модель для работы с таблицей users
type User struct {
	ID              int            `db:"id_user"`
	Taxes           int            `db:"taxes"`
	Username        string         `db:"username"`
	PasswordHash    string         `db:"password"` // В таблице поле password
	Email           sql.NullString `db:"email"`    // Email может быть NULL
	Admin           int            `db:"admin"`
	Block           int            `db:"block"`
	Pro             int            `db:"pro"`  // В таблице поле pro (int), а не pro_account (string)
	Name            sql.NullString `db:"name"` // Name может быть NULL
	Phone           sql.NullString `db:"phone"`
	WbKey           sql.NullString `db:"wb_key"`
	OzonKey         sql.NullString `db:"ozon_key"`
	U2782212Wbrosus int            `db:"u2782212_wbrosus"`
	OzonStatus      int            `db:"ozon_status"`
	CreatedAt       time.Time      `db:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at"`
	Del             int            `db:"del"`
	LastLogin       time.Time      `db:"last_login"`
}

type UserRepository struct {
	db *PostgresDB
}

func NewUserRepository(db *PostgresDB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByUsername(username string) (*User, error) {
	query := `SELECT 
        id_user, taxes, username, password, email, admin, block, pro, 
        name, phone, wb_key, ozon_key, u2782212_wbrosus, ozon_status,
        created_at, updated_at, del, last_login
        FROM users 
        WHERE username = $1 AND del = 0`

	row := r.db.QueryRow(query, username)

	var user User

	err := row.Scan(
		&user.ID,
		&user.Taxes,
		&user.Username,
		&user.PasswordHash,
		&user.Email,
		&user.Admin,
		&user.Block,
		&user.Pro,
		&user.Name,
		&user.Phone,
		&user.WbKey,
		&user.OzonKey,
		&user.U2782212Wbrosus,
		&user.OzonStatus,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Del,
		&user.LastLogin,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*User, error) {
	query := `SELECT 
        id_user, taxes, username, password, email, admin, block, pro, 
        name, phone, wb_key, ozon_key, u2782212_wbrosus, ozon_status,
        created_at, updated_at, del, last_login
        FROM users 
        WHERE email = $1 AND del = 0`

	row := r.db.QueryRow(query, email)

	var user User

	err := row.Scan(
		&user.ID,
		&user.Taxes,
		&user.Username,
		&user.PasswordHash,
		&user.Email,
		&user.Admin,
		&user.Block,
		&user.Pro,
		&user.Name,
		&user.Phone,
		&user.WbKey,
		&user.OzonKey,
		&user.U2782212Wbrosus,
		&user.OzonStatus,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Del,
		&user.LastLogin,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) Create(user *User) (*User, error) {
	// Проверяем существование username
	existingUser, _ := r.GetByUsername(user.Username)
	if existingUser != nil {
		return nil, fmt.Errorf("username already exists")
	}

	// Проверяем существование email если email предоставлен
	if user.Email.Valid && user.Email.String != "" {
		existingEmail, _ := r.GetByEmail(user.Email.String)
		if existingEmail != nil {
			return nil, fmt.Errorf("email already exists")
		}
	}

	// SQL запрос должен соответствовать структуре таблицы
	query := `INSERT INTO users (
        username, password, email, name, pro, admin, 
        taxes, block, phone, wb_key, ozon_key, 
        u2782212_wbrosus, ozon_status, del, last_login
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) 
    RETURNING id_user, created_at, updated_at`

	// Устанавливаем значения по умолчанию
	if !user.Email.Valid {
		user.Email = sql.NullString{String: "", Valid: false}
	}
	if !user.Name.Valid {
		user.Name = sql.NullString{String: "", Valid: false}
	}
	if !user.Phone.Valid {
		user.Phone = sql.NullString{String: "", Valid: false}
	}
	if !user.WbKey.Valid {
		user.WbKey = sql.NullString{String: "", Valid: false}
	}
	if !user.OzonKey.Valid {
		user.OzonKey = sql.NullString{String: "", Valid: false}
	}

	err := r.db.QueryRow(
		query,
		user.Username,
		user.PasswordHash,
		user.Email,
		user.Name,
		user.Pro, // pro вместо pro_account
		user.Admin,
		user.Taxes,
		user.Block,
		user.Phone,
		user.WbKey,
		user.OzonKey,
		user.U2782212Wbrosus,
		user.OzonStatus,
		user.Del,
		user.LastLogin,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

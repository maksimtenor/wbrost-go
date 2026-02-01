package user

import (
	"database/sql"
	"fmt"
	"time"
	"wbrost-go/internal/entity"
	"wbrost-go/internal/repository/database/postgres"
)

type UserRepository struct {
	db *postgres.PostgresDB
}

func NewUserRepository(db *postgres.PostgresDB) *UserRepository {
	return &UserRepository{db: db}
}
func (r *UserRepository) GetCountAllUsers() (int, error) {
	var count int
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM users",
	).Scan(&count)

	return count, err
}
func (r *UserRepository) GetAll(page, pageSize int) ([]entity.Users, error) {
	offset := (page - 1) * pageSize
	query := `SELECT 
        id_user, taxes, username, password, email, admin, block, pro, 
        name, phone, wb_key, ozon_key, u2782212_wbrosus, ozon_status,
        created_at, updated_at, del, last_login
        FROM users where del = 0
ORDER BY created_at DESC
LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []entity.Users
	for rows.Next() {
		var a entity.Users
		err := rows.Scan(
			&a.ID,
			&a.Taxes,
			&a.Username,
			&a.PasswordHash,
			&a.Email,
			&a.Admin,
			&a.Block,
			&a.Pro,
			&a.Name,
			&a.Phone,
			&a.WbKey,
			&a.OzonKey,
			&a.U2782212Wbrosus,
			&a.OzonStatus,
			&a.CreatedAt,
			&a.UpdatedAt,
			&a.Del,
			&a.LastLogin,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, a)
	}

	return users, nil
}

func (r *UserRepository) GetByUsername(username string) (*entity.Users, error) {
	query := `SELECT 
        id_user, taxes, username, password, email, admin, block, pro, 
        name, phone, wb_key, ozon_key, u2782212_wbrosus, ozon_status,
        created_at, updated_at, del, last_login
        FROM users 
        WHERE username = $1 AND del = 0`

	row := r.db.QueryRow(query, username)

	var user entity.Users

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

func (r *UserRepository) GetByUserId(userId int) (*entity.Users, error) {
	query := `SELECT 
        id_user, taxes, username, password, email, admin, block, pro, 
        name, phone, wb_key, ozon_key, u2782212_wbrosus, ozon_status,
        created_at, updated_at, del, last_login
        FROM users 
        WHERE id = $1 AND del = 0`

	row := r.db.QueryRow(query, userId)

	var user entity.Users

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

// UpdateUserPro обновляет пользователя PRO
func (r *UserRepository) UpdateUserPro(userId, value int) error {
	query := `
		UPDATE users 
		SET pro = $1, updated_at = $2
		WHERE id_user = $3
	`

	_, err := r.db.Exec(query, value, time.Now(), userId)

	return err
}

// UpdateUserAdmin обновляет пользователя Admin
func (r *UserRepository) UpdateUserAdmin(userId, value int) error {
	query := `
		UPDATE users 
		SET admin = $1, updated_at = $2
		WHERE id_user = $3
	`

	_, err := r.db.Exec(query, value, time.Now(), userId)

	return err
}

// UpdateUserBlock обновляет пользователя BLOCK
func (r *UserRepository) UpdateUserBlock(userId, value int) error {
	query := `
		UPDATE users 
		SET block = $1, updated_at = $2
		WHERE id_user = $3
	`

	_, err := r.db.Exec(query, value, time.Now(), userId)

	return err
}

// UpdateUserDel обновляет пользователя DEL
func (r *UserRepository) UpdateUserDel(userId, value int) error {
	query := `
		UPDATE users 
		SET del = $1, updated_at = $2
		WHERE id_user = $3
	`

	_, err := r.db.Exec(query, value, time.Now(), userId)

	return err
}

func (r *UserRepository) GetByEmail(email string) (*entity.Users, error) {
	query := `SELECT 
        id_user, taxes, username, password, email, admin, block, pro, 
        name, phone, wb_key, ozon_key, u2782212_wbrosus, ozon_status,
        created_at, updated_at, del, last_login
        FROM users 
        WHERE email = $1 AND del = 0`

	row := r.db.QueryRow(query, email)

	var user entity.Users

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

func (r *UserRepository) Create(user *entity.Users) (*entity.Users, error) {
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

func (r *UserRepository) GetByID(id int) (*entity.Users, error) {
	query := `SELECT * FROM users WHERE id_user = $1`
	row := r.db.QueryRow(query, id)

	var user entity.Users
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
		return nil, err
	}
	return &user, nil
}
func (r *UserRepository) UpdateUser(user *entity.Users) error {
	query := `
        UPDATE users 
        SET 
            name = $1,
            email = $2,
            phone = $3,
            taxes = $4,
            wb_key = $5,
            password = $6,
            updated_at = CURRENT_TIMESTAMP
        WHERE id_user = $7
    `

	// Подготавливаем значения для NULL полей
	nameValue := getNullStringValue(user.Name)
	emailValue := getNullStringValue(user.Email)
	phoneValue := getNullStringValue(user.Phone)
	wbKeyValue := getNullStringValue(user.WbKey)

	// Выполняем запрос
	_, err := r.db.Exec(query,
		nameValue,
		emailValue,
		phoneValue,
		user.Taxes,
		wbKeyValue,
		user.PasswordHash,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Вспомогательная функция для работы с NULL
func getNullStringValue(ns sql.NullString) interface{} {
	if ns.Valid {
		return ns.String
	}
	return nil
}

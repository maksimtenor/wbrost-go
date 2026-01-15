package repository

import (
	"context"
	"database/sql"
	"time"
	"wbrost-go/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id int) (*entity.User, error)
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]*entity.User, error)
	Count(ctx context.Context) (int, error)
	CountRealUsers(ctx context.Context) (int, error)
	UpdateLastLogin(ctx context.Context, id int) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (
			taxes, username, password, email, admin, block, pro, 
			name, phone, wb_key, ozon_key, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		user.Taxes,
		user.Username,
		user.Password,
		user.Email,
		user.Admin,
		user.Block,
		user.Pro,
		user.Name,
		user.Phone,
		user.WbKey,
		user.OzonKey,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

func (r *userRepository) FindByID(ctx context.Context, id int) (*entity.User, error) {
	query := `
		SELECT * FROM users 
		WHERE id_user = ? AND del = 0
	`

	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanUser(row)
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	query := `
		SELECT * FROM users 
		WHERE username = ? AND del = 0
	`

	row := r.db.QueryRowContext(ctx, query, username)
	return r.scanUser(row)
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users SET
			taxes = ?, username = ?, email = ?, admin = ?, block = ?, pro = ?,
			name = ?, phone = ?, wb_key = ?, ozon_key = ?, updated_at = ?
		WHERE id_user = ?
	`

	user.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		user.Taxes,
		user.Username,
		user.Email,
		user.Admin,
		user.Block,
		user.Pro,
		user.Name,
		user.Phone,
		user.WbKey,
		user.OzonKey,
		user.UpdatedAt,
		user.ID,
	)

	return err
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE users SET del = 1, updated_at = ? WHERE id_user = ?`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	query := `
		SELECT * FROM users 
		WHERE del = 0 
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		user, err := r.scanUserFromRows(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM users WHERE del = 0`
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (r *userRepository) CountRealUsers(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM users WHERE admin = 0 AND del = 0`
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, id int) error {
	query := `UPDATE users SET last_login = ? WHERE id_user = ?`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (r *userRepository) scanUser(row *sql.Row) (*entity.User, error) {
	var user entity.User
	err := row.Scan(
		&user.ID,
		&user.Taxes,
		&user.Username,
		&user.Password,
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

func (r *userRepository) scanUserFromRows(rows *sql.Rows) (*entity.User, error) {
	var user entity.User
	err := rows.Scan(
		&user.ID,
		&user.Taxes,
		&user.Username,
		&user.Password,
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

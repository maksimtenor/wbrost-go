package repository

import (
	"database/sql"
	"time"
)

// WBStatsGet - модель таблицы wb_stats_get
type WBStatsGet struct {
	ID        int            `db:"id"`
	UserID    int            `db:"id_user"`
	Status    sql.NullInt64  `db:"status"` // 0=в обработке, 1=готово, 2=ошибка
	DateFrom  string         `db:"date_from"`
	DateTo    string         `db:"date_to"`
	Created   time.Time      `db:"created"`
	Updated   time.Time      `db:"updated"`
	LastError sql.NullString `db:"last_error"`
}

type WBStatsRepository struct {
	db *PostgresDB
}

func NewWBStatsRepository(db *PostgresDB) *WBStatsRepository {
	return &WBStatsRepository{db: db}
}

// Create создает новую запись отчета
func (r *WBStatsRepository) Create(stats *WBStatsGet) error {
	query := `
		INSERT INTO wb_stats_get (id_user, status, date_from, date_to, last_error)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created, updated
	`

	return r.db.QueryRow(query,
		stats.UserID,
		stats.Status,
		stats.DateFrom,
		stats.DateTo,
		stats.LastError,
	).Scan(&stats.ID, &stats.Created, &stats.Updated)
}

// GetByUserID получает отчеты пользователя
func (r *WBStatsRepository) GetByUserID(userID int) ([]WBStatsGet, error) {
	query := `
		SELECT id, id_user, status, date_from, date_to, created, updated, last_error
		FROM wb_stats_get 
		WHERE id_user = $1 
		ORDER BY created DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []WBStatsGet
	for rows.Next() {
		var s WBStatsGet
		err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.Status,
			&s.DateFrom,
			&s.DateTo,
			&s.Created,
			&s.Updated,
			&s.LastError,
		)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}

	return stats, nil
}

// UpdateStatus обновляет статус отчета
func (r *WBStatsRepository) UpdateStatus(id int, status int, errorMsg string) error {
	query := `
		UPDATE wb_stats_get 
		SET status = $1, last_error = $2, updated = CURRENT_TIMESTAMP
		WHERE id = $3
	`

	_, err := r.db.Exec(query, status, errorMsg, id)
	return err
}

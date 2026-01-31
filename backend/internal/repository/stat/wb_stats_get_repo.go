package stat

import (
	"fmt"
	"time"
	"wbrost-go/internal/entity"
	"wbrost-go/internal/repository/database/postgres"
)

type WBStatsGetRepository struct {
	db *postgres.PostgresDB
}

func NewWBStatsGetRepository(db *postgres.PostgresDB) *WBStatsGetRepository {
	return &WBStatsGetRepository{db: db}
}

// Create создает новую запись отчета
func (r *WBStatsGetRepository) Create(stats *entity.WBStatsGet) error {
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
func (r *WBStatsGetRepository) GetByUserID(userID int) ([]entity.WBStatsGet, error) {
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

	var stats []entity.WBStatsGet
	for rows.Next() {
		var s entity.WBStatsGet
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
func (r *WBStatsGetRepository) UpdateStatus(id int, status int, errorMsg string) error {
	query := `
		UPDATE wb_stats_get 
		SET status = $1, last_error = $2, updated = CURRENT_TIMESTAMP
		WHERE id = $3
	`

	_, err := r.db.Exec(query, status, errorMsg, id)
	return err
}

// GetPendingOrders возвращает заказы с статусом 0 (в обработке)
func (r *WBStatsGetRepository) GetPendingOrders() ([]entity.WBStatsGet, error) {
	query := `
		SELECT id, id_user, status, date_from, date_to, created, updated, last_error
		FROM wb_stats_get 
		WHERE status = $1 
		ORDER BY created ASC
	`

	rows, err := r.db.Query(query, entity.StatusWait)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending orders: %w", err)
	}
	defer rows.Close()

	var orders []entity.WBStatsGet
	for rows.Next() {
		var order entity.WBStatsGet
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Status,
			&order.DateFrom,
			&order.DateTo,
			&order.Created,
			&order.Updated,
			&order.LastError,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// UpdateOrderStatus обновляет статус заказа
func (r *WBStatsGetRepository) UpdateOrderStatus(orderID int, status int, errorMsg string) error {
	query := `
		UPDATE wb_stats_get 
		SET status = $1, last_error = $2, updated = $3
		WHERE id = $4
	`

	_, err := r.db.Exec(query, status, errorMsg, time.Now(), orderID)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}

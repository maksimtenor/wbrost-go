package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"wbrost-go/internal/entity"
)

type StatRepository interface {
	GetDashboardData(ctx context.Context, userID int, dateFrom, dateTo string) ([]entity.StatSummary, error)
	GetTotal(ctx context.Context, userID int, dateFrom, dateTo string, field string) (float64, error)
	GetTotalClear(ctx context.Context, userID int, dateFrom, dateTo string) (float64, error)
	GetMonthlyStats(ctx context.Context, userID int, monthsBack int) ([]entity.StatSummary, error)
}

type statRepository struct {
	db *sql.DB
}

func NewStatRepository(db *sql.DB) StatRepository {
	return &statRepository{db: db}
}

func (r *statRepository) GetDashboardData(ctx context.Context, userID int, dateFrom, dateTo string) ([]entity.StatSummary, error) {
	query := `
		SELECT 
			TO_CHAR(rr_dt, 'YYYY-MM') as month,
			SUM(CAST(ppvz_for_pay AS DECIMAL)) as total_sell,
			SUM(CASE WHEN supplier_oper_name = $1 THEN 1 ELSE 0 END) as total_count_sales,
			SUM(CASE WHEN supplier_oper_name = $2 THEN 1 ELSE 0 END) as total_refund
		FROM stat 
		WHERE user_id = $3 
			AND rr_dt >= $4 
			AND rr_dt <= $5
		GROUP BY TO_CHAR(rr_dt, 'YYYY-MM')
		ORDER BY month
	`

	rows, err := r.db.QueryContext(ctx, query,
		entity.TypeSale,
		entity.TypeRefund,
		userID,
		dateFrom,
		dateTo,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []entity.StatSummary
	for rows.Next() {
		var summary entity.StatSummary
		err := rows.Scan(
			&summary.Month,
			&summary.TotalSell,
			&summary.TotalCountSales,
			&summary.TotalRefund,
		)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

func (r *statRepository) GetTotal(ctx context.Context, userID int, dateFrom, dateTo string, field string) (float64, error) {
	validFields := map[string]bool{
		"ppvz_for_pay": true,
		"count_sales":  true,
		"count_refund": true,
	}

	if !validFields[field] {
		return 0, fmt.Errorf("invalid field: %s", field)
	}

	var query string
	if field == "count_sales" {
		query = `
			SELECT COUNT(*) as total
			FROM stat 
			WHERE user_id = $1 
				AND rr_dt >= $2 
				AND rr_dt <= $3
				AND supplier_oper_name = $4
		`
	} else if field == "count_refund" {
		query = `
			SELECT COUNT(*) as total
			FROM stat 
			WHERE user_id = $1 
				AND rr_dt >= $2 
				AND rr_dt <= $3
				AND supplier_oper_name = $4
		`
	} else {
		query = `
			SELECT COALESCE(SUM(CAST($4 AS DECIMAL)), 0) as total
			FROM stat 
			WHERE user_id = $1 
				AND rr_dt >= $2 
				AND rr_dt <= $3
		`
	}

	var total float64
	var err error

	if field == "count_sales" {
		err = r.db.QueryRowContext(ctx, query, userID, dateFrom, dateTo, entity.TypeSale).Scan(&total)
	} else if field == "count_refund" {
		err = r.db.QueryRowContext(ctx, query, userID, dateFrom, dateTo, entity.TypeRefund).Scan(&total)
	} else {
		err = r.db.QueryRowContext(ctx, query, userID, dateFrom, dateTo, field).Scan(&total)
	}

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	return total, nil
}

func (r *statRepository) GetTotalClear(ctx context.Context, userID int, dateFrom, dateTo string) (float64, error) {
	// Здесь будет сложная логика расчёта чистой прибыли
	// Пока возвращаем заглушку
	return 0, nil
}

func (r *statRepository) GetMonthlyStats(ctx context.Context, userID int, monthsBack int) ([]entity.StatSummary, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, -monthsBack, 0)

	return r.GetDashboardData(ctx, userID,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
	)
}

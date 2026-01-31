package stat

import (
	"database/sql"
	"fmt"
	"time"
	"wbrost-go/internal/repository/database/postgres"
	"wbrost-go/internal/repository/user"
)

type DashboardRepository struct {
	db       *postgres.PostgresDB
	userRepo *user.UserRepository
}

func NewDashboardRepository(db *postgres.PostgresDB, userRepo *user.UserRepository) *DashboardRepository {
	return &DashboardRepository{db: db, userRepo: userRepo}
}

// GetDashboardStats получает статистику для дашборда
func (r *DashboardRepository) GetDashboardStats(userID int, dateFrom, dateTo string) (map[string]interface{}, error) {
	dateToWithTime := dateTo + " 23:59:59"

	query := `
        SELECT 
            -- Количество продаж
            COUNT(DISTINCT 
                CASE 
                    WHEN s.supplier_oper_name IN (1, 7) AND s.quantity > 0 
                    THEN s.shk_id 
                END
            ) as sales_count,
            
            -- Сумма к перечислению
            SUM(
                CASE 
                    WHEN s.ppvz_for_pay ~ '^[0-9]+\.?[0-9]*$' 
                    THEN CAST(s.ppvz_for_pay AS NUMERIC) 
                    ELSE 0 
                END
            ) as ppvz_for_pay_total,
            
            -- Возвраты
            COUNT(DISTINCT 
                CASE 
                    WHEN s.supplier_oper_name = 2 AND s.return_amount > 0 
                    THEN s.shk_id 
                END
            ) as returns_count,
            
            -- Чистая прибыль
            SUM(
                CASE 
                    WHEN s.ppvz_for_pay ~ '^[0-9]+\.?[0-9]*$' 
                    THEN CAST(s.ppvz_for_pay AS NUMERIC) 
                    ELSE 0 
                END
            ) - 
            SUM(COALESCE(s.delivery_rub, 0)) - 
            SUM(COALESCE(s.deduction, 0)) - 
            SUM(COALESCE(s.storage_fee, 0)) - 
            SUM(COALESCE(s.additional_payment, 0)) - 
            SUM(COALESCE(s.penalty, 0)) as net_profit
        FROM wb_stats s
        WHERE s.user_id = $1
            AND s.sale_dt BETWEEN $2 AND $3
    `

	var salesCount, returnsCount sql.NullInt64
	var ppvzForPayTotal, netProfit sql.NullFloat64

	err := r.db.QueryRow(query, userID, dateFrom, dateToWithTime).Scan(
		&salesCount,
		&ppvzForPayTotal,
		&returnsCount,
		&netProfit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard stats: %w", err)
	}

	// Форматируем значения
	stats := map[string]interface{}{
		"sales_count":        getIntValue(salesCount),
		"ppvz_for_pay_total": getFloatValue(ppvzForPayTotal),
		"returns_count":      getIntValue(returnsCount),
		"net_profit":         getFloatValue(netProfit),
	}

	return stats, nil
}

// GetChartData получает данные для графиков
func (r *DashboardRepository) GetChartData(userID int, dateFrom, dateTo string) (map[string]interface{}, error) {
	dateToWithTime := dateTo + " 23:59:59"

	// Данные для линейного графика (продажи по дням)
	queryLineChart := `
        SELECT 
            DATE(s.sale_dt) as sale_date,
            SUM(
                CASE 
                    WHEN s.ppvz_for_pay ~ '^[0-9]+\.?[0-9]*$' 
                    THEN CAST(s.ppvz_for_pay AS NUMERIC) 
                    ELSE 0 
                END
            ) as daily_revenue,
            COUNT(DISTINCT 
                CASE 
                    WHEN s.supplier_oper_name IN (1, 7) 
                    THEN s.shk_id 
                END
            ) as daily_sales
        FROM wb_stats s
        WHERE s.user_id = $1
            AND s.sale_dt BETWEEN $2 AND $3
            AND s.sale_dt IS NOT NULL
        GROUP BY DATE(s.sale_dt)
        ORDER BY sale_date
    `

	rows, err := r.db.Query(queryLineChart, userID, dateFrom, dateToWithTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get line chart data: %w", err)
	}
	defer rows.Close()

	lineChartLabels := []string{}
	lineChartRevenue := []float64{}
	lineChartSales := []float64{}

	for rows.Next() {
		var saleDate time.Time
		var dailyRevenue, dailySales sql.NullFloat64

		err := rows.Scan(&saleDate, &dailyRevenue, &dailySales)
		if err != nil {
			return nil, fmt.Errorf("failed to scan line chart data: %w", err)
		}

		lineChartLabels = append(lineChartLabels, saleDate.Format("02 Jan"))
		lineChartRevenue = append(lineChartRevenue, getFloatValue(dailyRevenue))
		lineChartSales = append(lineChartSales, getFloatValue(dailySales))
	}

	// Данные для столбчатого графика (топ категорий)
	queryBarChart := `
        SELECT 
            COALESCE(s.subject_name, 'Без категории') as category,
            COUNT(DISTINCT s.nm_id) as product_count,
            SUM(
                CASE 
                    WHEN s.ppvz_for_pay ~ '^[0-9]+\.?[0-9]*$' 
                    THEN CAST(s.ppvz_for_pay AS NUMERIC) 
                    ELSE 0 
                END
            ) as category_revenue
        FROM wb_stats s
        WHERE s.user_id = $1
            AND s.sale_dt BETWEEN $2 AND $3
            AND s.subject_name IS NOT NULL
        GROUP BY s.subject_name
        ORDER BY category_revenue DESC
        LIMIT 5
    `

	rowsBar, err := r.db.Query(queryBarChart, userID, dateFrom, dateToWithTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get bar chart data: %w", err)
	}
	defer rowsBar.Close()

	barChartLabels := []string{}
	barChartData := []float64{}

	for rowsBar.Next() {
		var category string
		var productCount sql.NullInt64
		var categoryRevenue sql.NullFloat64

		err := rowsBar.Scan(&category, &productCount, &categoryRevenue)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bar chart data: %w", err)
		}

		barChartLabels = append(barChartLabels, category)
		barChartData = append(barChartData, getFloatValue(categoryRevenue))
	}

	chartData := map[string]interface{}{
		"line_chart": map[string]interface{}{
			"labels":  lineChartLabels,
			"revenue": lineChartRevenue,
			"sales":   lineChartSales,
		},
		"bar_chart": map[string]interface{}{
			"labels": barChartLabels,
			"data":   barChartData,
		},
	}

	return chartData, nil
}

// GetMonthlyRevenue получает выручку по месяцам за последние 12 месяцев
func (r *DashboardRepository) GetMonthlyRevenue(userID int) (map[string]interface{}, error) {
	// Получаем текущую дату
	now := time.Now()

	// Начало - первый день месяца 11 месяцев назад
	startDate := time.Date(now.Year(), now.Month()-11, 1, 0, 0, 0, 0, now.Location())
	// Конец - последний день текущего месяца
	endDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// Запрос для получения данных
	query := `
        WITH months AS (
            SELECT generate_series(
                date_trunc('month', $2::timestamp),
                date_trunc('month', $3::timestamp),
                '1 month'::interval
            ) as month_start
        )
        SELECT 
            TO_CHAR(m.month_start, 'Mon YYYY') as month_name,
            EXTRACT(YEAR FROM m.month_start) as year,
            COALESCE(SUM(
                CASE 
                    WHEN s.ppvz_for_pay ~ '^[0-9]+\.?[0-9]*$' 
                    THEN CAST(s.ppvz_for_pay AS NUMERIC) 
                    ELSE 0 
                END
            ), 0) as monthly_revenue
        FROM months m
        LEFT JOIN wb_stats s ON 
            date_trunc('month', s.sale_dt) = m.month_start 
            AND s.user_id = $1
        GROUP BY m.month_start, year
        ORDER BY m.month_start
    `

	rows, err := r.db.Query(query, userID,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"))

	if err != nil {
		return nil, fmt.Errorf("failed to get monthly revenue: %w", err)
	}
	defer rows.Close()

	monthlyLabels := []string{}
	monthlyData := []float64{}

	for rows.Next() {
		var monthName string
		var year int
		var monthlyRevenue float64

		err := rows.Scan(&monthName, &year, &monthlyRevenue)
		if err != nil {
			return nil, fmt.Errorf("failed to scan monthly revenue: %w", err)
		}

		// Форматируем метку: "Фев 2025"
		monthLabel := formatMonthLabel(monthName, year)
		monthlyLabels = append(monthlyLabels, monthLabel)
		monthlyData = append(monthlyData, monthlyRevenue)

		// Отладка
		//fmt.Printf("Month: %s, Year: %d, Revenue: %.2f, Label: %s\n",
		//	monthName, year, monthlyRevenue, monthLabel)
	}

	// Если нет данных, создаем метки для последних 12 месяцев
	if len(monthlyLabels) == 0 {
		for i := 0; i < 12; i++ {
			date := startDate.AddDate(0, i, 0)
			monthName := date.Format("Jan")
			year := date.Year()
			monthLabel := formatMonthLabel(monthName, year)
			monthlyLabels = append(monthlyLabels, monthLabel)
			monthlyData = append(monthlyData, 0)
		}
	}

	return map[string]interface{}{
		"labels": monthlyLabels,
		"data":   monthlyData,
	}, nil
}

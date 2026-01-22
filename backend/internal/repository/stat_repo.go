package repository

import (
	"database/sql"
	"fmt"
	"time"
	"wbrost-go/internal/entity"
)

type StatRepository struct {
	db       *PostgresDB
	userRepo *UserRepository
}

type StatRepositoryInterface interface {
	Create(stat *entity.Stat) error
	GetStatDetails(userID int, dateFrom, dateTo string, page, pageSize int) ([]map[string]interface{}, error)
	GetStatDetailsCount(userID int, dateFrom, dateTo string) (int, error)
	GetStatSummary(userID int, dateFrom, dateTo string) (map[string]interface{}, error)
	ExistsByHash(hash string) (bool, error)
}

func NewStatRepository(db *PostgresDB, userRepo *UserRepository) *StatRepository {
	return &StatRepository{db: db, userRepo: userRepo}
}

// Create создает новую запись статистики в таблицу stat
func (r *StatRepository) Create(stat *entity.Stat) error {
	query := `
		INSERT INTO stat (
			hash_info, user_id, nm_id, ppvz_for_pay, supplier_oper_name,
			delivery_rub, penalty, additional_payment, storage_fee,
			rebill_logistic_cost, acquiring_fee, acquiring_percent,
			ppvz_sales_commission, deduction, ppvz_spp_prc, ppvz_kvw_prc_base,
			ppvz_kvw_prc, acceptance, dlv_prc, created_at, rr_dt,
			shk_id, sticker_id, gi_id, realizationreport_id, barcode,
			bonus_type_name, last_error, brand_name, ppvz_office_id,
			assembly_id, sa_name, ppvz_vw_nds, ppvz_vw, gi_box_type_name,
			subject_name, ts_name, quantity, retail_price, retail_amount,
			commission_percent, office_name, order_dt, sale_dt,
			delivery_amount, return_amount, report_type, srid, rid
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
			$21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
			$31, $32, $33, $34, $35, $36, $37, $38, $39, $40,
			$41, $42, $43, $44, $45, $46, $47, $48, $49
		)
		RETURNING id
	`

	err := r.db.QueryRow(query,
		// 1-10
		stat.HashInfo, // <-- HashInfo это string, передаем напрямую
		stat.UserID,
		getNullInt64(stat.Nmid),
		getNullString(stat.PpvzForPay),
		getNullInt64(stat.SupplierOperName),
		getNullFloat64(stat.DeliveryRub),
		getNullFloat64(stat.Penalty),
		getNullFloat64(stat.AdditionalPayment),
		getNullFloat64(stat.StorageFee),
		getNullString(stat.RebillLogisticCost),
		// 11-20
		getNullFloat64(stat.AcquiringFee),
		getNullFloat64(stat.AcquiringPercent),
		getNullFloat64(stat.PpvzSalesCommission),
		getNullFloat64(stat.Deduction),
		getNullString(stat.PpvzSppPrc),
		getNullString(stat.PpvzKvwPrcBase),
		getNullString(stat.PpvzKvwPrc),
		getNullFloat64(stat.Acceptance),
		getNullFloat64(stat.DlvPrc),
		stat.CreatedAt,
		// 21-30
		getNullTime(stat.RrDt),
		getNullInt64(stat.ShkID),
		getNullString(stat.StickerID),
		getNullInt64(stat.GiID),
		getNullInt64(stat.RealizationreportID),
		getNullString(stat.Barcode),
		getNullString(stat.BonusTypeName),
		getNullString(stat.LastError),
		getNullString(stat.BrandName),
		getNullInt64(stat.PpvzOfficeID),
		// 31-40
		getNullInt64(stat.AssemblyID),
		getNullString(stat.SaName),
		getNullString(stat.PpvzVwNds),
		getNullString(stat.PpvzVw),
		getNullString(stat.GiBoxTypeName),
		getNullString(stat.SubjectName),
		getNullString(stat.TsName),
		getNullInt64(stat.Quantity),
		getNullFloat64(stat.RetailPrice),
		getNullFloat64(stat.RetailAmount),
		// 41-50
		getNullFloat64(stat.CommissionPercent),
		getNullString(stat.OfficeName),
		getNullTime(stat.OrderDt),
		getNullTime(stat.SaleDt),
		getNullInt64(stat.DeliveryAmount),
		getNullInt64(stat.ReturnAmount),
		getNullInt64(stat.ReportType),
		getNullString(stat.Srid),
		getNullInt64(stat.Rid),
	).Scan(&stat.ID)

	if err != nil {
		return fmt.Errorf("failed to create stat: %w", err)
	}

	return nil
}

// GetStatDetails получает детальную статистику по фильтрам (с группировкой по nm_id)
func (r *StatRepository) GetStatDetails(userID int, dateFrom, dateTo string, page, pageSize int) ([]map[string]interface{}, error) {
	// Сначала получаем информацию о пользователе
	user, err := r.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	offset := (page - 1) * pageSize
	query := `
        SELECT 
            COALESCE(s.nm_id, 0) as nm_id,
            COALESCE(s.subject_name, 'Нет названия') as name,
            wa.photo, -- <-- Добавляем фото из wb_articles
            SUM(
                CASE 
                    WHEN s.ppvz_for_pay ~ '^[0-9]+\.?[0-9]*$' 
                    THEN CAST(s.ppvz_for_pay AS NUMERIC) 
                    ELSE 0 
                END
            ) as ppvz_for_pay,
            SUM(COALESCE(s.delivery_rub, 0)) as delivery_rub,
            SUM(COALESCE(s.deduction, 0)) as deduction,
            SUM(COALESCE(s.storage_fee, 0)) as storage_fee,
            SUM(COALESCE(s.additional_payment, 0)) as additional_payment,
            SUM(COALESCE(s.penalty, 0)) as penalty,
            SUM(
                CASE 
                    WHEN s.rebill_logistic_cost ~ '^[0-9]+\.?[0-9]*$' 
                    THEN CAST(s.rebill_logistic_cost AS NUMERIC)
                    ELSE 0 
                END
            ) as rebill_logistic_cost,
            SUM(
                CASE 
                    WHEN s.supplier_oper_name IN (1, 7) -- Продажа или Коррекция продаж
                    THEN 1 
                    ELSE 0 
                END
            ) as count_sales,
            SUM(
                CASE 
                    WHEN s.supplier_oper_name = 2 -- Возврат
                    THEN 1 
                    ELSE 0 
                END
            ) as count_refund,
            SUM(
                CASE 
                    WHEN s.supplier_oper_name IN (1, 7) -- Продажа или Коррекция продаж
                    THEN COALESCE(s.quantity, 0)
                    ELSE 0 
                END
            ) as sales,
            SUM(
                CASE 
                    WHEN s.supplier_oper_name = 2 -- Возврат
                    THEN COALESCE(s.return_amount, 0)
                    ELSE 0 
                END
            ) as returns
        FROM stat s
        LEFT JOIN wb_articles wa ON wa.articule::bigint = s.nm_id AND wa.id_user = s.user_id
        WHERE s.user_id = $1
            AND s.sale_dt BETWEEN $2 AND $3
            AND s.nm_id IS NOT NULL
            AND s.nm_id != 0
        GROUP BY s.nm_id, s.subject_name, wa.photo -- <-- Добавляем wa.photo в GROUP BY
        ORDER BY ppvz_for_pay DESC
        LIMIT $4 OFFSET $5
    `

	// Добавляем время для правильного диапазона дат
	dateToWithTime := dateTo + " 23:59:59"

	rows, err := r.db.Query(query, userID, dateFrom, dateToWithTime, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query stat details: %w", err)
	}
	defer rows.Close()

	var results []map[string]interface{}

	for rows.Next() {
		var nmID int64
		var name string
		var photo sql.NullString // <-- Добавляем переменную для фото
		var ppvzForPay, deliveryRub, deduction, storageFee, additionalPayment, penalty, rebillLogisticCost float64
		var countSales, countRefund, sales, returns int

		err := rows.Scan(
			&nmID,
			&name,
			&photo, // <-- Сканируем фото
			&ppvzForPay,
			&deliveryRub,
			&deduction,
			&storageFee,
			&additionalPayment,
			&penalty,
			&rebillLogisticCost,
			&countSales,
			&countRefund,
			&sales,
			&returns,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stat detail: %w", err)
		}

		// Расчет дополнительных полей
		deliveryPerUnit := 0.0
		if sales+returns > 0 {
			deliveryPerUnit = deliveryRub / float64(sales+returns)
		}
		costPriceTotal := 0.0

		// ВАШ ОРИГИНАЛЬНЫЙ РАСЧЕТ
		rebillLogisticCostInt := rebillLogisticCost
		netProfitBeforeTax := deliveryRub + penalty + deduction + storageFee + additionalPayment + rebillLogisticCostInt + float64(costPriceTotal)
		taxesAmount := 0.0

		// Учитываем налог пользователя (если он задан)
		if user.Taxes > 0 {
			taxesAmount = (ppvzForPay - netProfitBeforeTax) * (float64(user.Taxes) / 100)
		}

		netProfit := (ppvzForPay - netProfitBeforeTax) - taxesAmount

		item := map[string]interface{}{
			"nm_id":                nmID,
			"name":                 name,
			"photo":                getStringValue(photo), // <-- Используем реальное фото из базы
			"delivery_rub":         deliveryRub,
			"delivery_per_unit":    deliveryPerUnit,
			"deduction":            deduction,
			"storage_fee":          storageFee,
			"additional_payment":   additionalPayment,
			"penalty":              penalty,
			"rebill_logistic_cost": rebillLogisticCost,
			"ppvz_for_pay":         ppvzForPay,
			"count_sales":          countSales,
			"count_refund":         countRefund,
			"sales":                sales,
			"returns":              returns,
			"net_profit":           netProfit,
			"taxesAmount":          taxesAmount,
		}

		results = append(results, item)
	}

	return results, nil
}

// Добавь вспомогательную функцию если её нет
func getStringValue(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// GetStatDetailsCount получает общее количество записей для пагинации
func (r *StatRepository) GetStatDetailsCount(userID int, dateFrom, dateTo string) (int, error) {
	query := `
        SELECT COUNT(DISTINCT s.nm_id)
        FROM stat s
        WHERE s.user_id = $1
            AND s.sale_dt BETWEEN $2 AND $3
            AND s.nm_id IS NOT NULL
            AND s.nm_id != 0
    `

	dateToWithTime := dateTo + " 23:59:59"

	var count int
	err := r.db.QueryRow(query, userID, dateFrom, dateToWithTime).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get stat details count: %w", err)
	}

	return count, nil
}

// GetStatSummary получает итоговые суммы (уже агрегированные)
func (r *StatRepository) GetStatSummary(userID int, dateFrom, dateTo string) (map[string]interface{}, error) {
	query := `
        SELECT 
            SUM(
                CASE 
                    WHEN s.ppvz_for_pay ~ '^[0-9]+\.?[0-9]*$' 
                    THEN CAST(s.ppvz_for_pay AS NUMERIC) 
                    ELSE 0 
                END
            ) as total_ppvz_for_pay,
            SUM(COALESCE(s.delivery_rub, 0)) as total_delivery_rub,
            SUM(COALESCE(s.deduction, 0)) as total_deduction,
            SUM(COALESCE(s.storage_fee, 0)) as total_storage_fee,
            SUM(COALESCE(s.additional_payment, 0)) as total_additional_payment,
            SUM(COALESCE(s.penalty, 0)) as total_penalty,
            SUM(
                CASE 
                    WHEN s.supplier_oper_name = 1 
                    THEN 1 
                    ELSE 0 
                END
            ) as total_count_sales,
            SUM(
                CASE 
                    WHEN s.supplier_oper_name = 2 
                    THEN 1 
                    ELSE 0 
                END
            ) as total_count_refund,
            SUM(COALESCE(s.quantity, 0)) as total_quantity,
            SUM(COALESCE(s.return_amount, 0)) as total_return_amount,
            COUNT(DISTINCT s.nm_id) as unique_products
        FROM stat s
        WHERE s.user_id = $1
            AND s.sale_dt BETWEEN $2 AND $3
            AND s.nm_id IS NOT NULL
            AND s.nm_id != 0
    `

	dateToWithTime := dateTo + " 23:59:59"

	row := r.db.QueryRow(query, userID, dateFrom, dateToWithTime)

	var totalPpvzForPay, totalDeliveryRub, totalDeduction, totalStorageFee, totalAdditionalPayment, totalPenalty sql.NullFloat64
	var totalCountSales, totalCountRefund, totalQuantity, totalReturnAmount, uniqueProducts sql.NullInt64

	err := row.Scan(
		&totalPpvzForPay,
		&totalDeliveryRub,
		&totalDeduction,
		&totalStorageFee,
		&totalAdditionalPayment,
		&totalPenalty,
		&totalCountSales,
		&totalCountRefund,
		&totalQuantity,
		&totalReturnAmount,
		&uniqueProducts,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get stat summary: %w", err)
	}

	totalNetProfit := getFloatValue(totalPpvzForPay) - getFloatValue(totalDeliveryRub) -
		getFloatValue(totalDeduction) - getFloatValue(totalStorageFee) -
		getFloatValue(totalAdditionalPayment) - getFloatValue(totalPenalty)

	summary := map[string]interface{}{
		"total_ppvz_for_pay":       getFloatValue(totalPpvzForPay),
		"total_delivery_rub":       getFloatValue(totalDeliveryRub),
		"total_deduction":          getFloatValue(totalDeduction),
		"total_storage_fee":        getFloatValue(totalStorageFee),
		"total_additional_payment": getFloatValue(totalAdditionalPayment),
		"total_penalty":            getFloatValue(totalPenalty),
		"total_count_sales":        getIntValue(totalCountSales),
		"total_count_refund":       getIntValue(totalCountRefund),
		"total_quantity":           getIntValue(totalQuantity),
		"total_return_amount":      getIntValue(totalReturnAmount),
		"unique_products":          getIntValue(uniqueProducts),
		"total_net_profit":         totalNetProfit,
	}

	return summary, nil
}

// Вспомогательные функции
func getFloatValue(nf sql.NullFloat64) float64 {
	if nf.Valid {
		return nf.Float64
	}
	return 0.0
}

func getIntValue(ni sql.NullInt64) int64 {
	if ni.Valid {
		return ni.Int64
	}
	return 0
}

// Helper функция для генерации URL фото
func getPhotoURL(nmID int64) string {
	if nmID <= 0 {
		return ""
	}

	// Формируем URL по аналогии с вашим примером
	basketNum := nmID % 100
	vol := nmID / 1000000
	part := nmID / 1000

	return fmt.Sprintf("https://basket-%d.wbbasket.ru/vol%d/part%d/%d/images/big/1.webp",
		basketNum, vol, part, nmID)
}

// ExistsByHash проверяет существует ли запись с таким хешем в таблице stat
func (r *StatRepository) ExistsByHash(hash string) (bool, error) {
	query := `SELECT COUNT(*) FROM stat WHERE hash_info = $1`

	var count int
	err := r.db.QueryRow(query, hash).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check hash existence: %w", err)
	}

	return count > 0, nil
}

// Helper функции для работы с NULL значениями
func getNullString(ns sql.NullString) interface{} {
	if ns.Valid {
		return ns.String
	}
	return nil
}

func getNullInt64(ni sql.NullInt64) interface{} {
	if ni.Valid {
		return ni.Int64
	}
	return nil
}

func getNullFloat64(nf sql.NullFloat64) interface{} {
	if nf.Valid {
		return nf.Float64
	}
	return nil
}

func getNullTime(nt sql.NullTime) interface{} {
	if nt.Valid {
		return nt.Time
	}
	return nil
}

// GetDashboardStats получает статистику для дашборда
func (r *StatRepository) GetDashboardStats(userID int, dateFrom, dateTo string) (map[string]interface{}, error) {
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
        FROM stat s
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
func (r *StatRepository) GetChartData(userID int, dateFrom, dateTo string) (map[string]interface{}, error) {
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
        FROM stat s
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
        FROM stat s
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
func (r *StatRepository) GetMonthlyRevenue(userID int) (map[string]interface{}, error) {
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
        LEFT JOIN stat s ON 
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

// formatMonthLabel создает метку в формате "Янв 2025"
func formatMonthLabel(monthName string, year int) string {
	// Преобразуем английское название месяца в русское
	monthMap := map[string]string{
		"Jan": "Янв", "Feb": "Фев", "Mar": "Мар", "Apr": "Апр",
		"May": "Май", "Jun": "Июн", "Jul": "Июл", "Aug": "Авг",
		"Sep": "Сен", "Oct": "Окт", "Nov": "Ноя", "Dec": "Дек",
	}

	// Берем первые 3 символа
	shortName := ""
	if len(monthName) >= 3 {
		shortName = monthName[:3]
	}

	// Преобразуем в русское название
	ruMonth := shortName
	if ruName, ok := monthMap[shortName]; ok {
		ruMonth = ruName
	}

	// Возвращаем в формате "Янв 2025"
	return fmt.Sprintf("%s %d", ruMonth, year)
}

// formatMonthNameRU можно удалить или оставить, если она используется где-то еще
func formatMonthNameRU(monthName string) string {
	monthMap := map[string]string{
		"Jan": "Янв", "Feb": "Фев", "Mar": "Мар", "Apr": "Апр",
		"May": "Май", "Jun": "Июн", "Jul": "Июл", "Aug": "Авг",
		"Sep": "Сен", "Oct": "Окт", "Nov": "Ноя", "Dec": "Дек",
	}

	if len(monthName) >= 3 {
		shortName := monthName[:3]
		if ruName, ok := monthMap[shortName]; ok {
			return ruName
		}
	}

	return monthName
}

// Вспомогательная функция для генерации недостающих месяцев
func (r *StatRepository) generateMissingMonths(startDate, endDate time.Time, existingMonths map[string]bool,
	currentLabels []string, currentData []float64) ([]string, []float64) {

	allLabels := make([]string, 0)
	allData := make([]float64, 0)

	// Генерируем все месяцы в диапазоне
	for current := startDate; current.Before(endDate.AddDate(0, 1, 0)); current = current.AddDate(0, 1, 0) {
		monthName := formatMonthNameRU(current.Format("Jan"))
		label := fmt.Sprintf("%s %d", monthName, current.Year())

		// Проверяем, есть ли данные для этого месяца
		if existingMonths[label] {
			// Находим индекс в текущих данных
			for i, lbl := range currentLabels {
				if lbl == label {
					allLabels = append(allLabels, lbl)
					allData = append(allData, currentData[i])
					break
				}
			}
		} else {
			// Добавляем месяц с нулевыми данными
			allLabels = append(allLabels, label)
			allData = append(allData, 0)
		}

		// Ограничиваем 12 месяцами
		if len(allLabels) >= 12 {
			break
		}
	}

	return allLabels, allData
}

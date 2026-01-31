package stat

import (
	"database/sql"
	"fmt"
	"wbrost-go/internal/repository/database/postgres"
	"wbrost-go/internal/repository/user"
)

type AnalyticsRepository struct {
	db       *postgres.PostgresDB
	userRepo *user.UserRepository
}

func NewAnalyticsRepository(db *postgres.PostgresDB, userRepo *user.UserRepository) *AnalyticsRepository {
	return &AnalyticsRepository{db: db, userRepo: userRepo}
}

// GetStatDetails получает детальную статистику по фильтрам (с группировкой по nm_id)
func (r *AnalyticsRepository) GetStatDetails(userID int, dateFrom, dateTo string, page, pageSize int) ([]map[string]interface{}, error) {
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
        FROM wb_stats s
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

// GetStatDetailsCount получает общее количество записей для пагинации
func (r *AnalyticsRepository) GetStatDetailsCount(userID int, dateFrom, dateTo string) (int, error) {
	query := `
        SELECT COUNT(DISTINCT s.nm_id)
        FROM wb_stats s
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
func (r *AnalyticsRepository) GetStatSummary(userID int, dateFrom, dateTo string) (map[string]interface{}, error) {
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
        FROM wb_stats s
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

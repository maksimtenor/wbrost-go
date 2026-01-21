package repository

import (
	"database/sql"
	"fmt"
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
        WHERE s.user_id = $1
            AND s.sale_dt BETWEEN $2 AND $3
            AND s.nm_id IS NOT NULL
            AND s.nm_id != 0
        GROUP BY s.nm_id, s.subject_name
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
		var ppvzForPay, deliveryRub, deduction, storageFee, additionalPayment, penalty, rebillLogisticCost float64
		var countSales, countRefund, sales, returns int

		err := rows.Scan(
			&nmID,
			&name,
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
			"photo":                getPhotoURL(nmID),
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

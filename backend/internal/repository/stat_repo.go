package repository

import (
	"database/sql"
	"fmt"
	"wbrost-go/internal/entity"
)

type StatRepository struct {
	db *PostgresDB
}

func NewStatRepository(db *PostgresDB) *StatRepository {
	return &StatRepository{db: db}
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

// GetStatDetails получает детальную статистику по фильтрам
func (r *StatRepository) GetStatDetails(userID int, dateFrom, dateTo string) ([]map[string]interface{}, error) {
	query := `
        SELECT 
            COALESCE(s.nm_id, 0) as nm_id,
            COALESCE(s.subject_name, 'Нет названия') as name,
            COALESCE(s.delivery_rub, 0) as delivery_rub,
            COALESCE(s.deduction, 0) as deduction,
            COALESCE(s.storage_fee, 0) as storage_fee,
            COALESCE(s.additional_payment, 0) as additional_payment,
            COALESCE(s.penalty, 0) as penalty,
            COALESCE(s.rebill_logistic_cost, '0') as rebill_logistic_cost,
            COALESCE(s.ppvz_for_pay, 0) as ppvz_for_pay,
            COALESCE(s.quantity, 0) as sales,
            COALESCE(s.return_amount, 0) as returns,
            COALESCE(s.delivery_amount, 0) as delivery_amount
        FROM stat s
        WHERE s.user_id = $1
            AND s.sale_dt BETWEEN $2 AND $3
            AND s.nm_id IS NOT NULL
            AND s.nm_id != 0
        ORDER BY s.nm_id
    `

	// Если нужно включить все записи (включая nm_id = 0), уберите последнее условие

	rows, err := r.db.Query(query, userID, dateFrom, dateTo)
	if err != nil {
		return nil, fmt.Errorf("failed to query stat details: %w", err)
	}
	defer rows.Close()

	var results []map[string]interface{}

	for rows.Next() {
		var nmID int64
		var name string
		var deliveryRub, deduction, storageFee, additionalPayment, penalty, ppvzForPay float64
		var rebillLogisticCost string
		var sales, returns, deliveryAmount int

		err := rows.Scan(
			&nmID,
			&name,
			&deliveryRub,
			&deduction,
			&storageFee,
			&additionalPayment,
			&penalty,
			&rebillLogisticCost,
			&ppvzForPay,
			&sales,
			&returns,
			&deliveryAmount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stat detail: %w", err)
		}

		// Расчет дополнительных полей
		var deliveryPerUnit float64
		if sales > 0 {
			deliveryPerUnit = deliveryRub / float64(sales)
		}

		var netProfit float64
		netProfit = ppvzForPay - deliveryRub - deduction - storageFee - additionalPayment - penalty

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
			"sales":                sales,
			"returns":              returns,
			"delivery_amount":      deliveryAmount,
			"cost_price":           0, // Здесь можно добавить расчет себестоимости если есть данные
			"net_profit":           netProfit,
		}

		results = append(results, item)
	}

	return results, nil
}

// GetStatSummary получает итоговые суммы
func (r *StatRepository) GetStatSummary(userID int, dateFrom, dateTo string) (map[string]interface{}, error) {
	query := `
        SELECT 
            COALESCE(SUM(delivery_rub), 0) as total_delivery,
            COALESCE(SUM(deduction), 0) as total_deduction,
            COALESCE(SUM(storage_fee), 0) as total_storage,
            COALESCE(SUM(additional_payment), 0) as total_additional,
            COALESCE(SUM(penalty), 0) as total_penalty,
            COALESCE(SUM(ppvz_for_pay), 0) as total_revenue,
            COALESCE(SUM(quantity), 0) as total_sales,
            COALESCE(SUM(return_amount), 0) as total_returns,
            COUNT(DISTINCT nm_id) as unique_products
        FROM stat
        WHERE user_id = $1
            AND sale_dt BETWEEN $2 AND $3
            AND nm_id IS NOT NULL
            AND nm_id != 0
    `

	row := r.db.QueryRow(query, userID, dateFrom, dateTo)

	var totalDelivery, totalDeduction, totalStorage, totalAdditional, totalPenalty, totalRevenue float64
	var totalSales, totalReturns, uniqueProducts int

	err := row.Scan(
		&totalDelivery,
		&totalDeduction,
		&totalStorage,
		&totalAdditional,
		&totalPenalty,
		&totalRevenue,
		&totalSales,
		&totalReturns,
		&uniqueProducts,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get stat summary: %w", err)
	}

	totalNetProfit := totalRevenue - totalDelivery - totalDeduction - totalStorage - totalAdditional - totalPenalty

	summary := map[string]interface{}{
		"total_delivery_rub":       totalDelivery,
		"total_deduction":          totalDeduction,
		"total_storage_fee":        totalStorage,
		"total_additional_payment": totalAdditional,
		"total_penalty":            totalPenalty,
		"total_ppvz_for_pay":       totalRevenue,
		"total_sales":              totalSales,
		"total_returns":            totalReturns,
		"unique_products":          uniqueProducts,
		"total_net_profit":         totalNetProfit,
	}

	return summary, nil
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

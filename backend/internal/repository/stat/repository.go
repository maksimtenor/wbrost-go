package stat

import (
	"fmt"
	"wbrost-go/internal/entity"
	"wbrost-go/internal/repository/database/postgres"
)

type StatRepository struct {
	db *postgres.PostgresDB
}

//type StatRepositoryInterface interface {
//	Create(stat *entity.Stat) error
//	GetStatDetails(userID int, dateFrom, dateTo string, page, pageSize int) ([]map[string]interface{}, error)
//	GetStatDetailsCount(userID int, dateFrom, dateTo string) (int, error)
//	GetStatSummary(userID int, dateFrom, dateTo string) (map[string]interface{}, error)
//	ExistsByHash(hash string) (bool, error)
//}

func NewStatRepository(db *postgres.PostgresDB) *StatRepository {
	return &StatRepository{db: db}
}

// Create создает новую запись статистики в таблицу stat
func (r *StatRepository) Create(stat *entity.Stat) error {
	query := `
		INSERT INTO wb_stats (
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

// ExistsByHash проверяет существует ли запись с таким хешем в таблице stat
func (r *StatRepository) ExistsByHash(hash string) (bool, error) {
	query := `SELECT COUNT(*) FROM wb_stats WHERE hash_info = $1`

	var count int
	err := r.db.QueryRow(query, hash).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check hash existence: %w", err)
	}

	return count > 0, nil
}

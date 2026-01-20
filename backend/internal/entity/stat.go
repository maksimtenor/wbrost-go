package entity

import (
	"database/sql"
	"time"
)

// Stat - соответствует таблице stat из вашей БД
type Stat struct {
	ID                  int64           `json:"id" db:"id"`
	HashInfo            string          `json:"hash_info" db:"hash_info"`
	UserID              int             `json:"user_id" db:"user_id"`
	Nmid                sql.NullInt64   `json:"nm_id" db:"nm_id"`
	PpvzForPay          sql.NullString  `json:"ppvz_for_pay" db:"ppvz_for_pay"`
	SupplierOperName    sql.NullInt64   `json:"supplier_oper_name" db:"supplier_oper_name"` // Обратите внимание: integer в БД
	DeliveryRub         sql.NullFloat64 `json:"delivery_rub" db:"delivery_rub"`
	Penalty             sql.NullFloat64 `json:"penalty" db:"penalty"`
	AdditionalPayment   sql.NullFloat64 `json:"additional_payment" db:"additional_payment"`
	StorageFee          sql.NullFloat64 `json:"storage_fee" db:"storage_fee"`
	RebillLogisticCost  sql.NullString  `json:"rebill_logistic_cost" db:"rebill_logistic_cost"`
	AcquiringFee        sql.NullFloat64 `json:"acquiring_fee" db:"acquiring_fee"`
	AcquiringPercent    sql.NullFloat64 `json:"acquiring_percent" db:"acquiring_percent"`
	PpvzSalesCommission sql.NullFloat64 `json:"ppvz_sales_commission" db:"ppvz_sales_commission"`
	Deduction           sql.NullFloat64 `json:"deduction" db:"deduction"`
	PpvzSppPrc          sql.NullString  `json:"ppvz_spp_prc" db:"ppvz_spp_prc"`
	PpvzKvwPrcBase      sql.NullString  `json:"ppvz_kvw_prc_base" db:"ppvz_kvw_prc_base"`
	PpvzKvwPrc          sql.NullString  `json:"ppvz_kvw_prc" db:"ppvz_kvw_prc"`
	Acceptance          sql.NullFloat64 `json:"acceptance" db:"acceptance"`
	DlvPrc              sql.NullFloat64 `json:"dlv_prc" db:"dlv_prc"`
	CreatedAt           time.Time       `json:"created_at" db:"created_at"`
	RrDt                sql.NullTime    `json:"rr_dt" db:"rr_dt"`
	ShkID               sql.NullInt64   `json:"shk_id" db:"shk_id"`
	StickerID           sql.NullString  `json:"sticker_id" db:"sticker_id"`
	GiID                sql.NullInt64   `json:"gi_id" db:"gi_id"`
	RealizationreportID sql.NullInt64   `json:"realizationreport_id" db:"realizationreport_id"`
	Barcode             sql.NullString  `json:"barcode" db:"barcode"`
	BonusTypeName       sql.NullString  `json:"bonus_type_name" db:"bonus_type_name"`
	LastError           sql.NullString  `json:"last_error" db:"last_error"`
	BrandName           sql.NullString  `json:"brand_name" db:"brand_name"`
	PpvzOfficeID        sql.NullInt64   `json:"ppvz_office_id" db:"ppvz_office_id"`
	AssemblyID          sql.NullInt64   `json:"assembly_id" db:"assembly_id"` // Обратите внимание: bigint в БД
	SaName              sql.NullString  `json:"sa_name" db:"sa_name"`
	PpvzVwNds           sql.NullString  `json:"ppvz_vw_nds" db:"ppvz_vw_nds"`
	PpvzVw              sql.NullString  `json:"ppvz_vw" db:"ppvz_vw"`
	GiBoxTypeName       sql.NullString  `json:"gi_box_type_name" db:"gi_box_type_name"`
	SubjectName         sql.NullString  `json:"subject_name" db:"subject_name"`
	TsName              sql.NullString  `json:"ts_name" db:"ts_name"`
	Quantity            sql.NullInt64   `json:"quantity" db:"quantity"`
	RetailPrice         sql.NullFloat64 `json:"retail_price" db:"retail_price"`
	RetailAmount        sql.NullFloat64 `json:"retail_amount" db:"retail_amount"`
	CommissionPercent   sql.NullFloat64 `json:"commission_percent" db:"commission_percent"`
	OfficeName          sql.NullString  `json:"office_name" db:"office_name"`
	OrderDt             sql.NullTime    `json:"order_dt" db:"order_dt"`
	SaleDt              sql.NullTime    `json:"sale_dt" db:"sale_dt"`
	DeliveryAmount      sql.NullInt64   `json:"delivery_amount" db:"delivery_amount"`
	ReturnAmount        sql.NullInt64   `json:"return_amount" db:"return_amount"`
	ReportType          sql.NullInt64   `json:"report_type" db:"report_type"`
	Srid                sql.NullString  `json:"srid" db:"srid"`
	Rid                 sql.NullInt64   `json:"rid" db:"rid"`
}

package entity

import (
	"time"
)

type Stat struct {
	ID                  int        `json:"id" db:"id"`
	HashInfo            *string    `json:"hash_info,omitempty" db:"hash_info"`
	UserID              int        `json:"user_id" db:"user_id"`
	NmID                *int64     `json:"nm_id,omitempty" db:"nm_id"`
	PpvzForPay          *string    `json:"ppvz_for_pay,omitempty" db:"ppvz_for_pay"`
	SupplierOperName    *int       `json:"supplier_oper_name,omitempty" db:"supplier_oper_name"`
	DeliveryRub         float64    `json:"delivery_rub" db:"delivery_rub"`
	Penalty             float64    `json:"penalty" db:"penalty"`
	AdditionalPayment   float64    `json:"additional_payment" db:"additional_payment"`
	StorageFee          float64    `json:"storage_fee" db:"storage_fee"`
	RebillLogisticCost  string     `json:"rebill_logistic_cost" db:"rebill_logistic_cost"`
	AcquiringFee        float64    `json:"acquiring_fee" db:"acquiring_fee"`
	AcquiringPercent    float64    `json:"acquiring_percent" db:"acquiring_percent"`
	PpvzSalesCommission float64    `json:"ppvz_sales_commission" db:"ppvz_sales_commission"`
	Deduction           float64    `json:"deduction" db:"deduction"`
	PpvzSppPrc          *string    `json:"ppvz_spp_prc,omitempty" db:"ppvz_spp_prc"`
	PpvzKvwPrcBase      *string    `json:"ppvz_kvw_prc_base,omitempty" db:"ppvz_kvw_prc_base"`
	PpvzKvwPrc          *string    `json:"ppvz_kvw_prc,omitempty" db:"ppvz_kvw_prc"`
	Acceptance          float64    `json:"acceptance" db:"acceptance"`
	DlvPrc              float64    `json:"dlv_prc" db:"dlv_prc"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	RrDt                *time.Time `json:"rr_dt,omitempty" db:"rr_dt"`
	ShkID               *int64     `json:"shk_id,omitempty" db:"shk_id"`
	StickerID           *string    `json:"sticker_id,omitempty" db:"sticker_id"`
	GiID                *int64     `json:"gi_id,omitempty" db:"gi_id"`
	RealizationreportID *int64     `json:"realizationreport_id,omitempty" db:"realizationreport_id"`
	Barcode             *string    `json:"barcode,omitempty" db:"barcode"`
	BonusTypeName       *string    `json:"bonus_type_name,omitempty" db:"bonus_type_name"`
	LastError           *string    `json:"last_error,omitempty" db:"last_error"`
	BrandName           *string    `json:"brand_name,omitempty" db:"brand_name"`
	PpvzOfficeID        *int64     `json:"ppvz_office_id,omitempty" db:"ppvz_office_id"`
	AssemblyID          *int64     `json:"assembly_id,omitempty" db:"assembly_id"`
	SaName              *string    `json:"sa_name,omitempty" db:"sa_name"`
	PpvzVwNds           *string    `json:"ppvz_vw_nds,omitempty" db:"ppvz_vw_nds"`
	PpvzVw              *string    `json:"ppvz_vw,omitempty" db:"ppvz_vw"`
	GiBoxTypeName       *string    `json:"gi_box_type_name,omitempty" db:"gi_box_type_name"`
	SubjectName         *string    `json:"subject_name,omitempty" db:"subject_name"`
	TsName              *string    `json:"ts_name,omitempty" db:"ts_name"`
	Quantity            *int       `json:"quantity,omitempty" db:"quantity"`
	RetailPrice         float64    `json:"retail_price" db:"retail_price"`
	RetailAmount        float64    `json:"retail_amount" db:"retail_amount"`
	CommissionPercent   float64    `json:"commission_percent" db:"commission_percent"`
	OfficeName          *string    `json:"office_name,omitempty" db:"office_name"`
	OrderDt             *time.Time `json:"order_dt,omitempty" db:"order_dt"`
	SaleDt              *time.Time `json:"sale_dt,omitempty" db:"sale_dt"`
	DeliveryAmount      *int       `json:"delivery_amount,omitempty" db:"delivery_amount"`
	ReturnAmount        *int       `json:"return_amount,omitempty" db:"return_amount"`
	ReportType          *int       `json:"report_type,omitempty" db:"report_type"`
	Srid                *string    `json:"srid,omitempty" db:"srid"`
	Rid                 *int64     `json:"rid,omitempty" db:"rid"`
}

// Для агрегированных данных
type StatSummary struct {
	Month           string  `json:"month"`
	TotalSell       float64 `json:"total_sell"`
	TotalCountSales int     `json:"total_count_sales"`
	TotalRefund     int     `json:"total_refund"`
}

type ChartData struct {
	Labels   []string       `json:"labels"`
	Datasets []ChartDataset `json:"datasets"`
}

type ChartDataset struct {
	Label           string        `json:"label"`
	BackgroundColor string        `json:"backgroundColor,omitempty"`
	BorderColor     string        `json:"borderColor,omitempty"`
	Data            []interface{} `json:"data"`
}

// Константы из Yii2
const (
	TypeSale                                  = 1
	TypeRefund                                = 2
	TypeDelivery                              = 3
	TypeHold                                  = 4
	TypePenalty                               = 5
	TypeStorage                               = 6
	TypeSalesCorrection                       = 7
	TypeAdvancePaymentForGoodsWithoutMovement = 8
	TypeStorageRecalculation                  = 9
	TypePaidAcceptance                        = 10
	TypeDeliveryCorrection                    = 11
	TypeAcquiringCorrection                   = 12
	TypeDamageCompensation                    = 13
	TypeLostCompensations                     = 14
	TypeCompensationForDefects                = 15
	TypeCompensationUponReturn                = 16
	TypeCompensationForReplacedGoods          = 17
	TypeReimbursementOfTransportCost          = 18
	DataTypeLine                              = "line"
	DataTypeBar                               = "bar"
)

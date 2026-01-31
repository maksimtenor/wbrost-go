package entity

import "database/sql"

// WBArticles - соответствует таблице wb_articles в БД
type WBArticles struct {
	ID              int            `json:"id" db:"id"`
	UserID          int            `json:"id_user" db:"id_user"`
	Articule        string         `json:"articule" db:"articule"`
	Name            sql.NullString `json:"name" db:"name"`
	Photo           sql.NullString `json:"photo" db:"photo"`
	CostPrice       sql.NullString `json:"cost_price" db:"cost_price"`
	Created         sql.NullTime   `json:"created" db:"created"`
	Updated         sql.NullTime   `json:"updated" db:"updated"`
	UpdatedAt       sql.NullTime   `json:"updated_at" db:"updated_at"`
	SelfRansom      sql.NullInt64  `json:"self_ransom" db:"self_ransom"`
	SelfRansomPrice sql.NullString `json:"self_ransom_price" db:"self_ransom_price"`
	RusSize         sql.NullString `json:"rus_size" db:"rus_size"`
	EuSize          sql.NullString `json:"eu_size" db:"eu_size"`
	ChrtID          sql.NullInt64  `json:"chrt_id" db:"chrt_id"`
	Barcode         sql.NullString `json:"barcode" db:"barcode"`
	InternalID      sql.NullString `json:"internal_id" db:"internal_id"`
}

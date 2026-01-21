package entity

// WBArticle - структура для ответа API карточек WB
type WBArticle struct {
	NmID        int    `json:"nmID"`
	ImtID       int    `json:"imtID"`
	NmUUID      string `json:"nmUUID"`
	SubjectID   int    `json:"subjectID"`
	SubjectName string `json:"subjectName"`
	VendorCode  string `json:"vendorCode"`
	Brand       string `json:"brand"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Photos      []struct {
		Big      string `json:"big"`
		C246x328 string `json:"c246x328"`
		C516x688 string `json:"c516x688"`
		Square   string `json:"square"`
		Tm       string `json:"tm"`
	} `json:"photos"`
	Sizes []struct {
		ChrtID   int      `json:"chrtID"`
		TechSize string   `json:"techSize"`
		Skus     []string `json:"skus"`
		WbSize   string   `json:"wbSize"`
	} `json:"sizes"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// WBArticleResponse - структура полного ответа API
type WBArticleResponse struct {
	Cards  []WBArticle `json:"cards"`
	Cursor struct {
		UpdatedAt string `json:"updatedAt"`
		NmID      int    `json:"nmID"`
		Total     int    `json:"total"`
	} `json:"cursor"`
}

// WBArticleRequestCursor - курсор для пагинации
type WBArticleRequestCursor struct {
	Limit     int    `json:"limit"`
	UpdatedAt string `json:"updatedAt,omitempty"`
	NmID      int    `json:"nmID,omitempty"`
}

// WBArticleRequestSettings - настройки запроса
type WBArticleRequestSettings struct {
	Cursor WBArticleRequestCursor `json:"cursor"`
	Filter struct {
		WithPhoto int `json:"withPhoto"`
	} `json:"filter"`
}

// WBArticleRequest - полный запрос
type WBArticleRequest struct {
	Settings WBArticleRequestSettings `json:"settings"`
}

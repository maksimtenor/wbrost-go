package wb

// Article - структура для ответа API карточек WB
type Article struct {
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

// ArticleResponse - структура полного ответа API
type ArticleResponse struct {
	Cards  []Article `json:"cards"`
	Cursor struct {
		UpdatedAt string `json:"updatedAt"`
		NmID      int    `json:"nmID"`
		Total     int    `json:"total"`
	} `json:"cursor"`
}

// ArticleRequestCursor - курсор для пагинации
type ArticleRequestCursor struct {
	Limit     int    `json:"limit"`
	UpdatedAt string `json:"updatedAt,omitempty"`
	NmID      int    `json:"nmID,omitempty"`
}

// ArticleRequestSettings - настройки запроса
type ArticleRequestSettings struct {
	Cursor ArticleRequestCursor `json:"cursor"`
	Filter struct {
		WithPhoto int `json:"withPhoto"`
	} `json:"filter"`
}

// ArticleRequest - полный запрос
type ArticleRequest struct {
	Settings ArticleRequestSettings `json:"settings"`
}

// PassesResponse структура для ответа проверки токена
type PassesResponse []struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error"`
}

// ErrorResponse Структура для ошибки WB API (401/403/429)
type ErrorResponse struct {
	Status     int    `json:"status"`
	StatusText string `json:"statusText"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Detail     string `json:"detail"`
}

// Pass Структура для успешного ответа (пропуск)
type Pass struct {
	ID            int    `json:"id"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	CarModel      string `json:"carModel"`
	CarNumber     string `json:"carNumber"`
	OfficeName    string `json:"officeName"`
	OfficeAddress string `json:"officeAddress"`
	OfficeID      int    `json:"officeId"`
	DateEnd       string `json:"dateEnd"`
}

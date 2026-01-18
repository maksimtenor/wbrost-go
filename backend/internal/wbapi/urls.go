package wbapi

const (
	UrlStats       = "https://seller-analytics-api.wildberries.ru/"
	UrlStatsNew    = "https://statistics-api.wildberries.ru/"
	UrlCard        = "https://content-api.wildberries.ru/"
	UrlContent     = "https://seller-analytics-api.wildberries.ru/"
	UrlMarketplace = "https://marketplace-api.wildberries.ru/"

	// Endpoints
	incomes       = "api/v1/supplier/incomes"
	detailsV1     = "api/v1/supplier/reportDetailByPeriod"
	detailsV5     = "api/v5/supplier/reportDetailByPeriod"
	wborders      = "api/v1/supplier/orders"
	taskCreate    = "api/v1/delayed-gen/tasks/create"
	taskStatus    = "api/v1/delayed-gen/tasks"
	taskDownload  = "api/v1/delayed-gen/tasks/download"
	cardsList     = "content/v2/get/cards/list"
	detailHistory = "api/v2/nm-report/detail/history"
	passes        = "api/v3/passes" // Добавили!
)

// URLStatistic возвращает базовый URL статистики
func URLStatistic() string {
	return UrlStats
}

// URLIncomes возвращает URL для получения поставок
func URLIncomes() string {
	return UrlStats + incomes
}

// URLDetails возвращает URL для детального отчета v1
func URLDetails() string {
	return UrlStats + detailsV1
}

// URLDetail5 возвращает URL для детального отчета v5
func URLDetail5() string {
	return UrlStatsNew + detailsV5
}

// URLOrders возвращает URL для получения заказов
func URLOrders() string {
	return UrlStats + wborders
}

// URLTaskCreate возвращает URL для создания задачи
func URLTaskCreate() string {
	return UrlStats + taskCreate
}

// URLTaskStatus возвращает URL для проверки статуса задачи
func URLTaskStatus() string {
	return UrlStats + taskStatus
}

// URLTaskDownload возвращает URL для скачивания задачи
func URLTaskDownload() string {
	return UrlStats + taskDownload
}

// URLContent возвращает базовый URL контента
func URLContent() string {
	return UrlContent
}

// URLCardsList возвращает URL для получения списка карточек
func URLCardsList() string {
	return UrlCard + cardsList
}

// URLDetailHistory возвращает URL для истории детального отчета
func URLDetailHistory() string {
	return UrlContent + detailHistory
}

// URLPasses возвращает URL для проверки токена (НОВОЕ!)
func URLPasses() string {
	return UrlMarketplace + passes
}

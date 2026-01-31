package wb

const (
	// Base URLs
	BaseURLStats       = "https://seller-analytics-api.wildberries.ru/"
	BaseURLStatsNew    = "https://statistics-api.wildberries.ru/"
	BaseURLCard        = "https://content-api.wildberries.ru/"
	BaseURLContent     = "https://seller-analytics-api.wildberries.ru/"
	BaseURLMarketplace = "https://marketplace-api.wildberries.ru/"

	// Endpoints
	EndpointIncomes       = "api/v1/supplier/incomes"
	EndpointDetailsV1     = "api/v1/supplier/reportDetailByPeriod"
	EndpointDetailsV5     = "api/v5/supplier/reportDetailByPeriod"
	EndpointOrders        = "api/v1/supplier/orders"
	EndpointTaskCreate    = "api/v1/delayed-gen/tasks/create"
	EndpointTaskStatus    = "api/v1/delayed-gen/tasks"
	EndpointTaskDownload  = "api/v1/delayed-gen/tasks/download"
	EndpointCardsList     = "content/v2/get/cards/list"
	EndpointDetailHistory = "api/v2/nm-report/detail/history"
	EndpointPasses        = "api/v3/passes"
)

// Endpoint тип для эндпоинтов API Wildberries
type Endpoint string

// Константы эндпоинтов с типом Endpoint для type-safe использования
const (
	Incomes       Endpoint = EndpointIncomes
	DetailsV1     Endpoint = EndpointDetailsV1
	DetailsV5     Endpoint = EndpointDetailsV5
	Orders        Endpoint = EndpointOrders
	TaskCreate    Endpoint = EndpointTaskCreate
	TaskStatus    Endpoint = EndpointTaskStatus
	TaskDownload  Endpoint = EndpointTaskDownload
	CardsList     Endpoint = EndpointCardsList
	DetailHistory Endpoint = EndpointDetailHistory
	Passes        Endpoint = EndpointPasses
)

// URLFor возвращает полный URL для указанного эндпоинта
func URLFor(endpoint Endpoint) string {
	switch endpoint {
	case Incomes, DetailsV1, Orders, TaskCreate, TaskStatus, TaskDownload:
		return BaseURLStats + string(endpoint)
	case DetailsV5:
		return BaseURLStatsNew + string(endpoint)
	case CardsList, DetailHistory:
		return BaseURLCard + string(endpoint)
	case Passes:
		return BaseURLMarketplace + string(endpoint)
	default:
		// fallback на основной stats URL
		return BaseURLStats + string(endpoint)
	}
}

// Конкретные методы для удобства и обратной совместимости

// URLIncomes возвращает URL для получения поставок
func URLIncomes() string {
	return BaseURLStats + EndpointIncomes
}

// URLDetailsV1 возвращает URL для детального отчета v1
func URLDetailsV1() string {
	return BaseURLStats + EndpointDetailsV1
}

// URLDetailsV5 возвращает URL для детального отчета v5
func URLDetailsV5() string {
	return BaseURLStatsNew + EndpointDetailsV5
}

// URLOrders возвращает URL для получения заказов
func URLOrders() string {
	return BaseURLStats + EndpointOrders
}

// URLTaskCreate возвращает URL для создания задачи
func URLTaskCreate() string {
	return BaseURLStats + EndpointTaskCreate
}

// URLTaskStatus возвращает URL для проверки статуса задачи
func URLTaskStatus() string {
	return BaseURLStats + EndpointTaskStatus
}

// URLTaskDownload возвращает URL для скачивания задачи
func URLTaskDownload() string {
	return BaseURLStats + EndpointTaskDownload
}

// URLCardsList возвращает URL для получения списка карточек
func URLCardsList() string {
	return BaseURLCard + EndpointCardsList
}

// URLDetailHistory возвращает URL для истории детального отчета
func URLDetailHistory() string {
	return BaseURLCard + EndpointDetailHistory
}

// URLPasses возвращает URL для проверки токена
func URLPasses() string {
	return BaseURLMarketplace + EndpointPasses
}

// BaseURLs возвращает все базовые URL в виде map
func BaseURLs() map[string]string {
	return map[string]string{
		"stats":       BaseURLStats,
		"stats_new":   BaseURLStatsNew,
		"card":        BaseURLCard,
		"content":     BaseURLContent,
		"marketplace": BaseURLMarketplace,
	}
}

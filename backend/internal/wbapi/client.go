package wbapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// WBClient клиент для работы с API Wildberries
type WBClient struct {
	Token  string
	Client *http.Client
}

// NewWBClient создает новый клиент
func NewWBClient(token string) *WBClient {
	return &WBClient{
		Token: token,
		Client: &http.Client{
			Timeout: 10 * time.Second, // Добавляем таймаут
		},
	}
}

// PassesResponse структура для ответа проверки токена
type PassesResponse []struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error"`
}

// WBErrorResponse Структура для ошибки WB API (401/403/429)
type WBErrorResponse struct {
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

// CheckToken проверяет валидность токена через API WB
func (c *WBClient) CheckToken() (bool, error) {
	// Используем конструктор URL!
	url := URLPasses()

	// Создаем запрос
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("ошибка создания запроса: %v", err)
	}

	// Добавляем заголовки как в PHP
	req.Header.Set("Authorization", c.Token)
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := c.Client.Do(req)
	if err != nil {
		return false, fmt.Errorf("ошибка сети: %v", err)
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	// ГЛАВНАЯ ПРОВЕРКА: HTTP статус код!
	// Если статус 200 - токен валиден
	// Если 401, 403 - токен невалиден
	// Если 429 - лимит запросов (но токен валиден!)
	switch resp.StatusCode {
	case 200:
		// Проверяем что это действительно массив пропусков
		var passes []Pass
		if err := json.Unmarshal(body, &passes); err != nil {
			// Если не массив пропусков, но статус 200 - всё равно валиден
			return true, nil
		}
		return true, nil // Токен валиден

	case 401, 403:
		// Пробуем распарсить ошибку для деталей
		var wbError WBErrorResponse
		if err := json.Unmarshal(body, &wbError); err == nil {
			return false, fmt.Errorf("токен недействителен: %s", wbError.Message)
		}
		return false, nil // Токен невалиден

	case 429:
		return false, fmt.Errorf("лимит запросов к WB API. Подождите")

	default:
		return false, fmt.Errorf("неизвестный ответ от WB API: статус %d", resp.StatusCode)
	}
}

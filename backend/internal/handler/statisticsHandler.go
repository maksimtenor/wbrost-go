package handler

import (
	"encoding/json"
	"net/http"
	"wbrost-go/internal/wbapi"
)

// StatisticsHandler обрабатывает запросы к статистике
type StatisticsHandler struct {
	WBClient *wbapi.WBClient
}

// NewStatisticsHandler создает новый обработчик статистики
func NewStatisticsHandler(wbClient *wbapi.WBClient) *StatisticsHandler {
	return &StatisticsHandler{WBClient: wbClient}
}

// CheckTokenHandler проверяет токен Wildberries
func (h *StatisticsHandler) CheckTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Пока просто заглушка
	response := map[string]interface{}{
		"valid":   true,
		"message": "Token check placeholder",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetStatistics возвращает статистику
func (h *StatisticsHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	// Пока просто заглушка
	response := map[string]interface{}{
		"data": "Statistics placeholder",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

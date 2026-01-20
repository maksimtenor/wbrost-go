package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"wbrost-go/internal/entity"
	"wbrost-go/internal/repository"

	"github.com/golang-jwt/jwt/v4"
)

type WBStatsHandler struct {
	userRepo  *repository.UserRepository
	statsRepo *repository.WBStatsGetRepository
	jwtSecret []byte
}

func NewWBStatsHandler(userRepo *repository.UserRepository, statsRepo *repository.WBStatsGetRepository, jwtSecret string) *WBStatsHandler {
	return &WBStatsHandler{
		userRepo:  userRepo,
		statsRepo: statsRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

// GetWBReports - GET /api/wb/stats | Получение списка репортов(заказов отчетов из бд)
func (h *WBStatsHandler) GetWBReports(w http.ResponseWriter, r *http.Request) {
	// Получить пользователя по токену
	user, err := h.getUserFromRequest(r)
	if err != nil {
		respondWithJSON(w, http.StatusUnauthorized, entity.ErrorResponse{Error: "Unauthorized"})
		return
	}

	// Получить отчеты для этого пользователя
	reports, err := h.statsRepo.GetByUserID(user.ID)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, entity.ErrorResponse{Error: "Failed to get reports"})
		return
	}

	// Подготовить ответ
	response := make([]map[string]interface{}, len(reports))
	for i, report := range reports {
		response[i] = map[string]interface{}{
			"id":         report.ID,
			"user_id":    report.UserID,
			"status":     getStatusValue(report.Status),
			"date_from":  report.DateFrom,
			"date_to":    report.DateTo,
			"created":    report.Created.Format("2006-01-02 15:04:05"),
			"updated":    report.Updated.Format("2006-01-02 15:04:05"),
			"last_error": getStringValue(report.LastError),
		}
	}

	respondWithJSON(w, http.StatusOK, response)
}

// CreateWBReport - POST /api/wb/stats | Заказать репорт (добавить запись с заказом в бд)
func (h *WBStatsHandler) CreateWBReport(w http.ResponseWriter, r *http.Request) {
	// Get user from token
	user, err := h.getUserFromRequest(r)
	if err != nil {
		respondWithJSON(w, http.StatusUnauthorized, entity.ErrorResponse{Error: "Unauthorized"})
		return
	}

	// Парсим запрос
	var req struct {
		DateFrom string `json:"dateFrom"`
		DateTo   string `json:"dateTo"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithJSON(w, http.StatusBadRequest, entity.ErrorResponse{Error: "Invalid request body"})
		return
	}

	// Валидация данных
	if req.DateFrom == "" || req.DateTo == "" {
		respondWithJSON(w, http.StatusBadRequest, entity.ErrorResponse{Error: "DateFrom and DateTo are required"})
		return
	}

	// Создание репорта в бд
	stats := &entity.WBStatsGet{ // Меняем repository.WBStatsGet на entity.WBStatsGet
		UserID:   user.ID,
		Status:   getNullInt64(0), // 0 = в обработке
		DateFrom: req.DateFrom,
		DateTo:   req.DateTo,
	}

	if err := h.statsRepo.Create(stats); err != nil {
		respondWithJSON(w, http.StatusInternalServerError, entity.ErrorResponse{Error: "Failed to create report: " + err.Error()})
		return
	}

	// Запустите асинхронную задачу для получения данных из API WB.
	//go h.fetchWBDataAsync(stats.ID, user, req.DateFrom, req.DateTo)

	// Вернуть "Успех"
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"id":      stats.ID,
		"success": true,
		"message": "Отчет поставлен в очередь на формирование",
	})
}

// Асинхронная функция для получения данных из WB.
func (h *WBStatsHandler) fetchWBDataAsync(reportID int, user *repository.User, dateFrom, dateTo string) {
	// Здесь следует реализовать фактический вызов API WB
	// Пока что, имитируем обработку

	time.Sleep(2 * time.Second) // Имитировать вызов API

	// Обновление статуса на "успех"
	h.statsRepo.UpdateStatus(reportID, 1, "")

	fmt.Printf("Report %d processed successfully for user %d\n", reportID, user.ID)
}

// Вспомогательная функция для получения пользователя из JWT-токена.
func (h *WBStatsHandler) getUserFromRequest(r *http.Request) (*repository.User, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("no authorization header")
	}

	// Извлечь токен
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Парсинг JWT-токена
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return h.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid username in token")
	}

	// Получить пользователя из репозитория
	return h.userRepo.GetByUsername(username)
}

func getStatusValue(ns sql.NullInt64) int {
	if ns.Valid {
		return int(ns.Int64)
	}
	return 0
}

func getStringValue(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func getNullInt64(v int) sql.NullInt64 {
	return sql.NullInt64{Int64: int64(v), Valid: true}
}

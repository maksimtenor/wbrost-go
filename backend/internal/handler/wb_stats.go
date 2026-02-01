package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"wbrost-go/internal/dto"
	"wbrost-go/internal/entity"
	"wbrost-go/internal/repository/stat"
	"wbrost-go/internal/repository/user"

	"github.com/golang-jwt/jwt/v4"
)

type WBStatsHandler struct {
	userRepo       *user.UserRepository
	wbStatsGetRepo *stat.WBStatsGetRepository
	statRepo       *stat.StatRepository
	analyticsRepo  *stat.AnalyticsRepository
	dashboardRepo  *stat.DashboardRepository
	jwtSecret      []byte
}

func NewWBStatsHandler(
	userRepo *user.UserRepository,
	wbStatsGetRepo *stat.WBStatsGetRepository,
	statRepo *stat.StatRepository,
	analyticsRepo *stat.AnalyticsRepository,
	dashboardRepo *stat.DashboardRepository,
	jwtSecret string) *WBStatsHandler {
	return &WBStatsHandler{
		userRepo:       userRepo,
		wbStatsGetRepo: wbStatsGetRepo,
		statRepo:       statRepo,
		analyticsRepo:  analyticsRepo,
		dashboardRepo:  dashboardRepo,
		jwtSecret:      []byte(jwtSecret),
	}
}

// GetStatDetail - GET /api/stat/details | Получение детальной статистики
func (h *WBStatsHandler) GetStatDetail(w http.ResponseWriter, r *http.Request) {
	// Получить пользователя по токену
	user, err := h.getUserFromRequest(r)
	if err != nil {
		respondWithJSON(w, http.StatusUnauthorized, dto.ErrorResponse{Error: "Unauthorized"})
		return
	}

	// Получить параметры дат из запроса(query params)
	dateFrom := r.URL.Query().Get("dateFrom")
	dateTo := r.URL.Query().Get("dateTo")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	//Валидация дат(если не указаны, можно использовать значения по умолчанию или вернуть ошибку)
	if dateFrom == "" || dateTo == "" {
		respondWithJSON(w, http.StatusBadRequest, dto.ErrorResponse{
			Error: "Parameters dateFrom and dateTo are required",
		})
		return
	}

	// Проверяем формат дат (опционально)
	if _, err := time.Parse("2006-01-02", dateFrom); err != nil {
		respondWithJSON(w, http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid dateFrom format, expected YYYY-MM-DD",
		})
		return
	}

	if _, err := time.Parse("2006-01-02", dateTo); err != nil {
		respondWithJSON(w, http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid dateTo format, expected YYYY-MM-DD",
		})
		return
	}

	// Параметры пагинации
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 20 // По умолчанию
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	// Получить детальные данные статистики
	statDetails, err := h.analyticsRepo.GetStatDetails(user.ID, dateFrom, dateTo, page, pageSize)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get stat details: " + err.Error(),
		})
		return
	}

	// Получить общее количество записей для пагинации
	totalCount, err := h.analyticsRepo.GetStatDetailsCount(user.ID, dateFrom, dateTo)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get total count: " + err.Error(),
		})
		return
	}

	// Получить итоговые суммы
	summaryStatDetails, err := h.analyticsRepo.GetStatSummary(user.ID, dateFrom, dateTo)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get stat summary: " + err.Error(),
		})
		return
	}

	// Формируем ответ
	response := map[string]interface{}{
		"data":    statDetails,
		"summary": summaryStatDetails,
		"taxes":   user.Taxes,
		"pagination": map[string]interface{}{
			"current_page": page,
			"page_size":    pageSize,
			"total_items":  totalCount,
			"total_pages":  int(math.Ceil(float64(totalCount) / float64(pageSize))),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetWBReports - GET /api/wb/stats | Получение списка репортов(заказов отчетов из бд)
func (h *WBStatsHandler) GetWBReports(w http.ResponseWriter, r *http.Request) {
	// Получить пользователя по токену
	user, err := h.getUserFromRequest(r)
	if err != nil {
		respondWithJSON(w, http.StatusUnauthorized, dto.ErrorResponse{Error: "Unauthorized"})
		return
	}

	// Получить отчеты для этого пользователя
	reports, err := h.wbStatsGetRepo.GetByUserID(user.ID)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to get reports"})
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
		respondWithJSON(w, http.StatusUnauthorized, dto.ErrorResponse{Error: "Unauthorized"})
		return
	}

	// Парсим запрос
	var req struct {
		DateFrom string `json:"dateFrom"`
		DateTo   string `json:"dateTo"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid request body"})
		return
	}

	// Валидация данных
	if req.DateFrom == "" || req.DateTo == "" {
		respondWithJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "DateFrom and DateTo are required"})
		return
	}

	// Создание репорта в бд
	stats := &entity.WBStatsGet{ // Меняем repository.WBStatsGet на entity.WBStatsGet
		UserID:   user.ID,
		Status:   getNullInt64(0), // 0 = в обработке
		DateFrom: req.DateFrom,
		DateTo:   req.DateTo,
	}

	if err := h.wbStatsGetRepo.Create(stats); err != nil {
		respondWithJSON(w, http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to create report: " + err.Error()})
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
func (h *WBStatsHandler) fetchWBDataAsync(reportID int, user *entity.Users, dateFrom, dateTo string) {
	// Здесь следует реализовать фактический вызов API WB
	// Пока что, имитируем обработку

	time.Sleep(2 * time.Second) // Имитировать вызов API

	// Обновление статуса на "успех"
	h.wbStatsGetRepo.UpdateStatus(reportID, 1, "")

	fmt.Printf("Report %d processed successfully for user %d\n", reportID, user.ID)
}

// GetDashboardStats - GET /api/dashboard/stats | Получение статистики для дашборда
func (h *WBStatsHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	// Получить пользователя по токену
	user, err := h.getUserFromRequest(r)
	if err != nil {
		respondWithJSON(w, http.StatusUnauthorized, dto.ErrorResponse{Error: "Unauthorized"})
		return
	}

	// Получить параметры дат (по умолчанию - текущий месяц)
	now := time.Now()
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastDay := firstDay.AddDate(0, 1, -1)

	dateFrom := r.URL.Query().Get("dateFrom")
	dateTo := r.URL.Query().Get("dateTo")

	if dateFrom == "" {
		dateFrom = firstDay.Format("2006-01-02")
	}
	if dateTo == "" {
		dateTo = lastDay.Format("2006-01-02")
	}

	// Получить статистику для дашборда
	stats, err := h.dashboardRepo.GetDashboardStats(user.ID, dateFrom, dateTo)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get dashboard stats: " + err.Error(),
		})
		return
	}

	// Получить данные для графиков
	chartData, err := h.dashboardRepo.GetChartData(user.ID, dateFrom, dateTo)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get chart data: " + err.Error(),
		})
		return
	}

	// Получить данные по месяцам
	monthlyRevenue, err := h.dashboardRepo.GetMonthlyRevenue(user.ID)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get monthly revenue: " + err.Error(),
		})
		return
	}

	// Формируем ответ
	response := map[string]interface{}{
		"stats":           stats,
		"charts":          chartData,
		"monthly_revenue": monthlyRevenue, // Добавляем новые данные
		"period": map[string]string{
			"dateFrom": dateFrom,
			"dateTo":   dateTo,
		},
	}

	respondWithJSON(w, http.StatusOK, response)
}

// Вспомогательная функция для получения пользователя из JWT-токена.
func (h *WBStatsHandler) getUserFromRequest(r *http.Request) (*entity.Users, error) {
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

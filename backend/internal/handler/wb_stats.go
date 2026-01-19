package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
	"time"
	"wbrost-go/internal/repository"
)

type WBStatsHandler struct {
	userRepo  *repository.UserRepository
	statsRepo *repository.WBStatsRepository
	jwtSecret []byte
}

func NewWBStatsHandler(userRepo *repository.UserRepository, statsRepo *repository.WBStatsRepository, jwtSecret string) *WBStatsHandler {
	return &WBStatsHandler{
		userRepo:  userRepo,
		statsRepo: statsRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

// GetWBStats - GET /api/wb/stats
func (h *WBStatsHandler) GetWBStats(w http.ResponseWriter, r *http.Request) {
	// Get user from token
	user, err := h.getUserFromRequest(r)
	if err != nil {
		respondWithJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	// Get reports for this user
	reports, err := h.statsRepo.GetByUserID(user.ID)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Failed to get reports"})
		return
	}

	// Prepare response
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

// CreateWBReport - POST /api/wb/stats
func (h *WBStatsHandler) CreateWBReport(w http.ResponseWriter, r *http.Request) {
	// Get user from token
	user, err := h.getUserFromRequest(r)
	if err != nil {
		respondWithJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	// Parse request
	var req struct {
		DateFrom string `json:"dateFrom"`
		DateTo   string `json:"dateTo"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	// Validate dates
	if req.DateFrom == "" || req.DateTo == "" {
		respondWithJSON(w, http.StatusBadRequest, ErrorResponse{Error: "DateFrom and DateTo are required"})
		return
	}

	// Create report in DB
	stats := &repository.WBStatsGet{
		UserID:   user.ID,
		Status:   getNullInt64(0), // 0 = в обработке
		DateFrom: req.DateFrom,
		DateTo:   req.DateTo,
	}

	if err := h.statsRepo.Create(stats); err != nil {
		respondWithJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Failed to create report: " + err.Error()})
		return
	}

	// Start async task to fetch data from WB API
	//go h.fetchWBDataAsync(stats.ID, user, req.DateFrom, req.DateTo)

	// Return success
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"id":      stats.ID,
		"success": true,
		"message": "Отчет поставлен в очередь на формирование",
	})
}

// Async function to fetch data from WB
func (h *WBStatsHandler) fetchWBDataAsync(reportID int, user *repository.User, dateFrom, dateTo string) {
	// Here you would implement actual WB API call
	// For now, simulate processing

	time.Sleep(2 * time.Second) // Simulate API call

	// Update status to success
	h.statsRepo.UpdateStatus(reportID, 1, "")

	fmt.Printf("Report %d processed successfully for user %d\n", reportID, user.ID)
}

// Helper function to get user from JWT token
func (h *WBStatsHandler) getUserFromRequest(r *http.Request) (*repository.User, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("no authorization header")
	}

	// Extract token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse JWT token
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

	// Get user from repository
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

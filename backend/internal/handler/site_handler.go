package handler

import (
	"encoding/json"
	"net/http"
	"wbrost-go/internal/service"
)

type SiteHandler struct {
	siteService service.SiteService
}

func NewSiteHandler(siteService service.SiteService) *SiteHandler {
	return &SiteHandler{siteService: siteService}
}

// Homepage handler - возвращает JSON для Vue.js
func (h *SiteHandler) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Получаем userID из сессии/JWT (пока заглушка)
	userID := 0 // По умолчанию гость

	// Проверяем авторизацию (пример с cookie)
	cookie, err := r.Cookie("user_id")
	if err == nil && cookie.Value != "" {
		// Здесь должен быть парсинг JWT или получение из сессии
		// userID = parseUserID(cookie.Value)
	}

	dashboardData, err := h.siteService.GetDashboardData(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(dashboardData)
}

// Информационные страницы (отдаём HTML через Vue.js)
func (h *SiteHandler) Info(w http.ResponseWriter, r *http.Request) {
	// Эти страницы будут обрабатываться Vue Router на фронтенде
	// Go только отдаёт статику или редиректит
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *SiteHandler) Privacy(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *SiteHandler) Terms(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *SiteHandler) Donation(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

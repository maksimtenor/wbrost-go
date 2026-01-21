package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"wbrost-go/internal/entity"
	"wbrost-go/internal/repository"

	"github.com/golang-jwt/jwt/v4"
)

type WBArticlesHandler struct {
	userRepo        *repository.UserRepository
	articlesGetRepo *repository.WBArticlesGetRepository
	articleRepo     *repository.WBArticleRepository
	jwtSecret       []byte
}

func NewWBArticlesHandler(
	userRepo *repository.UserRepository,
	articlesGetRepo *repository.WBArticlesGetRepository,
	articleRepo *repository.WBArticleRepository,
	jwtSecret string,
) *WBArticlesHandler {
	return &WBArticlesHandler{
		userRepo:        userRepo,
		articlesGetRepo: articlesGetRepo,
		articleRepo:     articleRepo,
		jwtSecret:       []byte(jwtSecret),
	}
}

// GetArticles - GET /api/articles | Получение списка карточек товаров
func (h *WBArticlesHandler) GetArticles(w http.ResponseWriter, r *http.Request) {
	// Получить пользователя по токену
	user, err := h.getUserFromRequest(r)
	if err != nil {
		respondWithJSON(w, http.StatusUnauthorized, entity.ErrorResponse{Error: "Unauthorized"})
		return
	}

	// Параметры запроса
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	searchQuery := r.URL.Query().Get("search")

	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 10 // По умолчанию
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	var articles []entity.WBArticleDB
	var totalCount int

	if searchQuery != "" {
		// Поиск по запросу
		articles, err = h.articleRepo.SearchArticles(user.ID, searchQuery, page, pageSize)
		if err != nil {
			respondWithJSON(w, http.StatusInternalServerError, entity.ErrorResponse{
				Error: "Failed to search articles: " + err.Error(),
			})
			return
		}
		// Временное решение для поиска
		totalCount = len(articles)
	} else {
		// Получение всех карточек с пагинацией
		articles, err = h.articleRepo.GetByUserID(user.ID, page, pageSize)
		if err != nil {
			respondWithJSON(w, http.StatusInternalServerError, entity.ErrorResponse{
				Error: "Failed to get articles: " + err.Error(),
			})
			return
		}

		// Получить общее количество
		totalCount, err = h.articleRepo.GetCountByUserID(user.ID)
		if err != nil {
			respondWithJSON(w, http.StatusInternalServerError, entity.ErrorResponse{
				Error: "Failed to get total count: " + err.Error(),
			})
			return
		}
	}

	// Форматировать ответ
	response := make([]map[string]interface{}, len(articles))
	for i, article := range articles {
		// Генерация URL фото
		photoURL := ""
		if article.Photo.Valid && article.Photo.String != "" {
			photoURL = article.Photo.String
		} else if article.Articule != "" {
			// Генерация URL по аналогии со статистикой
			photoURL = h.generatePhotoURL(article.Articule)
		}

		// Форматирование дат
		createdDate := ""
		if article.Created.Valid {
			createdDate = article.Created.Time.Format("2006-01-02")
		}

		updatedDate := ""
		if article.Updated.Valid {
			updatedDate = article.Updated.Time.Format("2006-01-02")
		}

		response[i] = map[string]interface{}{
			"id":         article.ID,
			"articule":   article.Articule,
			"name":       getStringValue(article.Name),
			"photo":      photoURL,
			"cost_price": getStringValue(article.CostPrice),
			"created":    createdDate,
			"updated":    updatedDate,
			"rus_size":   getStringValue(article.RusSize),
			"eu_size":    getStringValue(article.EuSize),
			"chrt_id":    getIntValue(article.ChrtID),
			"barcode":    getStringValue(article.Barcode),
		}
	}

	// Полный ответ с пагинацией
	fullResponse := map[string]interface{}{
		"data": response,
		"pagination": map[string]interface{}{
			"current_page": page,
			"page_size":    pageSize,
			"total_items":  totalCount,
			"total_pages":  (totalCount + pageSize - 1) / pageSize,
		},
	}

	respondWithJSON(w, http.StatusOK, fullResponse)
}

// CreateArticlesRequest - POST /api/articles/request | Запрос обновления карточек товаров
func (h *WBArticlesHandler) CreateArticlesRequest(w http.ResponseWriter, r *http.Request) {
	// Получить пользователя по токену
	user, err := h.getUserFromRequest(r)
	if err != nil {
		respondWithJSON(w, http.StatusUnauthorized, entity.ErrorResponse{Error: "Unauthorized"})
		return
	}

	// Проверяем наличие WB ключа
	if !user.WbKey.Valid || user.WbKey.String == "" {
		respondWithJSON(w, http.StatusBadRequest, entity.ErrorResponse{
			Error: "WB ключ не указан. Добавьте ключ в настройках профиля.",
		})
		return
	}

	// Создаем запись в wb_articles_get
	articleRequest := &entity.WBArticlesGet{
		UserID: user.ID,
		Status: getNullInt64(entity.ArticlesStatusWait),
	}

	if err := h.articlesGetRepo.Create(articleRequest); err != nil {
		respondWithJSON(w, http.StatusInternalServerError, entity.ErrorResponse{
			Error: "Failed to create request: " + err.Error(),
		})
		return
	}

	// Запускаем асинхронную обработку
	// go h.fetchArticlesAsync(articleRequest.ID, user)

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"id":      articleRequest.ID,
		"success": true,
		"message": "Запрос на обновление карточек товаров поставлен в очередь",
	})
}

// UpdateCostPrice - POST /api/articles/cost-price | Обновление себестоимости
func (h *WBArticlesHandler) UpdateCostPrice(w http.ResponseWriter, r *http.Request) {
	// Получить пользователя по токену
	user, err := h.getUserFromRequest(r)
	if err != nil {
		respondWithJSON(w, http.StatusUnauthorized, entity.ErrorResponse{Error: "Unauthorized"})
		return
	}

	// Парсим запрос
	var req struct {
		Articule  string `json:"articule"`
		CostPrice string `json:"cost_price"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithJSON(w, http.StatusBadRequest, entity.ErrorResponse{Error: "Invalid request body"})
		return
	}

	// Валидация
	if req.Articule == "" {
		respondWithJSON(w, http.StatusBadRequest, entity.ErrorResponse{Error: "Articule is required"})
		return
	}

	// Обновляем себестоимость
	if err := h.articleRepo.UpdateCostPrice(user.ID, req.Articule, req.CostPrice); err != nil {
		respondWithJSON(w, http.StatusInternalServerError, entity.ErrorResponse{
			Error: "Failed to update cost price: " + err.Error(),
		})
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Себестоимость обновлена",
	})
}

// Вспомогательная функция для генерации URL фото
func (h *WBArticlesHandler) generatePhotoURL(articule string) string {
	if articule == "" {
		return ""
	}

	// Парсим артикул в число
	nmID, err := strconv.ParseInt(articule, 10, 64)
	if err != nil {
		return ""
	}

	// Формируем URL как в статистике
	basketNum := nmID % 100
	vol := nmID / 1000000
	part := nmID / 1000

	return fmt.Sprintf("https://basket-%d.wbbasket.ru/vol%d/part%d/%d/images/big/1.webp",
		basketNum, vol, part, nmID)
}

// Вспомогательная функция для получения пользователя из JWT
func (h *WBArticlesHandler) getUserFromRequest(r *http.Request) (*repository.User, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("no authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
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

	return h.userRepo.GetByUsername(username)
}

// Вспомогательные функции для работы с NULL значениями

func getIntValue(ni sql.NullInt64) int64 {
	if ni.Valid {
		return ni.Int64
	}
	return 0
}

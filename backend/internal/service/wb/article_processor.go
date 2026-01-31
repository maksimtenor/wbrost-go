package wb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"wbrost-go/internal/api/wb"
	"wbrost-go/internal/entity"
	"wbrost-go/internal/repository/user"
)

// ProcessPendingArticles обрабатывает все ожидающие запросы на получение карточек
func (s *WBService) ProcessPendingArticles() error {
	articles, err := s.articlesGetRepo.GetPendingArticles()
	if err != nil {
		return fmt.Errorf("failed to get pending articles: %w", err)
	}

	if len(articles) == 0 {
		fmt.Println("No pending articles requests found")
		return nil
	}

	for _, articleReq := range articles {
		fmt.Printf("Processing articles request ID: %d for user %d\n", articleReq.ID, articleReq.UserID)

		user, err := s.userRepo.GetByID(articleReq.UserID)
		if err != nil {
			s.updateArticleStatus(&articleReq, entity.ArticlesStatusError, "User not found")
			continue
		}

		if !user.WbKey.Valid || user.WbKey.String == "" {
			s.updateArticleStatus(&articleReq, entity.ArticlesStatusError, "WB key not found")
			continue
		}

		// Обрабатываем запрос
		result := s.processArticleRequest(&articleReq, user)

		if result.Status {
			s.updateArticleStatus(&articleReq, entity.ArticlesStatusSuccess, result.Error)
		} else {
			status := entity.ArticlesStatusError
			if result.Retake {
				status = entity.ArticlesStatusWait
			}
			s.updateArticleStatus(&articleReq, status, result.Error)
		}
	}

	return nil
}

func (s *WBService) updateArticleStatus(article *entity.WBArticlesGet, status int, errorMsg string) {
	err := s.articlesGetRepo.UpdateStatus(article.ID, status, errorMsg)
	if err != nil {
		fmt.Printf("Failed to update article %d status: %v\n", article.ID, err)
	} else {
		fmt.Printf("Article request %d updated to status %d\n", article.ID, status)
	}
}

func (s *WBService) processArticleRequest(articleReq *entity.WBArticlesGet, user *user.User) ProcessResult {
	// Получаем данные карточек от WB API
	articlesData, err := s.getWBArticles(user)
	if err != nil {
		return ProcessResult{
			Status: false,
			Error:  fmt.Sprintf("Failed to get WB articles: %v", err),
			Retake: false,
		}
	}

	// Обрабатываем и сохраняем данные
	success, message := s.saveArticles(articlesData, user.ID)

	return ProcessResult{
		Status: success,
		Error:  message,
		Retake: false,
	}
}

func (s *WBService) getWBArticles(user *user.User) ([]wb.Article, error) {
	if !user.WbKey.Valid || user.WbKey.String == "" {
		return nil, fmt.Errorf("токен WB не указан для пользователя %d", user.ID)
	}

	token := user.WbKey.String
	client := wb.NewWBClient(token)

	// Проверяем токен
	isValid, err := client.CheckToken()
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки токена: %v", err)
	}

	if !isValid {
		return nil, fmt.Errorf("токен недействителен или истек")
	}

	// Получаем карточки товаров
	return s.fetchArticlesFromWB(client)
}

func (s *WBService) fetchArticlesFromWB(client *wb.Client) ([]wb.Article, error) {
	var allCards []wb.Article
	var cursorUpdatedAt string
	var cursorNmID int

	limit := 100 // Максимальный лимит за один запрос
	totalProcessed := 0

	for {
		// Формируем тело запроса
		request := wb.ArticleRequest{
			Settings: wb.ArticleRequestSettings{
				Cursor: wb.ArticleRequestCursor{
					Limit: limit,
				},
				Filter: struct {
					WithPhoto int `json:"withPhoto"`
				}{
					WithPhoto: -1,
				},
			},
		}

		// Добавляем курсор, если он есть (для пагинации)
		if cursorUpdatedAt != "" && cursorNmID > 0 {
			request.Settings.Cursor.UpdatedAt = cursorUpdatedAt
			request.Settings.Cursor.NmID = cursorNmID
		}

		jsonBody, err := json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}

		url := wb.URLCardsList()
		req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonBody)))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Authorization", client.Token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %w", err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != 200 {
			fmt.Printf("Response body: %s\n", string(body))
			return nil, fmt.Errorf("WB API error: status %d", resp.StatusCode)
		}

		// Парсим ответ
		var response wb.ArticleResponse
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		// Добавляем полученные карточки
		allCards = append(allCards, response.Cards...)
		totalProcessed += len(response.Cards)

		fmt.Printf("Получено %d карточек (всего: %d)\n", len(response.Cards), totalProcessed)

		// Проверяем, нужно ли продолжать пагинацию
		if len(response.Cards) < limit || response.Cursor.Total < limit {
			fmt.Printf("Получены все карточки. Всего: %d\n", totalProcessed)
			break
		}

		// Обновляем курсор для следующего запроса
		cursorUpdatedAt = response.Cursor.UpdatedAt
		cursorNmID = response.Cursor.NmID

		// Небольшая задержка между запросами, чтобы не превысить лимиты
		time.Sleep(100 * time.Millisecond)
	}

	return allCards, nil
}

func (s *WBService) saveArticles(cards []wb.Article, userID int) (bool, string) {
	if len(cards) == 0 {
		return false, "No articles data received"
	}

	countSaved := 0
	countUnsaved := 0
	countTotal := len(cards)

	fmt.Printf("📦 Получено %d карточек товаров от WB API\n", countTotal)

	for _, card := range cards {
		// Создаем запись для базы данных (используем WBArticleDB)
		article := &entity.WBArticles{
			UserID:    userID,
			Articule:  strconv.Itoa(card.NmID),
			Created:   sql.NullTime{Time: time.Now(), Valid: true},
			Updated:   sql.NullTime{Time: time.Now(), Valid: true},
			UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		}

		// Заполняем основные поля
		if card.Title != "" {
			article.Name = sql.NullString{String: card.Title, Valid: true}
		}

		if card.VendorCode != "" {
			article.InternalID = sql.NullString{String: card.VendorCode, Valid: true}
		}

		if card.NmUUID != "" {
			article.InternalID = sql.NullString{String: card.NmUUID, Valid: true}
		}

		// Берем первую фотографию, если есть
		if len(card.Photos) > 0 && card.Photos[0].Big != "" {
			article.Photo = sql.NullString{String: card.Photos[0].Big, Valid: true}
		}

		// Обрабатываем размеры
		if len(card.Sizes) > 0 {
			size := card.Sizes[0]
			if size.TechSize != "" {
				article.EuSize = sql.NullString{String: size.TechSize, Valid: true}
			}
			if size.WbSize != "" {
				article.RusSize = sql.NullString{String: size.WbSize, Valid: true}
			}
			if size.ChrtID != 0 {
				article.ChrtID = sql.NullInt64{Int64: int64(size.ChrtID), Valid: true}
			}
			if len(size.Skus) > 0 {
				article.Barcode = sql.NullString{String: strings.Join(size.Skus, ", "), Valid: true}
			}
		}

		// Сохраняем в БД
		if err := s.articleRepo.CreateOrUpdate(article); err != nil {
			fmt.Printf("Error saving article %d: %v\n", card.NmID, err)
			countUnsaved++
			continue
		}

		countSaved++
	}

	message := fmt.Sprintf("Total: %d, Saved: %d, Not saved: %d", countTotal, countSaved, countUnsaved)
	success := countSaved > 0

	fmt.Printf("✅ Результат обработки карточек: %s\n", message)
	return success, message
}

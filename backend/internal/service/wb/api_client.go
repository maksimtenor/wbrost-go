package wb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"wbrost-go/internal/api/wb"
	"wbrost-go/internal/entity"
)

func (s *WBService) getWBData(order *entity.WBStatsGet, user *entity.Users) ([]interface{}, error) {
	if !user.WbKey.Valid || user.WbKey.String == "" {
		return nil, fmt.Errorf("токен WB не указан для пользователя %d", user.ID)
	}

	token := user.WbKey.String

	// Проверяем формат токена
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("неверный формат токена. Ожидается JWT токен")
	}

	client := wb.NewWBClient(token)

	// Сначала проверяем токен через WB API
	fmt.Println("🔐 Проверка токена...")
	isValid, err := client.CheckToken()
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки токена: %v", err)
	}

	if !isValid {
		return nil, fmt.Errorf("токен недействителен или истек")
	}

	fmt.Println("✅ Токен валиден")

	var allData []interface{}

	dateFrom, err := time.Parse("2006-01-02", order.DateFrom)
	if err != nil {
		return nil, fmt.Errorf("invalid date_from format: %w", err)
	}

	dateTo, err := time.Parse("2006-01-02", order.DateTo)
	if err != nil {
		return nil, fmt.Errorf("invalid date_to format: %w", err)
	}

	// Рассчитываем длительность периода
	days := int(dateTo.Sub(dateFrom).Hours() / 24)
	fmt.Printf("📅 Период: %s - %s (%d дней)\n",
		order.DateFrom, order.DateTo, days)

	// Для периодов больше 90 дней разбиваем на кварталы (3 месяца)
	if days > 90 {
		fmt.Printf("📦 Большой период (%d дней), разбиваем на кварталы\n", days)

		currentStart := dateFrom
		quarterCount := 0

		for currentStart.Before(dateTo) {
			quarterCount++
			// Конец квартала = +3 месяца -1 день
			currentEnd := currentStart.AddDate(0, 3, -1)
			if currentEnd.After(dateTo) {
				currentEnd = dateTo
			}

			fmt.Printf("\n🔍 Квартал %d: %s - %s\n",
				quarterCount,
				currentStart.Format("2006-01-02"),
				currentEnd.Format("2006-01-02"))

			// Определяем версию API
			useNewAPI := true
			if currentStart.Before(time.Date(2024, 1, 29, 0, 0, 0, 0, time.UTC)) {
				useNewAPI = false
				fmt.Println("📊 Используем старую версию API (до 29.01.2024)")
			} else {
				fmt.Println("📊 Используем новую версию API (после 29.01.2024)")
			}

			data, err := s.getReportByPeriod(
				client,
				currentStart.Format("2006-01-02"),
				currentEnd.Format("2006-01-02"),
				useNewAPI,
			)

			if err != nil {
				return nil, fmt.Errorf("ошибка за период %s-%s: %w",
					currentStart.Format("2006-01-02"),
					currentEnd.Format("2006-01-02"), err)
			}

			if len(data) > 0 {
				allData = append(allData, data...)
				fmt.Printf("✅ Получено %d записей за квартал\n", len(data))
			} else {
				fmt.Println("ℹ️  Нет данных за этот квартал")
			}

			// Переход к следующему кварталу
			currentStart = currentEnd.AddDate(0, 0, 1)

			// Пауза между кварталами
			if currentStart.Before(dateTo) {
				fmt.Println("⏸️  Пауза 3 секунды перед следующим кварталом...")
				time.Sleep(3 * time.Second)
			}
		}

		fmt.Printf("\n✅ Все кварталы обработаны. Всего кварталов: %d\n", quarterCount)
	} else {
		// Малый период - запрашиваем целиком
		fmt.Println("📦 Малый период, запрашиваем целиком")

		useNewAPI := true
		if dateFrom.Before(time.Date(2024, 1, 29, 0, 0, 0, 0, time.UTC)) {
			useNewAPI = false
			fmt.Println("📊 Используем старую версию API (до 29.01.2024)")
		} else {
			fmt.Println("📊 Используем новую версию API (после 29.01.2024)")
		}

		data, err := s.getReportByPeriod(
			client,
			dateFrom.Format("2006-01-02"),
			dateTo.Format("2006-01-02"),
			useNewAPI,
		)

		if err != nil {
			return nil, err
		}

		allData = data
	}

	fmt.Printf("\n🎉 Всего получено записей: %d\n", len(allData))
	return allData, nil
}

func (s *WBService) safeRequest(client *wb.Client, url string) (*http.Response, error) {
	maxRetries := 5

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Ждем разрешения от rate limiter
		if err := s.rateLimiter.Wait(); err != nil {
			return nil, fmt.Errorf("rate limiter error: %w", err)
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Authorization", client.Token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to execute request: %w", err)
		}

		// ОБРАБАТЫВАЕМ ЗАГОЛОВКИ ДЛЯ RATE LIMITING
		s.rateLimiter.ProcessHeaders(resp.Header, resp.StatusCode)

		// Проверяем статус код
		switch resp.StatusCode {
		case 200:
			return resp, nil

		case 429: // Too Many Requests
			// Rate limiter уже обработал заголовки в ProcessHeaders
			resp.Body.Close()

			// Не нужно ждать здесь - rate limiter.Wait() уже будет ждать нужное время
			// Просто логируем и продолжаем цикл
			fmt.Printf("🔄 429 - Rate limiter обработан, пробуем снова (попытка %d/%d)\n",
				attempt+1, maxRetries)
			continue

		case 404:
			// Обрабатываем 404 ошибку отдельно
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			// Парсим ответ для получения деталей
			var errorResp map[string]interface{}
			if err := json.Unmarshal(body, &errorResp); err == nil {
				if title, ok := errorResp["title"].(string); ok && title == "path not found" {
					// Это известная ошибка API - возвращаем как обычную ошибку
					return nil, fmt.Errorf("WB API error 404: %v", errorResp)
				}
			}

			// Возвращаем ответ для дальнейшей обработки
			return &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}, nil

		default:
			// Для других ошибок возвращаем ответ как есть
			return resp, nil
		}
	}

	return nil, fmt.Errorf("max retries (%d) exceeded", maxRetries)
}
func (s *WBService) getReportByPeriod(client *wb.Client, dateFrom, dateTo string, useNewAPI bool) ([]interface{}, error) {
	// Используем пагинированную версию
	return s.getReportByPeriodWithPagination(client, dateFrom, dateTo, useNewAPI)
}
func (s *WBService) getReportByPeriodWithPagination(client *wb.Client, dateFrom, dateTo string, useNewAPI bool) ([]interface{}, error) {
	var allData []interface{}
	var lastRrdID int64 = 0
	maxPages := 50 // Максимальное количество страниц пагинации

	for page := 1; page <= maxPages; page++ {
		var url string
		if useNewAPI {
			url = fmt.Sprintf("%s?dateFrom=%s&dateTo=%s&rrdid=%d&limit=100000",
				wb.URLFor(wb.DetailsV5), dateFrom, dateTo, lastRrdID)
		} else {
			url = fmt.Sprintf("%s?dateFrom=%s&dateTo=%s&rrdid=%d&limit=100000",
				wb.URLFor(wb.DetailsV1), dateFrom, dateTo, lastRrdID)
		}

		fmt.Printf("📄 Страница %d: запрос данных с rrdid=%d\n", page, lastRrdID)

		// Ждем разрешения от rate limiter
		if err := s.rateLimiter.Wait(); err != nil {
			return nil, fmt.Errorf("rate limiter error: %w", err)
		}

		resp, err := s.safeRequest(client, url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		// Проверяем статус код
		switch resp.StatusCode {
		case 200:
			var data []interface{}
			if err := json.Unmarshal(body, &data); err != nil {
				return nil, fmt.Errorf("failed to parse response: %w", err)
			}

			if len(data) == 0 {
				// Больше данных нет - завершаем пагинацию
				fmt.Printf("✅ Пагинация завершена. Всего страниц: %d, записей: %d\n", page-1, len(allData))
				return allData, nil
			}

			allData = append(allData, data...)

			// Получаем последний rrd_id для следующего запроса
			if lastItem, ok := data[len(data)-1].(map[string]interface{}); ok {
				if rrdID, ok := lastItem["rrd_id"].(float64); ok {
					lastRrdID = int64(rrdID)
					fmt.Printf("📊 Страница %d: получено %d записей (всего: %d), следующий rrd_id: %d\n",
						page, len(data), len(allData), lastRrdID)
				} else {
					// Если нет rrd_id, значит это последняя страница
					fmt.Printf("✅ Последняя страница. Всего записей: %d\n", len(allData))
					return allData, nil
				}
			} else {
				// Не удалось получить последний элемент
				fmt.Printf("✅ Пагинация завершена. Всего записей: %d\n", len(allData))
				return allData, nil
			}

		case 204:
			// Нет данных - завершение пагинации
			fmt.Printf("✅ Пагинация завершена (204 No Content). Всего страниц: %d, записей: %d\n",
				page-1, len(allData))
			return allData, nil

		case 429:
			// Rate limit - обрабатывается в safeRequest
			// Ждем и пробуем снова
			fmt.Println("🔄 429 - Rate limit, повторяем запрос...")
			continue

		default:
			// Другие ошибки
			fmt.Printf("❌ Ошибка API: статус %d\n", resp.StatusCode)
			fmt.Printf("Ответ: %s\n", string(body))
			return allData, fmt.Errorf("API error: status %d", resp.StatusCode)
		}

		// Короткая пауза между страницами
		if page < maxPages {
			time.Sleep(500 * time.Millisecond)
		}
	}

	fmt.Printf("⚠️ Достигнут лимит страниц (%d). Всего записей: %d\n", maxPages, len(allData))
	return allData, nil
}

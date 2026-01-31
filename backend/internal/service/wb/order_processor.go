package wb

import (
	"fmt"
	"strings"
	"wbrost-go/internal/entity"
	"wbrost-go/internal/repository/user"
)

// ProcessPendingOrders обрабатывает все ожидающие заказы
func (s *WBService) ProcessPendingOrders() error {
	orders, err := s.statsGetRepo.GetPendingOrders()
	if err != nil {
		return fmt.Errorf("failed to get pending orders: %w", err)
	}

	if len(orders) == 0 {
		fmt.Println("No pending orders found")
		return nil
	}

	for _, order := range orders {
		fmt.Printf("Processing order ID: %d for user %d\n", order.ID, order.UserID)

		user, err := s.userRepo.GetByID(order.UserID)
		if err != nil {
			s.updateOrderStatus(&order, entity.StatusError, "User not found")
			continue
		}

		if !user.WbKey.Valid || user.WbKey.String == "" {
			s.updateOrderStatus(&order, entity.StatusError, "WB key not found")
			continue
		}

		// Проверяем формат токена
		if !s.isValidTokenFormat(user.WbKey.String) {
			s.updateOrderStatus(&order, entity.StatusError, "Invalid WB token format")
			continue
		}

		// Обрабатываем заказ
		result := s.processOrder(&order, user)

		if result.Status {
			s.updateOrderStatus(&order, entity.StatusSuccess, result.Error)
		} else {
			status := entity.StatusError
			if result.Retake {
				status = entity.StatusWait
			}
			s.updateOrderStatus(&order, status, result.Error)
		}
	}

	return nil
}

func (s *WBService) updateOrderStatus(order *entity.WBStatsGet, status int, errorMsg string) {
	err := s.statsGetRepo.UpdateOrderStatus(order.ID, status, errorMsg)
	if err != nil {
		fmt.Printf("Failed to update order %d status: %v\n", order.ID, err)
	} else {
		fmt.Printf("Order %d updated to status %d\n", order.ID, status)
	}
}

func (s *WBService) processOrder(order *entity.WBStatsGet, user *user.User) ProcessResult {
	// Получаем данные от WB API
	reportData, err := s.getWBData(order, user)
	if err != nil {
		// Проверяем, является ли ошибка 404 "path not found"
		if strings.Contains(err.Error(), "path not found") {
			return ProcessResult{
				Status: false,
				Error:  "API endpoint not found (404). Возможно, у пользователя нет доступа к этому отчету.",
				Retake: false,
			}
		}

		return ProcessResult{
			Status: false,
			Error:  fmt.Sprintf("Failed to get WB data: %v", err),
			Retake: false,
		}
	}

	// Проверяем на лимит запросов
	if len(reportData) > 0 {
		if report, ok := reportData[0].(map[string]interface{}); ok {
			if title, ok := report["title"].(string); ok && title == entity.TooManyRequests {
				return ProcessResult{
					Status: false,
					Error:  entity.TooManyRequests,
					Retake: true,
				}
			}
		}
	}

	// Обрабатываем и сохраняем данные
	success, message := s.saveStats(reportData, user.ID)

	return ProcessResult{
		Status: success,
		Error:  message,
		Retake: false,
	}
}

func (s *WBService) isValidTokenFormat(token string) bool {
	if token == "" {
		return false
	}

	// JWT токен состоит из 3 частей, разделенных точками
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}

	// Каждая часть должна быть не пустой
	for _, part := range parts {
		if part == "" {
			return false
		}
	}

	return true
}

package service

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"wbrost-go/internal/entity"
	"wbrost-go/internal/repository"
	"wbrost-go/internal/wbapi"
)

type WBService struct {
	userRepo        *repository.UserRepository
	statsGetRepo    *repository.WBStatsGetRepository
	statRepo        *repository.StatRepository
	articlesGetRepo *repository.WBArticlesGetRepository
	articleRepo     *repository.WBArticleRepository
}

func NewWBService(
	userRepo *repository.UserRepository,
	statsGetRepo *repository.WBStatsGetRepository,
	statRepo *repository.StatRepository,
	articlesGetRepo *repository.WBArticlesGetRepository,
	articleRepo *repository.WBArticleRepository,
) *WBService {
	return &WBService{
		userRepo:        userRepo,
		statsGetRepo:    statsGetRepo,
		statRepo:        statRepo,
		articlesGetRepo: articlesGetRepo,
		articleRepo:     articleRepo,
	}
}

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

type ProcessResult struct {
	Status bool
	Error  string
	Retake bool
}

func (s *WBService) processOrder(order *entity.WBStatsGet, user *repository.User) ProcessResult {
	// Получаем данные от WB API
	reportData, err := s.getWBData(order, user)
	if err != nil {
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

func (s *WBService) getWBData(order *entity.WBStatsGet, user *repository.User) ([]interface{}, error) {
	if !user.WbKey.Valid || user.WbKey.String == "" {
		return nil, fmt.Errorf("токен WB не указан для пользователя %d", user.ID)
	}

	token := user.WbKey.String

	// Проверяем формат токена
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("неверный формат токена. Ожидается JWT токен")
	}

	client := wbapi.NewWBClient(token)

	// Сначала проверяем токен через WB API
	isValid, err := client.CheckToken()
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки токена: %v", err)
	}

	if !isValid {
		return nil, fmt.Errorf("токен недействителен или истек")
	}

	// Разбиваем период на интервалы по 30 дней (как в PHP)
	var allData []interface{}

	dateFrom, err := time.Parse("2006-01-02", order.DateFrom)
	if err != nil {
		return nil, fmt.Errorf("invalid date_from format: %w", err)
	}

	dateTo, err := time.Parse("2006-01-02", order.DateTo)
	if err != nil {
		return nil, fmt.Errorf("invalid date_to format: %w", err)
	}

	// Проверяем даты
	now := time.Now()
	if dateFrom.After(now) {
		return nil, fmt.Errorf("date_from (%s) в будущем", order.DateFrom)
	}

	if dateTo.After(now) {
		fmt.Printf("⚠️ date_to (%s) в будущем, используем текущую дату\n", order.DateTo)
		dateTo = now
	}

	swapDate, _ := time.Parse("2006-01-02", "2024-01-29")

	for current := dateFrom; current.Before(dateTo) || current.Equal(dateTo); {
		// Определяем конечную дату для интервала (макс 30 дней)
		endDate := current.AddDate(0, 0, 29)
		if endDate.After(dateTo) {
			endDate = dateTo
		}

		// Проверяем, какую версию API использовать
		useNewAPI := true
		if current.Before(swapDate) {
			useNewAPI = false
			if endDate.After(swapDate) {
				endDate = swapDate
			}
		}

		// Получаем данные за интервал
		data, err := s.getReportByPeriod(
			client,
			current.Format("2006-01-02"),
			endDate.Format("2006-01-02"),
			useNewAPI,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to get report for period %s - %s: %w",
				current.Format("2006-01-02"), endDate.Format("2006-01-02"), err)
		}

		if data != nil {
			allData = append(allData, data...)
		}

		// Переходим к следующему интервалу
		current = endDate.AddDate(0, 0, 1)
	}

	return allData, nil
}

// В service/wb_service.go добавьте:
func (s *WBService) convertSupplierOperName(supplierName interface{}) int64 {
	if supplierName == nil {
		return 0
	}

	var name string
	switch v := supplierName.(type) {
	case string:
		name = v
	case float64:
		name = strconv.FormatFloat(v, 'f', -1, 64)
	case int64:
		name = strconv.FormatInt(v, 10)
	default:
		name = fmt.Sprintf("%v", v)
	}

	// Логика как в Yii2 Stat::getSuplierType()
	switch name {
	case "Продажа":
		return 1
	case "Возврат":
		return 2
	case "Логистика":
		return 3
	case "Удержание":
		return 4
	case "Штраф":
		return 5
	case "Хранение":
		return 6
	case "Коррекция продаж":
		return 7
	case "Авансовая оплата за товар без движения":
		return 8
	case "Пересчет хранения":
		return 9
	case "Пересчет платной приемки":
		return 10
	case "Коррекция логистики":
		return 11
	case "Корректировка эквайринга":
		return 12
	case "Компенсация ущерба":
		return 13
	case "Компенсация потерянного товара":
		return 14
	case "Компенсация брака":
		return 15
	case "Добровольная компенсация при возврате":
		return 16
	case "Компенсация подмененного товара":
		return 17
	case "Возмещение издержек по перевозке/по складским операциям с товаром":
		return 18
	default:
		// Пробуем преобразовать как число
		if num, err := strconv.ParseInt(name, 10, 64); err == nil {
			return num
		}
		return 0
	}
}

func (s *WBService) getReportByPeriod(client *wbapi.WBClient, dateFrom, dateTo string, useNewAPI bool) ([]interface{}, error) {
	var url string
	if useNewAPI {
		url = fmt.Sprintf("%s?dateFrom=%s&dateTo=%s&rrdid=0&limit=100000",
			wbapi.URLDetail5(), dateFrom, dateTo)
	} else {
		url = fmt.Sprintf("%s?dateFrom=%s&dateTo=%s&rrdid=0&limit=100000",
			wbapi.URLDetails(), dateFrom, dateTo)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", client.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		var errorResp map[string]interface{}
		if err := json.Unmarshal(body, &errorResp); err == nil {
			return nil, fmt.Errorf("WB API error: %v", errorResp)
		}
		return nil, fmt.Errorf("WB API error: status %d", resp.StatusCode)
	}

	var data []interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return data, nil
}

func (s *WBService) saveStats(reportData []interface{}, userID int) (bool, string) {
	if len(reportData) == 0 {
		return false, "No data"
	}

	countSaved := 0
	countUnsaved := 0
	countTotal := len(reportData)

	fmt.Printf("📊 Получено %d записей от WB API\n", countTotal)

	for _, item := range reportData {
		order, ok := item.(map[string]interface{})
		if !ok {
			countUnsaved++
			continue
		}

		stat := s.mapToStat(order, userID)
		if stat == nil {
			countUnsaved++
			continue
		}

		// Генерируем хеш для проверки дубликатов
		hash := s.generateHash(stat)
		stat.HashInfo = hash

		// Проверяем существование
		exists, err := s.statRepo.ExistsByHash(hash)
		if err != nil {
			fmt.Printf("Error checking hash: %v\n", err)
			countUnsaved++
			continue
		}

		if exists {
			countUnsaved++
			continue
		}

		// Сохраняем
		if err := s.statRepo.Create(stat); err != nil {
			fmt.Printf("Error saving stat: %v\n", err)
			countUnsaved++
			continue
		}

		countSaved++
	}

	message := fmt.Sprintf("Total: %d, Saved: %d, Not saved (duplicates or errors): %d", countTotal, countSaved, countUnsaved)
	success := countSaved > 0

	fmt.Printf("✅ Результат: %s\n", message)
	return success, message
}

func (s *WBService) mapToStat(data map[string]interface{}, userID int) *entity.Stat {
	stat := &entity.Stat{
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	// Вспомогательная функция для преобразования значений
	setValue := func(value interface{}, setterFunc func(interface{})) {
		if value != nil {
			setterFunc(value)
		}
	}

	// Числовые поля (float64)
	setValue(data["delivery_rub"], func(v interface{}) {
		if num, ok := v.(float64); ok {
			stat.DeliveryRub = sql.NullFloat64{Float64: num, Valid: true}
		}
	})

	setValue(data["penalty"], func(v interface{}) {
		if num, ok := v.(float64); ok {
			stat.Penalty = sql.NullFloat64{Float64: num, Valid: true}
		}
	})

	setValue(data["additional_payment"], func(v interface{}) {
		if num, ok := v.(float64); ok {
			stat.AdditionalPayment = sql.NullFloat64{Float64: num, Valid: true}
		}
	})

	setValue(data["storage_fee"], func(v interface{}) {
		if num, ok := v.(float64); ok {
			stat.StorageFee = sql.NullFloat64{Float64: num, Valid: true}
		}
	})

	setValue(data["acquiring_fee"], func(v interface{}) {
		if num, ok := v.(float64); ok {
			stat.AcquiringFee = sql.NullFloat64{Float64: num, Valid: true}
		}
	})

	setValue(data["acquiring_percent"], func(v interface{}) {
		if num, ok := v.(float64); ok {
			stat.AcquiringPercent = sql.NullFloat64{Float64: num, Valid: true}
		}
	})

	setValue(data["ppvz_sales_commission"], func(v interface{}) {
		if num, ok := v.(float64); ok {
			stat.PpvzSalesCommission = sql.NullFloat64{Float64: num, Valid: true}
		}
	})

	setValue(data["deduction"], func(v interface{}) {
		if num, ok := v.(float64); ok {
			stat.Deduction = sql.NullFloat64{Float64: num, Valid: true}
		}
	})

	setValue(data["acceptance"], func(v interface{}) {
		if num, ok := v.(float64); ok {
			stat.Acceptance = sql.NullFloat64{Float64: num, Valid: true}
		}
	})

	setValue(data["dlv_prc"], func(v interface{}) {
		if num, ok := v.(float64); ok {
			stat.DlvPrc = sql.NullFloat64{Float64: num, Valid: true}
		}
	})

	setValue(data["retail_price"], func(v interface{}) {
		if num, ok := v.(float64); ok {
			stat.RetailPrice = sql.NullFloat64{Float64: num, Valid: true}
		}
	})

	setValue(data["retail_amount"], func(v interface{}) {
		if num, ok := v.(float64); ok {
			stat.RetailAmount = sql.NullFloat64{Float64: num, Valid: true}
		}
	})

	setValue(data["commission_percent"], func(v interface{}) {
		if num, ok := v.(float64); ok {
			stat.CommissionPercent = sql.NullFloat64{Float64: num, Valid: true}
		}
	})

	// Строковые поля, которые хранятся как VARCHAR
	setValue(data["ppvz_for_pay"], func(v interface{}) {
		if str, ok := v.(string); ok {
			stat.PpvzForPay = sql.NullString{String: str, Valid: true}
		} else if num, ok := v.(float64); ok {
			// Конвертируем число в строку
			stat.PpvzForPay = sql.NullString{String: fmt.Sprintf("%.2f", num), Valid: true}
		}
	})

	setValue(data["rebill_logistic_cost"], func(v interface{}) {
		if str, ok := v.(string); ok {
			stat.RebillLogisticCost = sql.NullString{String: str, Valid: true}
		} else if num, ok := v.(float64); ok {
			stat.RebillLogisticCost = sql.NullString{String: fmt.Sprintf("%.2f", num), Valid: true}
		}
	})

	setValue(data["ppvz_spp_prc"], func(v interface{}) {
		if str, ok := v.(string); ok {
			stat.PpvzSppPrc = sql.NullString{String: str, Valid: true}
		} else if num, ok := v.(float64); ok {
			stat.PpvzSppPrc = sql.NullString{String: fmt.Sprintf("%.2f", num), Valid: true}
		}
	})

	setValue(data["ppvz_kvw_prc_base"], func(v interface{}) {
		if str, ok := v.(string); ok {
			stat.PpvzKvwPrcBase = sql.NullString{String: str, Valid: true}
		} else if num, ok := v.(float64); ok {
			stat.PpvzKvwPrcBase = sql.NullString{String: fmt.Sprintf("%.2f", num), Valid: true}
		}
	})

	setValue(data["ppvz_kvw_prc"], func(v interface{}) {
		if str, ok := v.(string); ok {
			stat.PpvzKvwPrc = sql.NullString{String: str, Valid: true}
		} else if num, ok := v.(float64); ok {
			stat.PpvzKvwPrc = sql.NullString{String: fmt.Sprintf("%.2f", num), Valid: true}
		}
	})

	setValue(data["ppvz_vw_nds"], func(v interface{}) {
		if str, ok := v.(string); ok {
			stat.PpvzVwNds = sql.NullString{String: str, Valid: true}
		} else if num, ok := v.(float64); ok {
			stat.PpvzVwNds = sql.NullString{String: fmt.Sprintf("%.2f", num), Valid: true}
		}
	})

	setValue(data["ppvz_vw"], func(v interface{}) {
		if str, ok := v.(string); ok {
			stat.PpvzVw = sql.NullString{String: str, Valid: true}
		} else if num, ok := v.(float64); ok {
			stat.PpvzVw = sql.NullString{String: fmt.Sprintf("%.2f", num), Valid: true}
		}
	})

	// Целочисленные поля (bigint)
	setValue(data["nm_id"], func(v interface{}) {
		if num, ok := v.(float64); ok {
			stat.Nmid = sql.NullInt64{Int64: int64(num), Valid: true}
		}
	})

	setValue(data["quantity"], func(v interface{}) {
		if num, ok := v.(float64); ok {
			stat.Quantity = sql.NullInt64{Int64: int64(num), Valid: true}
		}
	})

	setValue(data["shk_id"], func(v interface{}) {
		if num, ok := v.(float64); ok {
			stat.ShkID = sql.NullInt64{Int64: int64(num), Valid: true}
		}
	})

	setValue(data["gi_id"], func(v interface{}) {
		if num, ok := v.(float64); ok {
			stat.GiID = sql.NullInt64{Int64: int64(num), Valid: true}
		}
	})

	setValue(data["realizationreport_id"], func(v interface{}) {
		if num, ok := v.(float64); ok {
			stat.RealizationreportID = sql.NullInt64{Int64: int64(num), Valid: true}
		}
	})

	setValue(data["ppvz_office_id"], func(v interface{}) {
		if num, ok := v.(float64); ok {
			stat.PpvzOfficeID = sql.NullInt64{Int64: int64(num), Valid: true}
		}
	})

	setValue(data["assembly_id"], func(v interface{}) {
		if num, ok := v.(float64); ok {
			stat.AssemblyID = sql.NullInt64{Int64: int64(num), Valid: true}
		}
	})

	setValue(data["delivery_amount"], func(v interface{}) {
		if num, ok := v.(float64); ok {
			stat.DeliveryAmount = sql.NullInt64{Int64: int64(num), Valid: true}
		}
	})

	setValue(data["return_amount"], func(v interface{}) {
		if num, ok := v.(float64); ok {
			stat.ReturnAmount = sql.NullInt64{Int64: int64(num), Valid: true}
		}
	})

	setValue(data["report_type"], func(v interface{}) {
		if num, ok := v.(float64); ok {
			stat.ReportType = sql.NullInt64{Int64: int64(num), Valid: true}
		}
	})

	setValue(data["rid"], func(v interface{}) {
		if num, ok := v.(float64); ok {
			stat.Rid = sql.NullInt64{Int64: int64(num), Valid: true}
		} else if str, ok := v.(string); ok {
			// Пробуем преобразовать строку в число
			if num, err := strconv.ParseInt(str, 10, 64); err == nil {
				stat.Rid = sql.NullInt64{Int64: num, Valid: true}
			}
		}
	})

	// supplier_oper_name - особый случай, в БД это integer
	setValue(data["supplier_oper_name"], func(v interface{}) {
		// Преобразуем как в Yii2
		supplierType := s.convertSupplierOperName(v)
		stat.SupplierOperName = sql.NullInt64{Int64: supplierType, Valid: true}

		// Отладка
		if supplierType > 0 {
			fmt.Printf("DEBUG: supplier_oper_name '%v' -> %d\n", v, supplierType)
		}
	})

	// Строковые поля
	setValue(data["subject_name"], func(v interface{}) {
		if str, ok := v.(string); ok {
			stat.SubjectName = sql.NullString{String: str, Valid: true}
		}
	})

	setValue(data["brand_name"], func(v interface{}) {
		if str, ok := v.(string); ok {
			stat.BrandName = sql.NullString{String: str, Valid: true}
		}
	})

	setValue(data["office_name"], func(v interface{}) {
		if str, ok := v.(string); ok {
			stat.OfficeName = sql.NullString{String: str, Valid: true}
		}
	})

	setValue(data["barcode"], func(v interface{}) {
		if str, ok := v.(string); ok {
			stat.Barcode = sql.NullString{String: str, Valid: true}
		}
	})

	setValue(data["bonus_type_name"], func(v interface{}) {
		if str, ok := v.(string); ok {
			stat.BonusTypeName = sql.NullString{String: str, Valid: true}
		}
	})

	setValue(data["last_error"], func(v interface{}) {
		if str, ok := v.(string); ok {
			stat.LastError = sql.NullString{String: str, Valid: true}
		}
	})

	setValue(data["sa_name"], func(v interface{}) {
		if str, ok := v.(string); ok {
			stat.SaName = sql.NullString{String: str, Valid: true}
		}
	})

	setValue(data["gi_box_type_name"], func(v interface{}) {
		if str, ok := v.(string); ok {
			stat.GiBoxTypeName = sql.NullString{String: str, Valid: true}
		}
	})

	setValue(data["ts_name"], func(v interface{}) {
		if str, ok := v.(string); ok {
			stat.TsName = sql.NullString{String: str, Valid: true}
		}
	})

	setValue(data["sticker_id"], func(v interface{}) {
		if str, ok := v.(string); ok {
			stat.StickerID = sql.NullString{String: str, Valid: true}
		}
	})

	setValue(data["srid"], func(v interface{}) {
		if str, ok := v.(string); ok {
			stat.Srid = sql.NullString{String: str, Valid: true}
		}
	})

	// Даты
	setValue(data["order_dt"], func(v interface{}) {
		if str, ok := v.(string); ok {
			// Пробуем разные форматы дат
			formats := []string{
				"2006-01-02T15:04:05",
				"2006-01-02T15:04:05Z",
				"2006-01-02 15:04:05",
				time.RFC3339,
			}

			for _, format := range formats {
				if t, err := time.Parse(format, str); err == nil {
					stat.OrderDt = sql.NullTime{Time: t, Valid: true}
					return
				}
			}
		}
	})

	setValue(data["sale_dt"], func(v interface{}) {
		if str, ok := v.(string); ok {
			formats := []string{
				"2006-01-02T15:04:05",
				"2006-01-02T15:04:05Z",
				"2006-01-02 15:04:05",
				time.RFC3339,
			}

			for _, format := range formats {
				if t, err := time.Parse(format, str); err == nil {
					stat.SaleDt = sql.NullTime{Time: t, Valid: true}
					return
				}
			}
		}
	})

	setValue(data["rr_dt"], func(v interface{}) {
		if str, ok := v.(string); ok {
			// Формат даты без времени
			if t, err := time.Parse("2006-01-02", str); err == nil {
				stat.RrDt = sql.NullTime{Time: t, Valid: true}
			}
		}
	})

	return stat
}

func (s *WBService) generateHash(stat *entity.Stat) string {
	// Воспроизводим логику PHP: создаем маску из всех полей
	// Адаптируйте эту функцию под вашу конкретную PHP логику генерации хеша

	hashParts := []string{}

	// Добавляем основные поля (аналогично PHP)
	if stat.UserID > 0 {
		hashParts = append(hashParts, strconv.Itoa(stat.UserID))
	}

	if stat.RealizationreportID.Valid {
		hashParts = append(hashParts, strconv.FormatInt(stat.RealizationreportID.Int64, 10))
	} else {
		hashParts = append(hashParts, "0")
	}

	if stat.Rid.Valid {
		hashParts = append(hashParts, strconv.FormatInt(stat.Rid.Int64, 10))
	} else {
		hashParts = append(hashParts, "0")
	}

	if stat.Srid.Valid {
		hashParts = append(hashParts, stat.Srid.String)
	} else {
		hashParts = append(hashParts, "")
	}

	if stat.SupplierOperName.Valid {
		hashParts = append(hashParts, strconv.FormatInt(stat.SupplierOperName.Int64, 10))
	} else {
		hashParts = append(hashParts, "0")
	}

	if stat.ReportType.Valid {
		hashParts = append(hashParts, strconv.FormatInt(stat.ReportType.Int64, 10))
	} else {
		hashParts = append(hashParts, "0")
	}

	if stat.ShkID.Valid {
		hashParts = append(hashParts, strconv.FormatInt(stat.ShkID.Int64, 10))
	} else {
		hashParts = append(hashParts, "0")
	}

	if stat.Nmid.Valid {
		hashParts = append(hashParts, strconv.FormatInt(stat.Nmid.Int64, 10))
	} else {
		hashParts = append(hashParts, "0")
	}

	if stat.GiID.Valid {
		hashParts = append(hashParts, strconv.FormatInt(stat.GiID.Int64, 10))
	} else {
		hashParts = append(hashParts, "0")
	}

	if stat.Quantity.Valid {
		hashParts = append(hashParts, strconv.FormatInt(stat.Quantity.Int64, 10))
	} else {
		hashParts = append(hashParts, "0")
	}

	// rebill_logistic_cost
	if stat.RebillLogisticCost.Valid {
		// Убираем запятые как в PHP
		value := strings.ReplaceAll(stat.RebillLogisticCost.String, ",", "")
		hashParts = append(hashParts, value)
	} else {
		hashParts = append(hashParts, "0")
	}

	// return_amount
	if stat.ReturnAmount.Valid {
		hashParts = append(hashParts, strconv.FormatInt(stat.ReturnAmount.Int64, 10))
	} else {
		hashParts = append(hashParts, "0")
	}

	// retail_price
	if stat.RetailPrice.Valid {
		hashParts = append(hashParts, fmt.Sprintf("%.2f", stat.RetailPrice.Float64))
	} else {
		hashParts = append(hashParts, "0")
	}

	// retail_amount
	if stat.RetailAmount.Valid {
		hashParts = append(hashParts, fmt.Sprintf("%.2f", stat.RetailAmount.Float64))
	} else {
		hashParts = append(hashParts, "0")
	}

	// subject_name
	if stat.SubjectName.Valid {
		hashParts = append(hashParts, stat.SubjectName.String)
	} else {
		hashParts = append(hashParts, "")
	}

	// sa_name
	if stat.SaName.Valid {
		hashParts = append(hashParts, stat.SaName.String)
	} else {
		hashParts = append(hashParts, "")
	}

	// ppvz_vw_nds
	if stat.PpvzVwNds.Valid {
		value := strings.ReplaceAll(stat.PpvzVwNds.String, ",", "")
		hashParts = append(hashParts, value)
	} else {
		hashParts = append(hashParts, "0")
	}

	// ppvz_vw
	if stat.PpvzVw.Valid {
		value := strings.ReplaceAll(stat.PpvzVw.String, ",", "")
		hashParts = append(hashParts, value)
	} else {
		hashParts = append(hashParts, "0")
	}

	// ppvz_spp_prc
	if stat.PpvzSppPrc.Valid {
		value := strings.ReplaceAll(stat.PpvzSppPrc.String, ",", "")
		hashParts = append(hashParts, value)
	} else {
		hashParts = append(hashParts, "0")
	}

	// ppvz_kvw_prc_base
	if stat.PpvzKvwPrcBase.Valid {
		value := strings.ReplaceAll(stat.PpvzKvwPrcBase.String, ",", "")
		hashParts = append(hashParts, value)
	} else {
		hashParts = append(hashParts, "0")
	}

	// ppvz_kvw_prc
	if stat.PpvzKvwPrc.Valid {
		value := strings.ReplaceAll(stat.PpvzKvwPrc.String, ",", "")
		hashParts = append(hashParts, value)
	} else {
		hashParts = append(hashParts, "0")
	}

	// ppvz_sales_commission
	if stat.PpvzSalesCommission.Valid {
		hashParts = append(hashParts, fmt.Sprintf("%.2f", stat.PpvzSalesCommission.Float64))
	} else {
		hashParts = append(hashParts, "0")
	}

	// acquiring_fee
	if stat.AcquiringFee.Valid {
		hashParts = append(hashParts, fmt.Sprintf("%.2f", stat.AcquiringFee.Float64))
	} else {
		hashParts = append(hashParts, "0")
	}

	// assembly_id
	if stat.AssemblyID.Valid {
		hashParts = append(hashParts, strconv.FormatInt(stat.AssemblyID.Int64, 10))
	} else {
		hashParts = append(hashParts, "0")
	}

	// acquiring_percent
	if stat.AcquiringPercent.Valid {
		hashParts = append(hashParts, fmt.Sprintf("%.2f", stat.AcquiringPercent.Float64))
	} else {
		hashParts = append(hashParts, "0")
	}

	// gi_box_type_name
	if stat.GiBoxTypeName.Valid {
		hashParts = append(hashParts, stat.GiBoxTypeName.String)
	} else {
		hashParts = append(hashParts, "")
	}

	// acceptance
	if stat.Acceptance.Valid {
		hashParts = append(hashParts, fmt.Sprintf("%.2f", stat.Acceptance.Float64))
	} else {
		hashParts = append(hashParts, "0")
	}

	// commission_percent
	if stat.CommissionPercent.Valid {
		value := strings.ReplaceAll(fmt.Sprintf("%.2f", stat.CommissionPercent.Float64), ",", "")
		hashParts = append(hashParts, value)
	} else {
		hashParts = append(hashParts, "0")
	}

	// delivery_amount
	if stat.DeliveryAmount.Valid {
		hashParts = append(hashParts, strconv.FormatInt(stat.DeliveryAmount.Int64, 10))
	} else {
		hashParts = append(hashParts, "0")
	}

	// delivery_rub
	if stat.DeliveryRub.Valid {
		hashParts = append(hashParts, fmt.Sprintf("%.2f", stat.DeliveryRub.Float64))
	} else {
		hashParts = append(hashParts, "0")
	}

	// bonus_type_name
	if stat.BonusTypeName.Valid {
		hashParts = append(hashParts, stat.BonusTypeName.String)
	} else {
		hashParts = append(hashParts, "")
	}

	// ppvz_for_pay
	if stat.PpvzForPay.Valid {
		value := strings.ReplaceAll(stat.PpvzForPay.String, ",", "")
		hashParts = append(hashParts, value)
	} else {
		hashParts = append(hashParts, "0")
	}

	// ppvz_office_id
	if stat.PpvzOfficeID.Valid {
		hashParts = append(hashParts, strconv.FormatInt(stat.PpvzOfficeID.Int64, 10))
	} else {
		hashParts = append(hashParts, "0")
	}

	// sticker_id
	if stat.StickerID.Valid {
		hashParts = append(hashParts, stat.StickerID.String)
	} else {
		hashParts = append(hashParts, "")
	}

	// office_name
	if stat.OfficeName.Valid {
		hashParts = append(hashParts, stat.OfficeName.String)
	} else {
		hashParts = append(hashParts, "")
	}

	// penalty
	if stat.Penalty.Valid {
		hashParts = append(hashParts, fmt.Sprintf("%.2f", stat.Penalty.Float64))
	} else {
		hashParts = append(hashParts, "0")
	}

	// ts_name
	if stat.TsName.Valid {
		hashParts = append(hashParts, stat.TsName.String)
	} else {
		hashParts = append(hashParts, "")
	}

	// order_dt
	if stat.OrderDt.Valid {
		hashParts = append(hashParts, stat.OrderDt.Time.Format("2006-01-02 15:04:05"))
	} else {
		hashParts = append(hashParts, "")
	}

	// sale_dt
	if stat.SaleDt.Valid {
		hashParts = append(hashParts, stat.SaleDt.Time.Format("2006-01-02 15:04:05"))
	} else {
		hashParts = append(hashParts, "")
	}

	hashMask := strings.Join(hashParts, "")
	hash := sha256.Sum256([]byte(hashMask))
	return fmt.Sprintf("%x", hash)
}

// GetSupplierType - конвертация типа поставщика (аналог PHP getSuplierType)
func (s *WBService) GetSupplierType(supplierType interface{}) string {
	if supplierType == nil {
		return ""
	}

	switch v := supplierType.(type) {
	case string:
		// Ваша логика преобразования (как в PHP Stat::getSuplierType)
		switch v {
		case "1":
			return "Продажа"
		case "2":
			return "Возврат"
		default:
			return v
		}
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", v)
	}
}

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

func (s *WBService) processArticleRequest(articleReq *entity.WBArticlesGet, user *repository.User) ProcessResult {
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

func (s *WBService) getWBArticles(user *repository.User) ([]entity.WBArticle, error) {
	if !user.WbKey.Valid || user.WbKey.String == "" {
		return nil, fmt.Errorf("токен WB не указан для пользователя %d", user.ID)
	}

	token := user.WbKey.String
	client := wbapi.NewWBClient(token)

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

func (s *WBService) fetchArticlesFromWB(client *wbapi.WBClient) ([]entity.WBArticle, error) {
	var allCards []entity.WBArticle
	var cursorUpdatedAt string
	var cursorNmID int

	limit := 100 // Максимальный лимит за один запрос
	totalProcessed := 0

	for {
		// Формируем тело запроса
		request := entity.WBArticleRequest{
			Settings: entity.WBArticleRequestSettings{
				Cursor: entity.WBArticleRequestCursor{
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

		url := wbapi.URLCardsList()
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
		var response entity.WBArticleResponse
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

func (s *WBService) saveArticles(cards []entity.WBArticle, userID int) (bool, string) {
	if len(cards) == 0 {
		return false, "No articles data received"
	}

	countSaved := 0
	countUnsaved := 0
	countTotal := len(cards)

	fmt.Printf("📦 Получено %d карточек товаров от WB API\n", countTotal)

	for _, card := range cards {
		// Создаем запись для базы данных (используем WBArticleDB)
		article := &entity.WBArticleDB{
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

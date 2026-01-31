package wb

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
	"wbrost-go/internal/entity"
)

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
		//if supplierType > 0 {
		//	fmt.Printf("DEBUG: supplier_oper_name '%v' -> %d\n", v, supplierType)
		//}
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

package stat

import (
	"database/sql"
	"fmt"
	"time"
)

// Helper функции для работы с NULL значениями
func getNullString(ns sql.NullString) interface{} {
	if ns.Valid {
		return ns.String
	}
	return nil
}

func getNullInt64(ni sql.NullInt64) interface{} {
	if ni.Valid {
		return ni.Int64
	}
	return nil
}

func getNullFloat64(nf sql.NullFloat64) interface{} {
	if nf.Valid {
		return nf.Float64
	}
	return nil
}

func getNullTime(nt sql.NullTime) interface{} {
	if nt.Valid {
		return nt.Time
	}
	return nil
}
func getStringValue(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
func getFloatValue(nf sql.NullFloat64) float64 {
	if nf.Valid {
		return nf.Float64
	}
	return 0.0
}

func getIntValue(ni sql.NullInt64) int64 {
	if ni.Valid {
		return ni.Int64
	}
	return 0
}
func getPhotoURL(nmID int64) string {
	if nmID <= 0 {
		return ""
	}

	// Формируем URL по аналогии с вашим примером
	basketNum := nmID % 100
	vol := nmID / 1000000
	part := nmID / 1000

	return fmt.Sprintf("https://basket-%d.wbbasket.ru/vol%d/part%d/%d/images/big/1.webp",
		basketNum, vol, part, nmID)
}

// formatMonthLabel создает метку в формате "Янв 2025"
func formatMonthLabel(monthName string, year int) string {
	// Преобразуем английское название месяца в русское
	monthMap := map[string]string{
		"Jan": "Янв", "Feb": "Фев", "Mar": "Мар", "Apr": "Апр",
		"May": "Май", "Jun": "Июн", "Jul": "Июл", "Aug": "Авг",
		"Sep": "Сен", "Oct": "Окт", "Nov": "Ноя", "Dec": "Дек",
	}

	// Берем первые 3 символа
	shortName := ""
	if len(monthName) >= 3 {
		shortName = monthName[:3]
	}

	// Преобразуем в русское название
	ruMonth := shortName
	if ruName, ok := monthMap[shortName]; ok {
		ruMonth = ruName
	}

	// Возвращаем в формате "Янв 2025"
	return fmt.Sprintf("%s %d", ruMonth, year)
}

// formatMonthNameRU можно удалить или оставить, если она используется где-то еще
func formatMonthNameRU(monthName string) string {
	monthMap := map[string]string{
		"Jan": "Янв", "Feb": "Фев", "Mar": "Мар", "Apr": "Апр",
		"May": "Май", "Jun": "Июн", "Jul": "Июл", "Aug": "Авг",
		"Sep": "Сен", "Oct": "Окт", "Nov": "Ноя", "Dec": "Дек",
	}

	if len(monthName) >= 3 {
		shortName := monthName[:3]
		if ruName, ok := monthMap[shortName]; ok {
			return ruName
		}
	}

	return monthName
}

// Вспомогательная функция для генерации недостающих месяцев
func (r *StatRepository) generateMissingMonths(startDate, endDate time.Time, existingMonths map[string]bool,
	currentLabels []string, currentData []float64) ([]string, []float64) {

	allLabels := make([]string, 0)
	allData := make([]float64, 0)

	// Генерируем все месяцы в диапазоне
	for current := startDate; current.Before(endDate.AddDate(0, 1, 0)); current = current.AddDate(0, 1, 0) {
		monthName := formatMonthNameRU(current.Format("Jan"))
		label := fmt.Sprintf("%s %d", monthName, current.Year())

		// Проверяем, есть ли данные для этого месяца
		if existingMonths[label] {
			// Находим индекс в текущих данных
			for i, lbl := range currentLabels {
				if lbl == label {
					allLabels = append(allLabels, lbl)
					allData = append(allData, currentData[i])
					break
				}
			}
		} else {
			// Добавляем месяц с нулевыми данными
			allLabels = append(allLabels, label)
			allData = append(allData, 0)
		}

		// Ограничиваем 12 месяцами
		if len(allLabels) >= 12 {
			break
		}
	}

	return allLabels, allData
}

package service

import (
	"context"
	"fmt"
	"time"
	"wbrost-go/internal/entity"
	"wbrost-go/internal/repository"
)

type SiteService interface {
	GetDashboardData(ctx context.Context, userID int) (*DashboardData, error)
	GetChartData(ctx context.Context, userID int, chartType string) (*entity.ChartData, error)
}

type siteService struct {
	statRepo repository.StatRepository
	userRepo repository.UserRepository
}

func NewSiteService(statRepo repository.StatRepository, userRepo repository.UserRepository) SiteService {
	return &siteService{
		statRepo: statRepo,
		userRepo: userRepo,
	}
}

type DashboardData struct {
	TotalSell       string            `json:"total_sell"`
	TotalCountSales int               `json:"total_count_sales"`
	TotalRefund     int               `json:"total_refund"`
	TotalProfit     string            `json:"total_profit"`
	DataLine        *entity.ChartData `json:"data_line"`
	DataBar         *entity.ChartData `json:"data_bar"`
	IsGuest         bool              `json:"is_guest"`
}

func (s *siteService) GetDashboardData(ctx context.Context, userID int) (*DashboardData, error) {
	// Если userID = 0, значит гость
	if userID == 0 {
		return &DashboardData{
			IsGuest: true,
		}, nil
	}

	// Даты как в Yii2
	dateFrom := time.Now().Format("2006-01-01") // Первое число текущего месяца
	dateTo := time.Now().Format("2006-01-02")   // Сегодня

	// Получаем данные за последние 3 месяца для графиков
	sumDataProviders, err := s.statRepo.GetMonthlyStats(ctx, userID, 3)
	if err != nil {
		return nil, err
	}

	// Получаем итоговые значения за текущий месяц
	totalSell, _ := s.statRepo.GetTotal(ctx, userID, dateFrom, dateTo, "ppvz_for_pay")
	totalCountSales, _ := s.statRepo.GetTotal(ctx, userID, dateFrom, dateTo, "count_sales")
	totalRefund, _ := s.statRepo.GetTotal(ctx, userID, dateFrom, dateTo, "count_refund")
	totalProfit, _ := s.statRepo.GetTotalClear(ctx, userID, dateFrom, dateTo)

	// Форматируем как в Yii2
	formatNumber := func(num float64) string {
		intPart := int(num)
		decPart := int((num - float64(intPart)) * 100)
		return fmt.Sprintf("%d,%02d", intPart, decPart)
	}

	dataLine := s.getChartData(sumDataProviders, entity.DataTypeLine)
	dataBar := s.getChartData(sumDataProviders, entity.DataTypeBar)

	return &DashboardData{
		TotalSell:       formatNumber(totalSell),
		TotalCountSales: int(totalCountSales),
		TotalRefund:     int(totalRefund),
		TotalProfit:     formatNumber(totalProfit),
		DataLine:        dataLine,
		DataBar:         dataBar,
		IsGuest:         false,
	}, nil
}

func (s *siteService) getChartData(sumDataProviders []entity.StatSummary, chartType string) *entity.ChartData {
	var labels []string
	var salesData []interface{}
	var totalSalesData []interface{}
	var refundsData []interface{}
	var backgroundColors []string
	var borderColors []string

	for idx, data := range sumDataProviders {
		labels = append(labels, data.Month)
		salesData = append(salesData, data.TotalSell)
		totalSalesData = append(totalSalesData, data.TotalCountSales)
		refundsData = append(refundsData, data.TotalRefund)

		backgroundColors = append(backgroundColors, "rgba(201, 203, 207, 0.7)")
		borderColors = append(borderColors, "rgb(201, 203, 207)")

		if idx == len(sumDataProviders)-1 {
			backgroundColors[idx] = "rgba(75, 192, 192, 0.2)"
			borderColors[idx] = "rgb(75, 192, 192)"
		}
	}

	if chartType == entity.DataTypeLine {
		return &entity.ChartData{
			Labels: labels,
			Datasets: []entity.ChartDataset{
				{
					Label:           "Продажи (кол-во)",
					BackgroundColor: "rgba(179,181,198,0.2)",
					BorderColor:     "rgba(44,163,203,255)",
					Data:            totalSalesData,
				},
				{
					Label:           "Возвраты (кол-во)",
					BackgroundColor: "rgba(255,99,132,0.2)",
					BorderColor:     "rgba(255,99,132,1)",
					Data:            refundsData,
				},
			},
		}
	}

	if chartType == entity.DataTypeBar {
		return &entity.ChartData{
			Labels: labels,
			Datasets: []entity.ChartDataset{
				{
					Label:           "Выручка в рублях",
					BackgroundColor: backgroundColors[0],
					BorderColor:     borderColors[0],
					Data:            salesData,
				},
			},
		}
	}

	return nil
}

func (s *siteService) GetChartData(ctx context.Context, userID int, chartType string) (*entity.ChartData, error) {
	sumDataProviders, err := s.statRepo.GetMonthlyStats(ctx, userID, 12) // За последний год
	if err != nil {
		return nil, err
	}

	return s.getChartData(sumDataProviders, chartType), nil
}

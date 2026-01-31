package wb

import (
	"wbrost-go/internal/repository/article"
	"wbrost-go/internal/repository/stat"
	"wbrost-go/internal/repository/user"
)

type WBService struct {
	userRepo        *user.UserRepository
	statsGetRepo    *stat.WBStatsGetRepository
	statRepo        *stat.StatRepository
	articlesGetRepo *article.WBArticlesGetRepository
	articleRepo     *article.WBArticlesRepository
	rateLimiter     *WBRateLimiter
}

func NewWBService(
	userRepo *user.UserRepository,
	statsGetRepo *stat.WBStatsGetRepository,
	statRepo *stat.StatRepository,
	articlesGetRepo *article.WBArticlesGetRepository,
	articleRepo *article.WBArticlesRepository,
) *WBService {
	// Используем новый rate limiter с поддержкой WB API
	rateLimiter := NewWBRateLimiter()

	return &WBService{
		userRepo:        userRepo,
		statsGetRepo:    statsGetRepo,
		statRepo:        statRepo,
		articlesGetRepo: articlesGetRepo,
		articleRepo:     articleRepo,
		rateLimiter:     rateLimiter,
	}
}

// GetLimiterStats возвращает статистику rate limiter
func (s *WBService) GetLimiterStats() map[string]interface{} {
	return s.rateLimiter.GetStats()
}

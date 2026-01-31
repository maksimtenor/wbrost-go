package wb

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type WBRateLimiter struct {
	mu sync.Mutex

	// Лимиты для WB API
	maxRequestsPerMinute int           // 50 запросов в минуту - безопасное значение
	minRequestInterval   time.Duration // 2 секунды между запросами

	// Состояние
	lastRequest    time.Time
	requestCount   int
	rateLimitStart time.Time

	// Для обработки 429 ошибок
	rateLimited    bool
	rateLimitUntil time.Time
	retryAfter     time.Duration

	// Статистика
	stats struct {
		TotalRequests   int64
		RateLimitHits   int64
		LastRequestTime time.Time
		Last429Time     time.Time
	}
}

func NewWBRateLimiter() *WBRateLimiter {
	return &WBRateLimiter{
		maxRequestsPerMinute: 50,              // Увеличено с 1 до 50!
		minRequestInterval:   2 * time.Second, // 2 секунды между запросами
		lastRequest:          time.Now().Add(-time.Minute),
	}
}

func (r *WBRateLimiter) Wait() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	// 1. Проверяем активный rate limit (после 429 ошибки)
	if r.rateLimited && now.Before(r.rateLimitUntil) {
		waitTime := r.rateLimitUntil.Sub(now)
		fmt.Printf("⏳ Активный rate limit. Ждем %v...\n", waitTime)
		time.Sleep(waitTime)
		r.rateLimited = false
		now = time.Now()
	}

	// 2. Сбрасываем счетчик если прошла минута
	if now.Sub(r.rateLimitStart) >= time.Minute {
		r.requestCount = 0
		r.rateLimitStart = now
	}

	// 3. Проверяем лимит запросов в минуту
	if r.requestCount >= r.maxRequestsPerMinute {
		waitTime := r.rateLimitStart.Add(time.Minute).Sub(now)
		if waitTime > 0 {
			fmt.Printf("⏳ Достигнут минутный лимит. Ждем %v...\n", waitTime)
			time.Sleep(waitTime)
			now = time.Now()
			r.requestCount = 0
			r.rateLimitStart = now
		}
	}

	// 4. Минимальный интервал между запросами
	if !r.lastRequest.IsZero() {
		elapsed := now.Sub(r.lastRequest)
		if elapsed < r.minRequestInterval {
			waitTime := r.minRequestInterval - elapsed
			if waitTime > 0 {
				fmt.Printf("⏳ Соблюдаем интервал. Ждем %v...\n", waitTime)
				time.Sleep(waitTime)
				now = time.Now()
			}
		}
	}

	r.lastRequest = now
	r.requestCount++
	r.stats.TotalRequests++
	r.stats.LastRequestTime = now

	return nil
}

func (r *WBRateLimiter) ProcessHeaders(headers http.Header, statusCode int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if statusCode == 429 {
		r.stats.RateLimitHits++
		r.stats.Last429Time = time.Now()
		r.rateLimited = true

		// Пробуем получить Retry-After из заголовков
		if retryAfter := headers.Get("Retry-After"); retryAfter != "" {
			if seconds, err := strconv.Atoi(retryAfter); err == nil && seconds > 0 {
				r.retryAfter = time.Duration(seconds) * time.Second
			} else {
				// Если не получилось, используем 30 секунд
				r.retryAfter = 30 * time.Second
			}
		} else {
			r.retryAfter = 30 * time.Second
		}

		r.rateLimitUntil = time.Now().Add(r.retryAfter)
		fmt.Printf("⚠️ Получен 429 Too Many Requests. Ждем %v\n", r.retryAfter)
	} else if statusCode == 200 {
		// При успешном запросе сбрасываем флаг rate limit
		r.rateLimited = false
	}
}

func (r *WBRateLimiter) GetStats() map[string]interface{} {
	r.mu.Lock()
	defer r.mu.Unlock()

	stats := map[string]interface{}{
		"total_requests":    r.stats.TotalRequests,
		"rate_limit_hits":   r.stats.RateLimitHits,
		"last_request":      r.stats.LastRequestTime.Format("2006-01-02 15:04:05"),
		"requests_this_min": r.requestCount,
		"max_per_minute":    r.maxRequestsPerMinute,
		"rate_limited":      r.rateLimited,
	}

	if r.stats.Last429Time.IsZero() {
		stats["last_429"] = "never"
	} else {
		stats["last_429"] = r.stats.Last429Time.Format("2006-01-02 15:04:05")
	}

	if r.rateLimited {
		stats["rate_limit_until"] = r.rateLimitUntil.Format("15:04:05")
		stats["retry_after"] = r.retryAfter.String()
	}

	return stats
}

package middleware

import (
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/cesarbmathec/medical-exams-backend/utils"
	"github.com/gin-gonic/gin"
)

type limiter struct {
	mu       sync.Mutex
	limiters map[string]*clientLimiter
	limit    float64
	burst    int
	window   time.Duration
}

type clientLimiter struct {
	tokens     float64
	lastRefill time.Time
}

func NewRateLimiter() *limiter {
	limit := getEnvInt("RATE_LIMIT_PER_MINUTE", 120)
	burst := getEnvInt("RATE_LIMIT_BURST", 30)
	window := time.Minute
	if limit <= 0 {
		limit = 120
	}
	if burst <= 0 {
		burst = 30
	}

	return &limiter{
		limiters: make(map[string]*clientLimiter),
		limit:    float64(limit) / window.Seconds(),
		burst:    burst,
		window:   window,
	}
}

func (l *limiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if stringsEqualFold(os.Getenv("RATE_LIMIT_ENABLED"), "false") {
			c.Next()
			return
		}

		key := c.ClientIP()
		if key == "" {
			key = "unknown"
		}

		if !l.allow(key) {
			utils.Error(c, http.StatusTooManyRequests, "Demasiadas solicitudes", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

func (l *limiter) allow(key string) bool {
	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	cl, ok := l.limiters[key]
	if !ok {
		l.limiters[key] = &clientLimiter{tokens: float64(l.burst) - 1, lastRefill: now}
		return true
	}

	elapsed := now.Sub(cl.lastRefill).Seconds()
	cl.tokens += elapsed * l.limit
	if cl.tokens > float64(l.burst) {
		cl.tokens = float64(l.burst)
	}
	cl.lastRefill = now

	if cl.tokens < 1 {
		return false
	}

	cl.tokens -= 1
	return true
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func stringsEqualFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca := a[i]
		cb := b[i]
		if ca == cb {
			continue
		}
		if ca >= 'A' && ca <= 'Z' {
			ca = ca + ('a' - 'A')
		}
		if cb >= 'A' && cb <= 'Z' {
			cb = cb + ('a' - 'A')
		}
		if ca != cb {
			return false
		}
	}
	return true
}

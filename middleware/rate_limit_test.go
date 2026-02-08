package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRateLimiterBlocksAfterBurst(t *testing.T) {
	os.Setenv("RATE_LIMIT_ENABLED", "true")
	os.Setenv("RATE_LIMIT_PER_MINUTE", "1")
	os.Setenv("RATE_LIMIT_BURST", "2")
	defer os.Unsetenv("RATE_LIMIT_ENABLED")
	defer os.Unsetenv("RATE_LIMIT_PER_MINUTE")
	defer os.Unsetenv("RATE_LIMIT_BURST")

	gin.SetMode(gin.TestMode)
	r := gin.New()
	limiter := NewRateLimiter()
	r.Use(limiter.Middleware())
	r.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })

	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", w.Code)
	}
}

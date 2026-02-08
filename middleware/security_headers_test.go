package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSecurityHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(SecurityHeaders())
	r.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	if w.Header().Get("X-Content-Type-Options") == "" {
		t.Fatal("missing X-Content-Type-Options")
	}
	if w.Header().Get("X-Frame-Options") == "" {
		t.Fatal("missing X-Frame-Options")
	}
	if w.Header().Get("Referrer-Policy") == "" {
		t.Fatal("missing Referrer-Policy")
	}
}

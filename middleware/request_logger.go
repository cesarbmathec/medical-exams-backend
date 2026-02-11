package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		path := param.Path
		if path == "" && param.Request != nil {
			path = param.Request.URL.Path
		}

		// Extraemos el Origin de la solicitud
		origin := param.Request.Header.Get("Origin")
		if origin == "" {
			origin = "No-Origin (Direct/Same-Domain)"
		}

		return fmt.Sprintf(
			"%s - %s \"%s %s\" %d %s | Origin: %s\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC3339),
			param.Method,
			path,
			param.StatusCode,
			param.Latency,
			origin,
		)
	})
}

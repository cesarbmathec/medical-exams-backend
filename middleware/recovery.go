package middleware

import (
	"net/http"
	"os"
	"runtime/debug"

	"github.com/cesarbmathec/medical-exams-backend/utils"
	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log del error y el stack trace solo en entornos no-release
				if os.Getenv("GIN_MODE") != "release" {
					debug.PrintStack()
				}

				// Enviamos la respuesta estandarizada
				utils.Error(c, http.StatusInternalServerError, "Ha ocurrido un error inesperado en el servidor", nil)

				// Detenemos la ejecución de los siguientes handlers
				c.Abort()
			}
		}()
		c.Next()
	}
}

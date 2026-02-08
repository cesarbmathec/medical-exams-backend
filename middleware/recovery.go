package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/cesarbmathec/medical-exams-backend/utils"
	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log del error y el stack trace en la consola para depuración
				fmt.Printf("--- PANIC RECOVERED ---\nError: %v\nStack: %s\n-----------------------\n", err, debug.Stack())

				// Enviamos la respuesta estandarizada usando tu utilidad
				utils.Error(
					c,
					http.StatusInternalServerError,
					"Ha ocurrido un error inesperado en el servidor",
					fmt.Sprintf("%v", err),
				)

				// Detenemos la ejecución de los siguientes handlers
				c.Abort()
			}
		}()
		c.Next()
	}
}

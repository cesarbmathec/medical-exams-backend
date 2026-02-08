package utils

import (
	"github.com/gin-gonic/gin"
)

// Response es la estructura base para todas las respuestas
type Response struct {
	Status  string      `json:"status"`           // "success" o "error"
	Code    int         `json:"code"`             // Código HTTP (200, 400, 500, etc.)
	Message string      `json:"message"`          // Mensaje legible para el usuario/dev
	Data    interface{} `json:"data,omitempty"`   // Datos (solo en respuestas exitosas)
	Errors  interface{} `json:"errors,omitempty"` // Detalles técnicos (opcional)
}

// Success envía una respuesta 200/201 estandarizada
func Success(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Status:  "success",
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// Error envía una respuesta de error estandarizada
func Error(c *gin.Context, code int, message string, errs interface{}) {
	c.JSON(code, Response{
		Status:  "error",
		Code:    code,
		Message: message,
		Errors:  errs,
	})
}

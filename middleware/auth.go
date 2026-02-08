package middleware

import (
	"net/http"
	"strings"

	"github.com/cesarbmathec/medical-exams-backend/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Se requiere token de autorización"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
			c.Abort()
			return
		}
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Guardamos el ID del usuario en el contexto para saber quién opera
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roleID", claims.RoleID)
		c.Next()
	}
}

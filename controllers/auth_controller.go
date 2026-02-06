package controllers

import (
	"net/http"

	"github.com/cesarbmathec/medical-exams-backend/config"
	"github.com/cesarbmathec/medical-exams-backend/models"
	"github.com/cesarbmathec/medical-exams-backend/utils"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	RoleID   uint   `json:"role_id" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	var user models.User
	db := config.GetDB()

	// Buscar usuario e incluir el rol
	if err := db.Preload("Role").Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario o contraseña incorrectos"})
		return
	}

	// Verificar password usando el método que definiste en user.go
	if !user.CheckPassword(input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario o contraseña incorrectos"})
		return
	}

	// Generar Token
	token, err := utils.GenerateToken(user.ID, user.Username, user.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo generar el token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Inicio de sesión exitoso",
		"token":   token,
		"user":    user.ToResponse(),
	})
}

func Register(c *gin.Context) {
	var input RegisterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.GetDB()

	// Creamos la instancia del modelo
	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password, // El hook BeforeCreate en user.go hará el Hash
		FullName: input.FullName,
		RoleID:   input.RoleID,
		IsActive: true,
	}

	// Guardamos en DB
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el usuario (posible duplicado)"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Usuario registrado exitosamente",
		"user":    user.ToResponse(),
	})
}

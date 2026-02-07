package controllers

import (
	"net/http"

	"github.com/cesarbmathec/medical-exams-backend/config"
	"github.com/cesarbmathec/medical-exams-backend/dtos"
	"github.com/cesarbmathec/medical-exams-backend/models"
	"github.com/cesarbmathec/medical-exams-backend/utils"
	"github.com/gin-gonic/gin"

	_ "github.com/cesarbmathec/medical-exams-backend/docs"
)

// Login godoc
// @Summary      Iniciar sesión
// @Description  Autentica al usuario y devuelve un token JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dtos.LoginRequest true "Credenciales de usuario"
// @Success      200 {object} map[string]interface{}
// @Failure      401 {object} map[string]string
// @Router       /login [post]
func Login(c *gin.Context) {
	var input dtos.LoginRequest
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

// Register godoc
// @Summary      Registrar nuevo usuario
// @Description  Crea un nuevo usuario en el sistema
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dtos.RegisterRequest true "Datos del nuevo usuario"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /register [post]
func Register(c *gin.Context) {
	var input dtos.RegisterRequest
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

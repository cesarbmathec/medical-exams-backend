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
// @Success      200 {object} utils.Response{data=dtos.LoginResponse}
// @Failure      400 {object} utils.Response{errors=string}
// @Failure      401 {object} utils.Response{errors=string}
// @Failure      500 {object} utils.Response{errors=string} "Error interno del servidor"
// @Router       /login [post]
func Login(c *gin.Context) {
	var input dtos.LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, "Error de validación", err.Error())
		return
	}

	var user models.User
	db := config.GetDB()

	// Buscar usuario e incluir el rol
	if err := db.Preload("Role").Where("username = ?", input.Username).First(&user).Error; err != nil {
		utils.Error(c, http.StatusUnauthorized, "Usuario o contraseña incorrectos", nil)
		return
	}

	// Verificar password usando el método que definiste en user.go
	if !user.CheckPassword(input.Password) {
		utils.Error(c, http.StatusUnauthorized, "error", "Usuario o contraseña incorrectos")
		return
	}

	// Generar Token
	token, err := utils.GenerateToken(user.ID, user.Username, user.RoleID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Inicio de sesión exitoso", dtos.LoginResponse{
		User:  user.ToResponse(),
		Token: token,
	})
}

// Register godoc
// @Summary      Registrar nuevo usuario
// @Description  Crea un nuevo usuario en el sistema
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dtos.RegisterRequest true "Datos del nuevo usuario"
// @Success      201 {object} utils.Response{data=dtos.RegisterResponse}
// @Success      200 {object} utils.Response{data=dtos.LoginResponse}
// @Failure      400 {object} utils.Response{errors=string}
// @Failure      500 {object} utils.Response{errors=string} "Error interno del servidor"
// @Router       /register [post]
func Register(c *gin.Context) {
	var input dtos.RegisterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, "error", err.Error())
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
		utils.Error(c, http.StatusInternalServerError, "error", "No se pudo crear el usuario (posible duplicado)")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.RoleID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.Success(c, http.StatusCreated,
		"Usuario registrado exitosamente",
		dtos.RegisterResponse{
			User:  user.ToResponse(),
			Token: token,
		})
}

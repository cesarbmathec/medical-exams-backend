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

// CreatePatient godoc
// @Summary      Crear paciente
// @Description  Crea un nuevo paciente en el sistema
// @Tags         patients
// @Accept       json
// @Produce      json
// @Param        request body dtos.CreatePatientRequest true "Datos para crear el paciente"
// @Success      201 {object} utils.Response{data=models.Patient}
// @Failure      400 {object} utils.Response{errors=string}
// @Failure      500 {object} utils.Response{errors=string}
// @Router       /patients [post]
// @Security BearerAuth
func CreatePatient(c *gin.Context) {
	var input dtos.CreatePatientRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, "Error de validación", err.Error())
		return
	}

	// Extraemos el ID del usuario del token (inyectado por el middleware)
	userID, _ := c.Get("userID")

	// Mapeamos los datos al modelo real de base de datos
	patient := models.Patient{
		DocumentType:   input.DocumentType,
		DocumentNumber: input.DocumentNumber,
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		DateOfBirth:    input.DateOfBirth,
		Gender:         input.Gender,
		Phone:          input.Phone,
		Email:          input.Email,
		BloodType:      input.BloodType,
		CreatedBy:      userID.(uint), // Asignamos el ID del usuario autenticado
	}

	db := config.GetDB()
	if err := db.Create(&patient).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "No se pudo crear el paciente", err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, "Paciente creado exitosamente", patient)
}

// GetPatients godoc
// @Summary      Listar pacientes
// @Description  Obtiene una lista de pacientes, con opción de filtrar por número de documento
// @Tags         patients
// @Accept       json
// @Produce      json
// @Param        document query string false "Número de documento para filtrar"
// @Success      200 {array} utils.Response{data=[]models.Patient}
// @Failure      500 {object} utils.Response{errors=string} "Error interno del servidor"
// @Router       /patients [get]
// @Security BearerAuth
func GetPatients(c *gin.Context) {
	var patients []models.Patient
	db := config.GetDB()

	// Filtro por número de documento si viene en el query
	doc := c.Query("document")
	query := db.Model(&models.Patient{})
	if doc != "" {
		query = query.Where("document_number LIKE ?", "%"+doc+"%")
	}

	if err := query.Find(&patients).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Error al obtener pacientes", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Pacientes obtenidos exitosamente", patients)
}

// GetPatientByID godoc
// @Summary      Obtener paciente por ID
// @Description  Obtiene los detalles de un paciente específico por su ID
// @Tags         patients
// @Accept       json
// @Produce      json
// @Param        id path int true "ID del paciente"
// @Success      200 {object} utils.Response{data=models.Patient}
// @Failure      404 {object} utils.Response{errors=string}
// @Router       /patients/{id} [get]
// @Security BearerAuth
func GetPatientByID(c *gin.Context) {
	id := c.Param("id")
	var patient models.Patient
	db := config.GetDB()

	if err := db.First(&patient, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Paciente no encontrado", nil)
		return
	}

	utils.Success(c, http.StatusOK, "Paciente obtenido exitosamente", patient)
}

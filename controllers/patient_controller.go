package controllers

import (
	"net/http"

	"github.com/cesarbmathec/medical-exams-backend/config"
	"github.com/cesarbmathec/medical-exams-backend/dtos"
	"github.com/cesarbmathec/medical-exams-backend/models"
	"github.com/gin-gonic/gin"
)

func CreatePatient(c *gin.Context) {
	var input dtos.CreatePatientRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, patient)
}

// GetPatients lista pacientes con filtros opcionales
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los pacientes"})
		return
	}

	c.JSON(http.StatusOK, patients)
}

// GetPatientByID busca un paciente específico
func GetPatientByID(c *gin.Context) {
	id := c.Param("id")
	var patient models.Patient
	db := config.GetDB()

	if err := db.First(&patient, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Paciente no encontrado"})
		return
	}

	c.JSON(http.StatusOK, patient)
}

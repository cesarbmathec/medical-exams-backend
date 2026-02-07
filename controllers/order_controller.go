package controllers

import (
	"net/http"

	"github.com/cesarbmathec/medical-exams-backend/config"
	"github.com/cesarbmathec/medical-exams-backend/dtos"
	"github.com/cesarbmathec/medical-exams-backend/models"
	"github.com/gin-gonic/gin"

	_ "github.com/cesarbmathec/medical-exams-backend/docs"
)

// CreateOrder godoc
// @Summary      Crear orden de examen
// @Description  Crea una nueva orden de examen para un paciente específico
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        request body dtos.CreateOrderRequest true "Datos para crear la orden"
// @Success      201 {object} models.Order
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /orders [post]
func CreateOrder(c *gin.Context) {
	var input dtos.CreateOrderRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	db := config.GetDB()

	// Iniciamos una Transacción para asegurar que se cree la orden Y sus exámenes
	tx := db.Begin()

	order := models.Order{
		PatientID:       input.PatientID,
		Priority:        input.Priority,
		ReferringDoctor: input.ReferringDoctor,
		Diagnosis:       input.Diagnosis,
		CreatedBy:       userID.(uint),
		Status:          "pendiente",
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear la orden"})
		return
	}

	// Crear cada examen dentro de la orden
	for _, exInput := range input.Exams {
		exam := models.OrderExam{
			OrderID:    order.ID,
			ExamTypeID: exInput.ExamTypeID,
			Price:      exInput.Price,
			Status:     "pendiente",
		}
		if err := tx.Create(&exam).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al asignar exámenes"})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusCreated, order)
}

// GetOrders godoc
// @Summary      Listar órdenes de examen
// @Description  Obtiene una lista de todas las órdenes de examen con sus detalles
// @Tags         orders
// @Accept       json
// @Produce      json
// @Success      200 {array} models.Order
// @Failure      500 {object} map[string]string
// @Router       /orders [get]
func GetOrders(c *gin.Context) {
	var orders []models.Order
	db := config.GetDB()

	// Preload carga los exámenes asociados a cada orden
	if err := db.Preload("Exams").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener órdenes"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

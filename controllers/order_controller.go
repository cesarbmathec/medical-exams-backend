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

// CreateOrder godoc
// @Summary      Crear orden de examen
// @Description  Crea una nueva orden de examen para un paciente específico
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        request body dtos.CreateOrderRequest true "Datos para crear la orden"
// @Success      201 {object} utils.Response{data=models.Order}
// @Failure      400 {object} utils.Response{errors=string}
// @Failure      500 {object} utils.Response{errors=string} "Error interno del servidor"
// @Router       /orders [post]
// @Security BearerAuth
func CreateOrder(c *gin.Context) {
	var input dtos.CreateOrderRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, "Error de validación", err.Error())
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
		utils.Error(c, http.StatusInternalServerError, "No se pudo crear la orden", err.Error())
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
			utils.Error(c, http.StatusInternalServerError, "No se pudo crear el examen", err.Error())
			return
		}
	}

	tx.Commit()
	utils.Success(c, http.StatusCreated, "Orden creada exitosamente", order)
}

// GetOrders godoc
// @Summary      Listar órdenes con filtros
// @Description  Obtiene órdenes filtradas por rango de fechas, estado, prioridad o paciente
// @Tags         orders
// @Security     BearerAuth
// @Param        status query string false "Estado (pendiente, completado, cancelado)"
// @Param        priority query string false "Prioridad (normal, urgente, stat)"
// @Param        start_date query string false "Fecha inicio (YYYY-MM-DD)"
// @Param        end_date query string false "Fecha fin (YYYY-MM-DD)"
// @Param        patient_id query int false "ID del Paciente"
// @Success      200 {array} utils.Response{data=[]models.Order}
// @Failure      500 {object} utils.Response{errors=string} "Error interno del servidor"
// @Router       /orders [get]
func GetOrders(c *gin.Context) {
	db := config.GetDB()
	var orders []models.Order

	// Iniciamos el query cargando la relación con el paciente para el Front
	query := db.Preload("Patient").Preload("OrderExams.ExamType")

	// Aplicar filtros dinámicos
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if priority := c.Query("priority"); priority != "" {
		query = query.Where("priority = ?", priority)
	}
	if pID := c.Query("patient_id"); pID != "" {
		query = query.Where("patient_id = ?", pID)
	}

	// Filtro por rango de fechas
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate != "" && endDate != "" {
		query = query.Where("order_date BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59")
	}

	if err := query.Order("created_at DESC").Find(&orders).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Error al obtener órdenes", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Órdenes obtenidas exitosamente", orders)
}

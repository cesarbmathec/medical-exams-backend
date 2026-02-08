package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/cesarbmathec/medical-exams-backend/config"
	"github.com/cesarbmathec/medical-exams-backend/dtos"
	"github.com/cesarbmathec/medical-exams-backend/models"
	"github.com/cesarbmathec/medical-exams-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetOrderExamDetails godoc
// @Summary      Detalle de un examen específico de una orden
// @Description  Retorna el examen con sus parámetros y resultados previos (si existen)
// @Tags         lab
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID del examen dentro de la orden"
// @Success      200 {object} utils.Response{data=models.OrderExam}
// @Failure      404 {object} utils.Response{errors=string}
// @Router       /lab/exams/{id} [get]
func GetOrderExamDetails(c *gin.Context) {
	id := c.Param("id")
	var orderExam models.OrderExam
	db := config.GetDB()

	// Cargamos el examen, su tipo, y los parámetros definidos para ese tipo
	if err := db.Preload("ExamType.Parameters").
		Preload("Results").
		First(&orderExam, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Examen no encontrado", nil)
		return
	}

	utils.Success(c, http.StatusOK, "Examen obtenido exitosamente", orderExam)
}

// UpdateExamStatus godoc
// @Summary      Actualizar estado de un examen (Toma de muestra / Análisis)
// @Tags         lab
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Order Exam ID"
// @Param        status body string true "Nuevo estado: muestra_tomada, en_analisis"
// @Success      200 {object} utils.Response{data=nil}
// @Failure      400 {object} utils.Response{errors=string}
// @Failure      500 {object} utils.Response{errors=string} "Error interno del servidor"
// @Router       /lab/exams/{id}/status [patch]
func UpdateExamStatus(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Status string `json:"status" binding:"required,oneof=pendiente muestra_tomada en_analisis"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, "Error de validación", err.Error())
		return
	}

	userID, _ := c.Get("userID")
	db := config.GetDB()

	updates := map[string]interface{}{
		"status": input.Status,
	}

	if input.Status == "muestra_tomada" {
		now := time.Now()
		updates["sample_collected_at"] = &now
		updates["sample_collected_by"] = userID.(uint)
	}

	if err := db.Model(&models.OrderExam{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Error al actualizar estado", err.Error())
		return
	}
	utils.Success(c, http.StatusOK, "Estado actualizado exitosamente", nil)
}

// ValidateResults godoc
// @Summary      Validar resultados de un examen
// @Description  Marca los resultados como validados y finaliza el examen para su impresión
// @Tags         lab
// @Param        id path int true "ID del examen de la orden"
// @Success      200 {object} utils.Response{data=nil}
// @Failure      500 {object} utils.Response{errors=string} "Error interno del servidor"
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Router       /lab/exams/{id}/validate [post]
func ValidateResults(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	db := config.GetDB()

	now := time.Now()
	err := db.Transaction(func(tx *gorm.DB) error {
		// 1. Marcar el examen de la orden como completado y validado
		if err := tx.Model(&models.OrderExam{}).Where("id = ?", id).Updates(map[string]interface{}{
			"status":       "completado",
			"validated_at": &now,
			"validated_by": userID.(uint),
		}).Error; err != nil {
			return err
		}

		// 2. Opcional: Verificar si todos los exámenes de la ORDEN están listos
		// Si todos están 'completado', marcar la Orden principal como 'completada'
		return nil
	})

	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Error al validar resultados", err.Error())
		return
	}
	utils.Success(c, http.StatusOK, "Resultados validados exitosamente", nil)
}

// GetExamCatalog godoc
// @Summary      Catálogo de exámenes
// @Description  Obtiene la lista de tipos de exámenes disponibles con sus categorías y parámetros
// @Tags         lab
// @Accept       json
// @Produce      json
// @Success      200 {array} utils.Response{data=[]models.ExamType}
// @Failure      500 {object} utils.Response{errors=string} "Error interno del servidor"
// @Router       /lab/exams/catalog [get]
// @Security BearerAuth
func GetExamCatalog(c *gin.Context) {
	var exams []models.ExamType
	db := config.GetDB()

	// Preload carga la categoría y los parámetros de cada examen
	if err := db.Preload("Category").Find(&exams).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Error al obtener el catálogo", err.Error())
		return
	}
	utils.Success(c, http.StatusOK, "Catálogo obtenido exitosamente", exams)
}

// SubmitResults godoc
// @Summary      Registrar resultados de examen
// @Description  Permite a un técnico de laboratorio registrar los resultados de un examen específico dentro de una orden
// @Tags         lab
// @Accept       json
// @Produce      json
// @Param        id path int true "ID del OrderExam"
// @Param        request body []dtos.UpdateResultRequest true "Resultados a registrar"
// @Success      200 {object} utils.Response{data=nil}
// @Failure      400 {object} utils.Response{errors=string}
// @Router       /lab/exams/{id}/results [post]
// @Security BearerAuth
func SubmitResults(c *gin.Context) {
	orderExamIDParam := c.Param("id")
	var input []dtos.UpdateResultRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, "error", err.Error())
		return
	}

	orderExamID, err := parseUint(orderExamIDParam)
	if err != nil || orderExamID == 0 {
		utils.Error(c, http.StatusBadRequest, "ID de examen inválido", nil)
		return
	}

	userID, _ := c.Get("userID")
	db := config.GetDB()

	err = db.Transaction(func(tx *gorm.DB) error {
		var orderExam models.OrderExam
		if err := tx.First(&orderExam, orderExamID).Error; err != nil {
			return err
		}

		for _, res := range input {
			result := models.ExamResult{
				OrderExamID:     orderExamID,
				ExamParameterID: res.ParameterID,
				ValueNumeric:    res.ValueNumeric,
				ValueText:       res.ValueText,
				EnteredBy:       userID.(uint),
			}
			if err := tx.Create(&result).Error; err != nil {
				return err
			}
		}

		if err := tx.Model(&models.OrderExam{}).Where("id = ?", orderExamID).Update("status", "completado").Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Error al registrar resultados", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Resultados registrados exitosamente", nil)
}

// Función auxiliar para convertir string a uint
func parseUint(s string) (uint, error) {
	parsed, err := strconv.ParseUint(s, 10, 0)
	if err != nil {
		return 0, err
	}
	return uint(parsed), nil
}

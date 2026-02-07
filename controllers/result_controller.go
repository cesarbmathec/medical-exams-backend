package controllers

import (
	"fmt"
	"net/http"

	"github.com/cesarbmathec/medical-exams-backend/config"
	"github.com/cesarbmathec/medical-exams-backend/dtos"
	"github.com/cesarbmathec/medical-exams-backend/models"
	"github.com/gin-gonic/gin"

	_ "github.com/cesarbmathec/medical-exams-backend/docs"
)

// SubmitResults godoc
// @Summary      Registrar resultados de examen
// @Description  Permite a un técnico de laboratorio registrar los resultados de un examen específico dentro de una orden
// @Tags         results
// @Accept       json
// @Produce      json
// @Param        id path int true "ID del OrderExam"
// @Param        request body []dtos.UpdateResultRequest true "Resultados a registrar"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /results/{id} [post]
func SubmitResults(c *gin.Context) {
	orderExamID := c.Param("id")
	var input []dtos.UpdateResultRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	db := config.GetDB()

	for _, res := range input {
		result := models.ExamResult{
			OrderExamID:     parseUint(orderExamID), // Función auxiliar para convertir string a uint
			ExamParameterID: res.ParameterID,
			ValueNumeric:    res.ValueNumeric,
			ValueText:       res.ValueText,
			EnteredBy:       userID.(uint),
		}
		db.Save(&result)
	}

	// Actualizar estado del examen a 'completado'
	db.Model(&models.OrderExam{}).Where("id = ?", orderExamID).Update("status", "completado")

	c.JSON(http.StatusOK, gin.H{"message": "Resultados registrados exitosamente"})
}

// Función auxiliar para convertir string a uint
func parseUint(s string) uint {
	var num uint
	fmt.Sscanf(s, "%d", &num)
	return num
}

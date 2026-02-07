package controllers

import (
	"net/http"

	"github.com/cesarbmathec/medical-exams-backend/config"
	"github.com/cesarbmathec/medical-exams-backend/models"
	"github.com/gin-gonic/gin"

	_ "github.com/cesarbmathec/medical-exams-backend/docs"
)

// GetExamCatalog godoc
// @Summary      Catálogo de exámenes
// @Description  Obtiene la lista de tipos de exámenes disponibles con sus categorías y parámetros
// @Tags         exams
// @Accept       json
// @Produce      json
// @Success      200 {array} models.ExamType
// @Failure      500 {object} map[string]string
// @Router       /exams/catalog [get]
func GetExamCatalog(c *gin.Context) {
	var exams []models.ExamType
	db := config.GetDB()

	// Preload carga la categoría y los parámetros de cada examen
	if err := db.Preload("Category").Find(&exams).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener catálogo"})
		return
	}
	c.JSON(http.StatusOK, exams)
}

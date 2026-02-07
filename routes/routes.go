package routes

// Aquí definimos todas las rutas de la API y los controladores asociados

import (
	"github.com/cesarbmathec/medical-exams-backend/controllers"
	"github.com/cesarbmathec/medical-exams-backend/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// --- RUTAS PÚBLICAS ---
	api := r.Group("/api")
	{
		api.POST("/login", controllers.Login)
		api.POST("/register", controllers.Register)
	}

	// --- RUTAS PROTEGIDAS (Requieren Token) ---
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// Ejemplo: Perfil del usuario actual
		protected.GET("/me", func(c *gin.Context) {
			userID, _ := c.Get("userID")
			c.JSON(200, gin.H{"user_id": userID, "message": "Acceso concedido"})
		})

		// RUTAS DE PACIENTES
		patients := protected.Group("/patients")
		{
			patients.POST("/", controllers.CreatePatient)    // Registrar
			patients.GET("/", controllers.GetPatients)       // Listar/Buscar
			patients.GET("/:id", controllers.GetPatientByID) // Ver detalle
		}

		// Catálogo
		protected.GET("/exams/catalog", controllers.GetExamCatalog)

		// Órdenes
		protected.POST("/orders", controllers.CreateOrder)
		protected.GET("/orders", controllers.GetOrders)

		// Resultados (Ruta para Bioanalistas)
		protected.POST("/exams/:id/results", controllers.SubmitResults)
	}
	return r
}

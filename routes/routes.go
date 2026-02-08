package routes

// Aquí definimos todas las rutas de la API y los controladores asociados

import (
	"github.com/cesarbmathec/medical-exams-backend/controllers"
	"github.com/cesarbmathec/medical-exams-backend/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/cesarbmathec/medical-exams-backend/docs"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(middleware.RecoveryMiddleware())

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// --- RUTAS PÚBLICAS ---
	api := r.Group("/api/v1")
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

		// Órdenes
		orders := protected.Group("/orders")
		{
			orders.POST("/", controllers.CreateOrder)
			orders.GET("/", controllers.GetOrders)
		}

		lab := protected.Group("/lab")
		{
			lab.GET("/exams/:id", controllers.GetOrderExamDetails)
			lab.PATCH("/exams/:id/status", controllers.UpdateExamStatus)
			lab.POST("/exams/:id/validate", controllers.ValidateResults) // Nueva ruta para validar resultados
			lab.POST("/exams/:id/results", controllers.SubmitResults)
			lab.GET("/exams/catalog", controllers.GetExamCatalog) // Para que los bioanalistas puedan ver el catálogo de exámenes y sus parámetros
		}
	}
	return r
}

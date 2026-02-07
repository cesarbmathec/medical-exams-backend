// @title           Laboratorio Clínico API
// @version         1.0
// @description     API para gestión de exámenes, pacientes y resultados.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Soporte Técnico
// @contact.email  soporte@laboratorio.com

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Escribe 'Bearer ' seguido de tu token JWT

package main

import (
	"log"
	"os"

	"github.com/cesarbmathec/medical-exams-backend/config"
	"github.com/cesarbmathec/medical-exams-backend/controllers"
	"github.com/cesarbmathec/medical-exams-backend/middleware"
	"github.com/cesarbmathec/medical-exams-backend/migrations"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	// Swagger

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/cesarbmathec/medical-exams-backend/docs"
)

func main() {
	// Cargamos variables de entorno (.env)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando el archivo .env")
	}

	// Conectamos a la base de datos
	config.ConnectDatabase()
	db := config.GetDB()

	// Ejecutamos Migraciones y Seeding
	migrations.RunMigrations(db)

	// Crear instancia de Gin
	r := gin.Default()

	// Configurar CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("🚀 Servidor corriendo en http://localhost:" + port)
	r.Run(":" + port)
}

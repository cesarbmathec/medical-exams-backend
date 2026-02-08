// @title           Laboratorio Clínico API
// @version         1.0
// @description     API para gestión de exámenes, pacientes y resultados.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Soporte Técnico
// @contact.email  soporte@laboratorio.com

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Escribe 'Bearer ' seguido de tu token JWT

package main

import (
	"log"
	"os"

	"github.com/cesarbmathec/medical-exams-backend/config"
	"github.com/cesarbmathec/medical-exams-backend/migrations"
	"github.com/cesarbmathec/medical-exams-backend/routes"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"

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

	// Usamos el router definido en routes.go
	r := routes.SetupRouter()

	// Configurar CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("🚀 Servidor corriendo en http://localhost:" + port)
	r.Run(":" + port)
}

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
	"strings"

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
	allowedOrigins := parseCSVEnv("CORS_ALLOWED_ORIGINS")
	if len(allowedOrigins) == 0 && os.Getenv("GIN_MODE") != "release" {
		allowedOrigins = []string{"http://localhost:3000", "http://localhost:5173"}
	}
	if len(allowedOrigins) == 0 && os.Getenv("GIN_MODE") == "release" {
		log.Fatal("CORS_ALLOWED_ORIGINS requerido en GIN_MODE=release")
	}

	allowCredentials := true
	if strings.EqualFold(os.Getenv("CORS_ALLOW_CREDENTIALS"), "false") {
		allowCredentials = false
	}
	for _, origin := range allowedOrigins {
		if origin == "*" {
			allowCredentials = false
			break
		}
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: allowCredentials,
	}))

	if os.Getenv("GIN_MODE") == "release" {
		r.SetTrustedProxies(parseCSVEnv("TRUSTED_PROXIES"))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("🚀 Servidor corriendo en http://localhost:" + port)
	r.Run(":" + port)
}

func parseCSVEnv(key string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item != "" {
			out = append(out, item)
		}
	}
	return out
}

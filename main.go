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
	"fmt"
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
		allowedOrigins = []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:34115"}
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

	fmt.Printf("log.Logger AllowCredentials: %v\n", allowCredentials)
	fmt.Printf("log.Logger AllowedOrigins: %v\n", allowedOrigins)

	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			// Esto permite localhost y el esquema especial de Wails, y también los orígenes del .env
			if contains(allowedOrigins, origin) {
				return true
			}
			return origin == "wails://wails.localhost:34115" ||
				origin == "wails://wails.127.0.0.1:34115" || origin == "http://localhost:34115" || origin == "http://127.0.0.1:5173"
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
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

// contains checks if a string is present in a slice of strings.
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

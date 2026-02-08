package migrations

import (
	"log"
	"os"
	"strings"

	"github.com/cesarbmathec/medical-exams-backend/models"

	"gorm.io/gorm"
)

// RunMigrations ejecuta todas las migraciones
func RunMigrations(db *gorm.DB) {
	log.Println("🔄 Running database migrations...")

	// Orden de migraciones (importante por las foreign keys)
	err := db.AutoMigrate(
		// Primero las tablas sin dependencias
		&models.Role{},
		&models.ExamCategory{},
		&models.SampleType{},

		// Luego las que dependen de las anteriores
		&models.User{},
		&models.Patient{},
		&models.ExamType{},

		// Luego las que tienen más dependencias
		&models.ExamParameter{},
		&models.Order{},

		// Finalmente las tablas dependientes
		&models.OrderExam{},
		&models.ExamResult{},
		&models.Payment{},
		&models.Invoice{},
		&models.AuditLog{},
		&models.Reagent{},
		&models.Equipment{},
	)

	if err != nil {
		log.Fatal("❌ Migration failed:", err)
	}

	log.Println("✅ Migrations completed successfully!")

	// Crear datos iniciales
	seedData(db)
}

// seedData crea los datos iniciales del sistema
func seedData(db *gorm.DB) {
	if !shouldSeed() {
		log.Println("🌱 Seeding skipped (disabled for current environment)")
		return
	}

	log.Println("🌱 Seeding initial data...")

	// Crear roles si no existen
	createRoles(db)
	createAdminUser(db)
	createExamCategories(db)
	createSampleTypes(db)

	log.Println("✅ Seeding completed!")
}

func shouldSeed() bool {
	if strings.EqualFold(os.Getenv("SEED_DB"), "false") {
		return false
	}

	if strings.EqualFold(os.Getenv("SEED_DB"), "true") {
		return true
	}

	ginMode := strings.ToLower(os.Getenv("GIN_MODE"))
	if ginMode == "release" || ginMode == "production" {
		return false
	}

	return true
}

func createRoles(db *gorm.DB) {
	roles := []models.Role{
		{
			Name:        "admin",
			Description: "Administrador del sistema",
			Permissions: models.Permissions{"all": {"*"}},
			IsActive:    true,
		},
		{
			Name:        "bioanalista",
			Description: "Bioanalista - Análisis y validación",
			Permissions: models.Permissions{
				"orders":   {"read"},
				"results":  {"read", "write"},
				"patients": {"read"},
			},
			IsActive: true,
		},
		{
			Name:        "recepcionista",
			Description: "Recepcionista - Gestión de pacientes y órdenes",
			Permissions: models.Permissions{
				"patients": {"read", "write"},
				"orders":   {"read", "write"},
				"payments": {"read", "write"},
			},
			IsActive: true,
		},
	}

	for _, role := range roles {
		var existingRole models.Role
		if err := db.Where("name = ?", role.Name).First(&existingRole).Error; err == gorm.ErrRecordNotFound {
			db.Create(&role)
			log.Printf("  ✓ Created role: %s\n", role.Name)
		}
	}
}

func createAdminUser(db *gorm.DB) {
	var adminRole models.Role
	db.Where("name = ?", "admin").First(&adminRole)

	adminUser := models.User{
		Username: "admin",
		Email:    "admin@laboratorio.com",
		Password: "Admin123!",
		FullName: "Administrador del Sistema",
		RoleID:   adminRole.ID,
		IsActive: true,
	}

	var existingUser models.User
	if err := db.Where("username = ?", adminUser.Username).First(&existingUser).Error; err == gorm.ErrRecordNotFound {
		db.Create(&adminUser)
		log.Println("  ✓ Created admin user (username: admin, password: Admin123!)")
	}
}

func createExamCategories(db *gorm.DB) {
	categories := []models.ExamCategory{
		{Name: "Hematología", Code: "HEM", Description: "Estudios de sangre y componentes", DisplayOrder: 1},
		{Name: "Química Sanguínea", Code: "QS", Description: "Análisis bioquímicos en sangre", DisplayOrder: 2},
		{Name: "Serología", Code: "SER", Description: "Detección de anticuerpos y antígenos", DisplayOrder: 3},
	}

	for _, cat := range categories {
		var existing models.ExamCategory
		if err := db.Where("code = ?", cat.Code).First(&existing).Error; err == gorm.ErrRecordNotFound {
			db.Create(&cat)
			log.Printf("  ✓ Created category: %s\n", cat.Name)
		}
	}
}

func createSampleTypes(db *gorm.DB) {
	types := []models.SampleType{
		{Name: "Sangre Venosa", Description: "Muestra de sangre obtenida por venopunción"},
		{Name: "Orina", Description: "Muestra de orina"},
		{Name: "Heces", Description: "Muestra de materia fecal"},
	}

	for _, sampleType := range types {
		var existing models.SampleType
		if err := db.Where("name = ?", sampleType.Name).First(&existing).Error; err == gorm.ErrRecordNotFound {
			db.Create(&sampleType)
			log.Printf("  ✓ Created sample type: %s\n", sampleType.Name)
		}
	}
}

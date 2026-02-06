package models

import "time"

// Patient representa un paciente del laboratorio
type Patient struct {
	BaseModel
	DocumentType          string    `gorm:"size:20;not null" json:"document_type" binding:"required,oneof=cedula pasaporte rif otro"`
	DocumentNumber        string    `gorm:"size:50;not null" json:"document_number" binding:"required"`
	FirstName             string    `gorm:"size:100;not null" json:"first_name" binding:"required"`
	LastName              string    `gorm:"size:100;not null" json:"last_name" binding:"required"`
	DateOfBirth           time.Time `gorm:"type:date;not null" json:"date_of_birth" binding:"required"`
	Gender                string    `gorm:"size:1" json:"gender" binding:"omitempty,oneof=M F O"`
	Phone                 string    `gorm:"size:20" json:"phone"`
	Email                 string    `gorm:"size:100" json:"email" binding:"omitempty,email"`
	Address               string    `gorm:"type:text" json:"address"`
	City                  string    `gorm:"size:100" json:"city"`
	State                 string    `gorm:"size:100" json:"state"`
	Country               string    `gorm:"size:100;default:'Venezuela'" json:"country"`
	EmergencyContactName  string    `gorm:"size:150" json:"emergency_contact_name"`
	EmergencyContactPhone string    `gorm:"size:20" json:"emergency_contact_phone"`
	BloodType             string    `gorm:"size:5" json:"blood_type" binding:"omitempty,oneof=A+ A- B+ B- AB+ AB- O+ O-"`
	Allergies             string    `gorm:"type:text" json:"allergies"`
	MedicalConditions     string    `gorm:"type:text" json:"medical_conditions"`
	IsActive              bool      `gorm:"default:true" json:"is_active"`
	CreatedBy             uint      `json:"created_by"`

	// Relaciones
	Creator User    `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Orders  []Order `gorm:"foreignKey:PatientID" json:"orders,omitempty"`
}

// TableName especifica el nombre de la tabla
func (Patient) TableName() string {
	return "patients"
}

// GetAge calcula la edad del paciente
func (p *Patient) GetAge() int {
	now := time.Now()
	age := now.Year() - p.DateOfBirth.Year()
	if now.YearDay() < p.DateOfBirth.YearDay() {
		age--
	}
	return age
}

// GetFullName retorna el nombre completo
func (p *Patient) GetFullName() string {
	return p.FirstName + " " + p.LastName
}

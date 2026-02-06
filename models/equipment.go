package models

import (
	"time"

	"gorm.io/gorm"
)

// Equipment representa un equipo del laboratorio
type Equipment struct {
	BaseModel
	Name                     string     `gorm:"size:200;not null" json:"name" binding:"required"`
	Code                     string     `gorm:"size:50;uniqueIndex;not null" json:"code" binding:"required"`
	Description              string     `gorm:"type:text" json:"description"`
	Manufacturer             string     `gorm:"size:150" json:"manufacturer"`
	Model                    string     `gorm:"size:100" json:"model"`
	SerialNumber             string     `gorm:"size:100;uniqueIndex" json:"serial_number"`
	PurchaseDate             *time.Time `gorm:"type:date" json:"purchase_date"`
	WarrantyExpiration       *time.Time `gorm:"type:date" json:"warranty_expiration"`
	LastMaintenanceDate      *time.Time `gorm:"type:date" json:"last_maintenance_date"`
	NextMaintenanceDate      *time.Time `gorm:"type:date" json:"next_maintenance_date"`
	MaintenanceFrequencyDays int        `json:"maintenance_frequency_days"`
	Status                   string     `gorm:"size:30;default:'operativo'" json:"status"`
	Location                 string     `gorm:"size:100" json:"location"`
	ResponsibleUserID        *uint      `json:"responsible_user_id"`
	Notes                    string     `gorm:"type:text" json:"notes"`

	// Relaciones
	ResponsibleUser    *User                  `gorm:"foreignKey:ResponsibleUserID" json:"responsible_user,omitempty"`
	MaintenanceHistory []EquipmentMaintenance `gorm:"foreignKey:EquipmentID" json:"maintenance_history,omitempty"`
}

// TableName especifica el nombre de la tabla
func (Equipment) TableName() string {
	return "equipment"
}

// IsOperational verifica si el equipo está operativo
func (e *Equipment) IsOperational() bool {
	return e.Status == "operativo"
}

// NeedsMaintenance verifica si necesita mantenimiento
func (e *Equipment) NeedsMaintenance() bool {
	if e.NextMaintenanceDate == nil {
		return false
	}
	return time.Now().After(*e.NextMaintenanceDate)
}

// IsUnderWarranty verifica si está bajo garantía
func (e *Equipment) IsUnderWarranty() bool {
	if e.WarrantyExpiration == nil {
		return false
	}
	return time.Now().Before(*e.WarrantyExpiration)
}

// ScheduleNextMaintenance programa el próximo mantenimiento
func (e *Equipment) ScheduleNextMaintenance() {
	if e.MaintenanceFrequencyDays > 0 {
		nextDate := time.Now().AddDate(0, 0, e.MaintenanceFrequencyDays)
		e.NextMaintenanceDate = &nextDate
	}
}

// GetMaintenanceStatus retorna el estado de mantenimiento
func (e *Equipment) GetMaintenanceStatus() string {
	if e.NextMaintenanceDate == nil {
		return "no_programado"
	}
	if e.NeedsMaintenance() {
		return "vencido"
	}

	// Verificar si está próximo (7 días)
	sevenDaysFromNow := time.Now().AddDate(0, 0, 7)
	if e.NextMaintenanceDate.Before(sevenDaysFromNow) {
		return "proximo"
	}

	return "al_dia"
}

// EquipmentMaintenance representa el historial de mantenimiento
type EquipmentMaintenance struct {
	BaseModel
	EquipmentID         uint       `gorm:"not null" json:"equipment_id" binding:"required"`
	MaintenanceType     string     `gorm:"size:30;not null" json:"maintenance_type" binding:"required,oneof=preventivo correctivo calibracion verificacion"`
	MaintenanceDate     time.Time  `gorm:"type:date;not null" json:"maintenance_date" binding:"required"`
	PerformedBy         string     `gorm:"size:150" json:"performed_by"`
	TechnicianCompany   string     `gorm:"size:150" json:"technician_company"`
	Description         string     `gorm:"type:text" json:"description"`
	Findings            string     `gorm:"type:text" json:"findings"`
	ActionsTaken        string     `gorm:"type:text" json:"actions_taken"`
	Cost                float64    `gorm:"type:decimal(10,2)" json:"cost"`
	NextMaintenanceDate *time.Time `gorm:"type:date" json:"next_maintenance_date"`
	CreatedBy           *uint      `json:"created_by"`

	// Relaciones
	Equipment Equipment `gorm:"foreignKey:EquipmentID" json:"equipment,omitempty"`
	Creator   *User     `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// TableName especifica el nombre de la tabla
func (EquipmentMaintenance) TableName() string {
	return "equipment_maintenance"
}

// IsPreventive verifica si es mantenimiento preventivo
func (em *EquipmentMaintenance) IsPreventive() bool {
	return em.MaintenanceType == "preventivo"
}

// IsCorrective verifica si es mantenimiento correctivo
func (em *EquipmentMaintenance) IsCorrective() bool {
	return em.MaintenanceType == "correctivo"
}

// AfterCreate hook para actualizar el equipo después de crear mantenimiento
func (em *EquipmentMaintenance) AfterCreate(tx *gorm.DB) error {
	// Actualizar la fecha de último mantenimiento del equipo
	if err := tx.Model(&Equipment{}).Where("id = ?", em.EquipmentID).Updates(map[string]interface{}{
		"last_maintenance_date": em.MaintenanceDate,
		"next_maintenance_date": em.NextMaintenanceDate,
	}).Error; err != nil {
		return err
	}
	return nil
}

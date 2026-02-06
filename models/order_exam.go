package models

import (
	"time"

	"gorm.io/gorm"
)

// OrderExam representa un examen específico dentro de una orden
type OrderExam struct {
	BaseModel
	OrderID           uint       `gorm:"not null" json:"order_id" binding:"required"`
	ExamTypeID        uint       `gorm:"not null" json:"exam_type_id" binding:"required"`
	Status            string     `gorm:"size:20;not null;default:'pendiente'" json:"status"`
	SampleCollectedAt *time.Time `json:"sample_collected_at"`
	SampleCollectedBy *uint      `json:"sample_collected_by"`
	SampleBarcode     string     `gorm:"size:100" json:"sample_barcode"`
	AnalyzedAt        *time.Time `json:"analyzed_at"`
	AnalyzedBy        *uint      `json:"analyzed_by"`
	ValidatedAt       *time.Time `json:"validated_at"`
	ValidatedBy       *uint      `json:"validated_by"`
	Price             float64    `gorm:"type:decimal(10,2);not null" json:"price" binding:"required,gt=0"`
	Discount          float64    `gorm:"type:decimal(10,2);default:0" json:"discount"`
	FinalPrice        float64    `gorm:"type:decimal(10,2);not null" json:"final_price"`
	Notes             string     `gorm:"type:text" json:"notes"`
	RejectionReason   string     `gorm:"type:text" json:"rejection_reason"`

	// Relaciones
	Order           Order        `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	ExamType        ExamType     `gorm:"foreignKey:ExamTypeID" json:"exam_type,omitempty"`
	CollectedByUser *User        `gorm:"foreignKey:SampleCollectedBy" json:"collected_by_user,omitempty"`
	AnalyzedByUser  *User        `gorm:"foreignKey:AnalyzedBy" json:"analyzed_by_user,omitempty"`
	ValidatedByUser *User        `gorm:"foreignKey:ValidatedBy" json:"validated_by_user,omitempty"`
	Results         []ExamResult `gorm:"foreignKey:OrderExamID" json:"results,omitempty"`
}

// TableName especifica el nombre de la tabla
func (OrderExam) TableName() string {
	return "order_exams"
}

// BeforeCreate calcula el precio final antes de crear
func (oe *OrderExam) BeforeCreate(tx *gorm.DB) error {
	if oe.FinalPrice == 0 {
		oe.FinalPrice = oe.Price - oe.Discount
	}
	return nil
}

// BeforeUpdate calcula el precio final antes de actualizar
func (oe *OrderExam) BeforeUpdate(tx *gorm.DB) error {
	oe.FinalPrice = oe.Price - oe.Discount
	return nil
}

// IsCompleted verifica si el examen está completado
func (oe *OrderExam) IsCompleted() bool {
	return oe.Status == "completado" && oe.ValidatedAt != nil
}

// IsPending verifica si el examen está pendiente
func (oe *OrderExam) IsPending() bool {
	return oe.Status == "pendiente"
}

// CanBeAnalyzed verifica si el examen puede ser analizado
func (oe *OrderExam) CanBeAnalyzed() bool {
	return oe.Status == "muestra_tomada" && oe.SampleCollectedAt != nil
}

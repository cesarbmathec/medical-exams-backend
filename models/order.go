package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Order representa una orden de exámenes
type Order struct {
	BaseModel
	OrderNumber        string     `gorm:"size:50;uniqueIndex;not null" json:"order_number"`
	PatientID          uint       `gorm:"not null" json:"patient_id" binding:"required"`
	OrderDate          time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"order_date"`
	Status             string     `gorm:"size:20;not null;default:'pendiente'" json:"status"`
	Priority           string     `gorm:"size:20;default:'normal'" json:"priority" binding:"omitempty,oneof=normal urgente stat"`
	ReferringDoctor    string     `gorm:"size:150" json:"referring_doctor"`
	DoctorPhone        string     `gorm:"size:20" json:"doctor_phone"`
	Diagnosis          string     `gorm:"type:text" json:"diagnosis"`
	ClinicalNotes      string     `gorm:"type:text" json:"clinical_notes"`
	Subtotal           float64    `gorm:"type:decimal(10,2);default:0" json:"subtotal"`
	DiscountPercentage float64    `gorm:"type:decimal(5,2);default:0" json:"discount_percentage"`
	DiscountAmount     float64    `gorm:"type:decimal(10,2);default:0" json:"discount_amount"`
	TaxPercentage      float64    `gorm:"type:decimal(5,2);default:0" json:"tax_percentage"`
	TaxAmount          float64    `gorm:"type:decimal(10,2);default:0" json:"tax_amount"`
	TotalAmount        float64    `gorm:"type:decimal(10,2);default:0" json:"total_amount"`
	PaidAmount         float64    `gorm:"type:decimal(10,2);default:0" json:"paid_amount"`
	Balance            float64    `gorm:"type:decimal(10,2);default:0" json:"balance"`
	PaymentStatus      string     `gorm:"size:20;default:'pendiente'" json:"payment_status"`
	CreatedBy          uint       `gorm:"not null" json:"created_by"`
	CompletedAt        *time.Time `json:"completed_at"`
	CancelledAt        *time.Time `json:"cancelled_at"`
	CancelledBy        *uint      `json:"cancelled_by"`
	CancellationReason string     `gorm:"type:text" json:"cancellation_reason"`

	// Relaciones
	Patient    Patient     `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	Creator    User        `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	OrderExams []OrderExam `gorm:"foreignKey:OrderID" json:"order_exams,omitempty"`
	Payments   []Payment   `gorm:"foreignKey:OrderID" json:"payments,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}

// BeforeCreate genera el número de orden automáticamente
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.OrderNumber == "" {
		// Formato: ORD-YYYYMMDD-000001
		now := time.Now()
		dateStr := now.Format("20060102")

		// Contar órdenes del día
		var count int64
		tx.Model(&Order{}).Where("order_number LIKE ?", "ORD-"+dateStr+"%").Count(&count)

		o.OrderNumber = fmt.Sprintf("ORD-%s-%06d", dateStr, count+1)
	}
	return nil
}

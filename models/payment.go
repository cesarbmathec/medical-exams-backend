package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Payment representa un pago realizado
type Payment struct {
	BaseModel
	PaymentNumber      string     `gorm:"size:50;uniqueIndex;not null" json:"payment_number"`
	OrderID            uint       `gorm:"not null" json:"order_id" binding:"required"`
	PaymentDate        time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"payment_date"`
	Amount             float64    `gorm:"type:decimal(10,2);not null" json:"amount" binding:"required,gt=0"`
	PaymentMethod      string     `gorm:"size:30;not null" json:"payment_method" binding:"required,oneof=efectivo tarjeta_debito tarjeta_credito transferencia pago_movil cheque otro"`
	ReferenceNumber    string     `gorm:"size:100" json:"reference_number"`
	BankName           string     `gorm:"size:100" json:"bank_name"`
	CardLastDigits     string     `gorm:"size:4" json:"card_last_digits"`
	Status             string     `gorm:"size:20;default:'aprobado'" json:"status"`
	Notes              string     `gorm:"type:text" json:"notes"`
	CreatedBy          uint       `gorm:"not null" json:"created_by"`
	ApprovedBy         *uint      `json:"approved_by"`
	ApprovedAt         *time.Time `json:"approved_at"`
	CancelledBy        *uint      `json:"cancelled_by"`
	CancelledAt        *time.Time `json:"cancelled_at"`
	CancellationReason string     `gorm:"type:text" json:"cancellation_reason"`

	// Relaciones
	Order           Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Creator         User  `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	ApprovedByUser  *User `gorm:"foreignKey:ApprovedBy" json:"approved_by_user,omitempty"`
	CancelledByUser *User `gorm:"foreignKey:CancelledBy" json:"cancelled_by_user,omitempty"`
}

// TableName especifica el nombre de la tabla
func (Payment) TableName() string {
	return "payments"
}

// BeforeCreate genera el número de pago automáticamente
func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.PaymentNumber == "" {
		now := time.Now()
		dateStr := now.Format("20060102")

		var count int64
		tx.Model(&Payment{}).Where("payment_number LIKE ?", "PAY-"+dateStr+"%").Count(&count)

		p.PaymentNumber = fmt.Sprintf("PAY-%s-%06d", dateStr, count+1)
	}

	// Auto-aprobar pagos en efectivo
	if p.PaymentMethod == "efectivo" && p.Status == "" {
		p.Status = "aprobado"
		now := time.Now()
		p.ApprovedAt = &now
		p.ApprovedBy = &p.CreatedBy
	}

	return nil
}

// IsApproved verifica si el pago está aprobado
func (p *Payment) IsApproved() bool {
	return p.Status == "aprobado"
}

// IsCancelled verifica si el pago está cancelado
func (p *Payment) IsCancelled() bool {
	return p.Status == "anulado"
}

// Approve aprueba el pago
func (p *Payment) Approve(userID uint) {
	p.Status = "aprobado"
	now := time.Now()
	p.ApprovedAt = &now
	p.ApprovedBy = &userID
}

// Cancel cancela el pago
func (p *Payment) Cancel(userID uint, reason string) {
	p.Status = "anulado"
	now := time.Now()
	p.CancelledAt = &now
	p.CancelledBy = &userID
	p.CancellationReason = reason
}

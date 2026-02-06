package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Invoice representa una factura emitida
type Invoice struct {
	BaseModel
	InvoiceNumber      string     `gorm:"size:50;uniqueIndex;not null" json:"invoice_number"`
	OrderID            uint       `gorm:"not null" json:"order_id" binding:"required"`
	PatientID          uint       `gorm:"not null" json:"patient_id" binding:"required"`
	InvoiceDate        time.Time  `gorm:"type:date;not null;default:CURRENT_DATE" json:"invoice_date"`
	DueDate            *time.Time `gorm:"type:date" json:"due_date"`
	Subtotal           float64    `gorm:"type:decimal(10,2);not null" json:"subtotal" binding:"required,gte=0"`
	DiscountAmount     float64    `gorm:"type:decimal(10,2);default:0" json:"discount_amount"`
	TaxPercentage      float64    `gorm:"type:decimal(5,2);default:0" json:"tax_percentage"`
	TaxAmount          float64    `gorm:"type:decimal(10,2);default:0" json:"tax_amount"`
	TotalAmount        float64    `gorm:"type:decimal(10,2);not null" json:"total_amount" binding:"required,gt=0"`
	Status             string     `gorm:"size:20;default:'pendiente'" json:"status"`
	Notes              string     `gorm:"type:text" json:"notes"`
	CreatedBy          uint       `gorm:"not null" json:"created_by"`
	CancelledBy        *uint      `json:"cancelled_by"`
	CancelledAt        *time.Time `json:"cancelled_at"`
	CancellationReason string     `gorm:"type:text" json:"cancellation_reason"`

	// Relaciones
	Order           Order   `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Patient         Patient `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	Creator         User    `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	CancelledByUser *User   `gorm:"foreignKey:CancelledBy" json:"cancelled_by_user,omitempty"`
}

// TableName especifica el nombre de la tabla
func (Invoice) TableName() string {
	return "invoices"
}

// BeforeCreate genera el número de factura automáticamente
func (i *Invoice) BeforeCreate(tx *gorm.DB) error {
	if i.InvoiceNumber == "" {
		now := time.Now()
		dateStr := now.Format("20060102")

		var count int64
		tx.Model(&Invoice{}).Where("invoice_number LIKE ?", "INV-"+dateStr+"%").Count(&count)

		i.InvoiceNumber = fmt.Sprintf("INV-%s-%06d", dateStr, count+1)
	}

	// Calcular totales si no están definidos
	if i.TotalAmount == 0 {
		i.calculateTotals()
	}

	return nil
}

// BeforeUpdate recalcula totales antes de actualizar
func (i *Invoice) BeforeUpdate(tx *gorm.DB) error {
	i.calculateTotals()
	return nil
}

// calculateTotals calcula los totales de la factura
func (i *Invoice) calculateTotals() {
	afterDiscount := i.Subtotal - i.DiscountAmount
	i.TaxAmount = afterDiscount * (i.TaxPercentage / 100)
	i.TotalAmount = afterDiscount + i.TaxAmount
}

// IsPaid verifica si la factura está pagada
func (i *Invoice) IsPaid() bool {
	return i.Status == "pagada"
}

// IsOverdue verifica si la factura está vencida
func (i *Invoice) IsOverdue() bool {
	if i.DueDate == nil {
		return false
	}
	return time.Now().After(*i.DueDate) && i.Status != "pagada"
}

// MarkAsPaid marca la factura como pagada
func (i *Invoice) MarkAsPaid() {
	i.Status = "pagada"
}

// MarkAsOverdue marca la factura como vencida
func (i *Invoice) MarkAsOverdue() {
	if i.IsOverdue() {
		i.Status = "vencida"
	}
}

// Cancel cancela la factura
func (i *Invoice) Cancel(userID uint, reason string) {
	i.Status = "anulada"
	now := time.Now()
	i.CancelledAt = &now
	i.CancelledBy = &userID
	i.CancellationReason = reason
}

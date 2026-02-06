package models

import (
	"fmt"
	"time"
)

// Reagent representa un reactivo o insumo del laboratorio
type Reagent struct {
	BaseModel
	Name               string     `gorm:"size:200;not null" json:"name" binding:"required"`
	Code               string     `gorm:"size:50;uniqueIndex;not null" json:"code" binding:"required"`
	Description        string     `gorm:"type:text" json:"description"`
	Manufacturer       string     `gorm:"size:150" json:"manufacturer"`
	Supplier           string     `gorm:"size:150" json:"supplier"`
	LotNumber          string     `gorm:"size:100" json:"lot_number"`
	ExpirationDate     *time.Time `gorm:"type:date" json:"expiration_date"`
	QuantityAvailable  float64    `gorm:"type:decimal(10,2);not null;default:0" json:"quantity_available"`
	MinimumStock       float64    `gorm:"type:decimal(10,2);default:0" json:"minimum_stock"`
	UnitOfMeasure      string     `gorm:"size:50;not null" json:"unit_of_measure" binding:"required"`
	CostPerUnit        float64    `gorm:"type:decimal(10,2)" json:"cost_per_unit"`
	StorageLocation    string     `gorm:"size:100" json:"storage_location"`
	StorageTemperature string     `gorm:"size:50" json:"storage_temperature"`
	IsActive           bool       `gorm:"default:true" json:"is_active"`
}

// TableName especifica el nombre de la tabla
func (Reagent) TableName() string {
	return "reagents"
}

// IsLowStock verifica si el stock está bajo
func (r *Reagent) IsLowStock() bool {
	return r.QuantityAvailable <= r.MinimumStock
}

// IsExpired verifica si el reactivo está vencido
func (r *Reagent) IsExpired() bool {
	if r.ExpirationDate == nil {
		return false
	}
	return time.Now().After(*r.ExpirationDate)
}

// IsExpiringSoon verifica si el reactivo vence pronto (30 días)
func (r *Reagent) IsExpiringSoon() bool {
	if r.ExpirationDate == nil {
		return false
	}
	thirtyDaysFromNow := time.Now().AddDate(0, 0, 30)
	return r.ExpirationDate.Before(thirtyDaysFromNow) && r.ExpirationDate.After(time.Now())
}

// AddStock agrega stock al reactivo
func (r *Reagent) AddStock(quantity float64) {
	r.QuantityAvailable += quantity
}

// RemoveStock remueve stock del reactivo
func (r *Reagent) RemoveStock(quantity float64) error {
	if r.QuantityAvailable < quantity {
		return fmt.Errorf("stock insuficiente: disponible=%.2f, solicitado=%.2f", r.QuantityAvailable, quantity)
	}
	r.QuantityAvailable -= quantity
	return nil
}

// GetStockStatus retorna el estado del stock
func (r *Reagent) GetStockStatus() string {
	if r.QuantityAvailable == 0 {
		return "agotado"
	}
	if r.IsLowStock() {
		return "bajo"
	}
	return "suficiente"
}

// GetExpirationStatus retorna el estado de expiración
func (r *Reagent) GetExpirationStatus() string {
	if r.ExpirationDate == nil {
		return "sin_fecha"
	}
	if r.IsExpired() {
		return "vencido"
	}
	if r.IsExpiringSoon() {
		return "proximo_a_vencer"
	}
	return "vigente"
}

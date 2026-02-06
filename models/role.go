package models

import (
	"database/sql/driver"
	"encoding/json"
)

// Permissions representa los permisos en formato JSON
type Permissions map[string][]string

// Scan implementa la interfaz sql.Scanner para JSONB
func (p *Permissions) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, p)
}

// Value implementa la interfaz driver.Valuer para JSONB
func (p Permissions) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return json.Marshal(p)
}

// Role representa los roles del sistema
type Role struct {
	BaseModel
	Name        string      `gorm:"size:50;uniqueIndex;not null" json:"name" binding:"required"`
	Description string      `gorm:"type:text" json:"description"`
	Permissions Permissions `gorm:"type:jsonb;default:'{}'" json:"permissions"`
	IsActive    bool        `gorm:"default:true" json:"is_active"`

	// Relaciones
	Users []User `gorm:"foreignKey:RoleID" json:"-"`
}

// TableName especifica el nombre de la tabla
func (Role) TableName() string {
	return "roles"
}

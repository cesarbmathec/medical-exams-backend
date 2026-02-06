package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// JSONMap representa datos JSON flexibles
type JSONMap map[string]interface{}

// Scan implementa la interfaz sql.Scanner para JSONB
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Value implementa la interfaz driver.Valuer para JSONB
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// AuditLog representa el registro de auditoría del sistema
type AuditLog struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Table     string    `gorm:"column:table_name;size:100;not null;index:idx_audit_table_record" json:"table_name"`
	RecordID  uint      `gorm:"not null;index:idx_audit_table_record" json:"record_id"`
	Action    string    `gorm:"size:10;not null" json:"action"` // INSERT, UPDATE, DELETE
	OldValues JSONMap   `gorm:"type:jsonb" json:"old_values"`
	NewValues JSONMap   `gorm:"type:jsonb" json:"new_values"`
	UserID    *uint     `gorm:"index" json:"user_id"`
	IPAddress string    `gorm:"type:inet" json:"ip_address"`
	UserAgent string    `gorm:"type:text" json:"user_agent"`
	CreatedAt time.Time `gorm:"index:idx_audit_created" json:"created_at"`

	// Relaciones
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName especifica el nombre de la tabla
func (AuditLog) TableName() string {
	return "audit_logs"
}

// CreateAuditLog crea un nuevo registro de auditoría
func CreateAuditLog(table string, recordID uint, action string, oldValues, newValues map[string]interface{}, userID *uint, ipAddress, userAgent string) *AuditLog {
	return &AuditLog{
		Table:     table,
		RecordID:  recordID,
		Action:    action,
		OldValues: oldValues,
		NewValues: newValues,
		UserID:    userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
	}
}

// GetChangedFields retorna los campos que cambiaron
func (al *AuditLog) GetChangedFields() []string {
	if al.Action != "UPDATE" {
		return []string{}
	}

	changed := make([]string, 0)
	for key := range al.NewValues {
		if oldVal, exists := al.OldValues[key]; !exists || oldVal != al.NewValues[key] {
			changed = append(changed, key)
		}
	}
	return changed
}

// IsInsert verifica si es una inserción
func (al *AuditLog) IsInsert() bool {
	return al.Action == "INSERT"
}

// IsUpdate verifica si es una actualización
func (al *AuditLog) IsUpdate() bool {
	return al.Action == "UPDATE"
}

// IsDelete verifica si es una eliminación
func (al *AuditLog) IsDelete() bool {
	return al.Action == "DELETE"
}

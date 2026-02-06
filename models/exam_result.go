package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ExamResult representa el resultado de un parámetro de un examen
type ExamResult struct {
	BaseModel
	OrderExamID     uint       `gorm:"not null" json:"order_exam_id" binding:"required"`
	ExamParameterID uint       `gorm:"not null" json:"exam_parameter_id" binding:"required"`
	ValueNumeric    *float64   `gorm:"type:decimal(12,4)" json:"value_numeric"`
	ValueText       string     `gorm:"type:text" json:"value_text"`
	ValueBoolean    *bool      `json:"value_boolean"`
	IsAbnormal      bool       `gorm:"default:false" json:"is_abnormal"`
	IsCritical      bool       `gorm:"default:false" json:"is_critical"`
	AbnormalityType string     `gorm:"size:20" json:"abnormality_type"` // 'low', 'high', 'abnormal'
	TechnicianNotes string     `gorm:"type:text" json:"technician_notes"`
	Flags           string     `gorm:"size:50" json:"flags"` // 'L' (low), 'H' (high), 'C' (critical)
	Version         int        `gorm:"default:1" json:"version"`
	IsCurrent       bool       `gorm:"default:true" json:"is_current"`
	EnteredBy       uint       `gorm:"not null" json:"entered_by" binding:"required"`
	EnteredAt       time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"entered_at"`
	ValidatedBy     *uint      `json:"validated_by"`
	ValidatedAt     *time.Time `json:"validated_at"`

	// Relaciones
	OrderExam       OrderExam     `gorm:"foreignKey:OrderExamID" json:"order_exam,omitempty"`
	ExamParameter   ExamParameter `gorm:"foreignKey:ExamParameterID" json:"exam_parameter,omitempty"`
	EnteredByUser   User          `gorm:"foreignKey:EnteredBy" json:"entered_by_user,omitempty"`
	ValidatedByUser *User         `gorm:"foreignKey:ValidatedBy" json:"validated_by_user,omitempty"`
}

// TableName especifica el nombre de la tabla
func (ExamResult) TableName() string {
	return "exam_results"
}

// BeforeCreate hook para calcular si es anormal
func (er *ExamResult) BeforeCreate(tx *gorm.DB) error {
	// Cargar el parámetro del examen para obtener rangos de referencia
	var param ExamParameter
	if err := tx.First(&param, er.ExamParameterID).Error; err == nil {
		er.checkAbnormality(param)
	}
	return nil
}

// BeforeUpdate hook para calcular si es anormal
func (er *ExamResult) BeforeUpdate(tx *gorm.DB) error {
	// Cargar el parámetro del examen para obtener rangos de referencia
	var param ExamParameter
	if err := tx.First(&param, er.ExamParameterID).Error; err == nil {
		er.checkAbnormality(param)
	}
	return nil
}

// checkAbnormality verifica si el valor está fuera del rango normal
func (er *ExamResult) checkAbnormality(param ExamParameter) {
	if param.DataType == "numeric" && er.ValueNumeric != nil {
		if param.ReferenceMin != nil && *er.ValueNumeric < *param.ReferenceMin {
			er.IsAbnormal = true
			er.AbnormalityType = "low"
			er.Flags = "L"
		} else if param.ReferenceMax != nil && *er.ValueNumeric > *param.ReferenceMax {
			er.IsAbnormal = true
			er.AbnormalityType = "high"
			er.Flags = "H"
		} else {
			er.IsAbnormal = false
			er.AbnormalityType = ""
			er.Flags = ""
		}

		// Marcar como crítico si el parámetro lo requiere y está muy fuera del rango
		if param.IsCritical && er.IsAbnormal {
			er.IsCritical = true
			er.Flags += "C"
		}
	}
}

// IsValidated verifica si el resultado ha sido validado
func (er *ExamResult) IsValidated() bool {
	return er.ValidatedAt != nil && er.ValidatedBy != nil
}

// GetDisplayValue retorna el valor en el formato apropiado
func (er *ExamResult) GetDisplayValue() string {
	if er.ValueNumeric != nil {
		return fmt.Sprintf("%.2f", *er.ValueNumeric)
	}
	if er.ValueBoolean != nil {
		if *er.ValueBoolean {
			return "Positivo"
		}
		return "Negativo"
	}
	return er.ValueText
}

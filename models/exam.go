package models

import (
	"database/sql/driver"
	"encoding/json"
)

// ExamCategory representa una categoría de exámenes
type ExamCategory struct {
	BaseModel
	Name         string `gorm:"size:100;uniqueIndex;not null" json:"name" binding:"required"`
	Description  string `gorm:"type:text" json:"description"`
	Code         string `gorm:"size:20;uniqueIndex;not null" json:"code" binding:"required"`
	DisplayOrder int    `gorm:"default:0" json:"display_order"`
	IsActive     bool   `gorm:"default:true" json:"is_active"`

	// Relaciones
	ExamTypes []ExamType `gorm:"foreignKey:CategoryID" json:"exam_types,omitempty"`
}

func (ExamCategory) TableName() string {
	return "exam_categories"
}

// SampleType representa un tipo de muestra
type SampleType struct {
	BaseModel
	Name                   string `gorm:"size:100;uniqueIndex;not null" json:"name" binding:"required"`
	Description            string `gorm:"type:text" json:"description"`
	CollectionInstructions string `gorm:"type:text" json:"collection_instructions"`
	StorageRequirements    string `gorm:"type:text" json:"storage_requirements"`
	StorageTemperature     string `gorm:"size:50" json:"storage_temperature"`
	MaxStorageTimeHours    int    `json:"max_storage_time_hours"`
	IsActive               bool   `gorm:"default:true" json:"is_active"`

	// Relaciones
	ExamTypes []ExamType `gorm:"foreignKey:SampleTypeID" json:"exam_types,omitempty"`
}

func (SampleType) TableName() string {
	return "sample_types"
}

// ExamType representa un tipo de examen
type ExamType struct {
	BaseModel
	Code                    string  `gorm:"size:50;uniqueIndex;not null" json:"code" binding:"required"`
	Name                    string  `gorm:"size:200;not null" json:"name" binding:"required"`
	Description             string  `gorm:"type:text" json:"description"`
	CategoryID              uint    `gorm:"not null" json:"category_id" binding:"required"`
	SampleTypeID            uint    `gorm:"not null" json:"sample_type_id" binding:"required"`
	BasePrice               float64 `gorm:"type:decimal(10,2);not null" json:"base_price" binding:"required,gt=0"`
	PreparationInstructions string  `gorm:"type:text" json:"preparation_instructions"`
	ProcessingTimeHours     int     `gorm:"default:24" json:"processing_time_hours"`
	RequiresFasting         bool    `gorm:"default:false" json:"requires_fasting"`
	FastingHours            int     `json:"fasting_hours"`
	RequiresAppointment     bool    `gorm:"default:false" json:"requires_appointment"`
	IsActive                bool    `gorm:"default:true" json:"is_active"`

	// Relaciones
	Category   ExamCategory    `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	SampleType SampleType      `gorm:"foreignKey:SampleTypeID" json:"sample_type,omitempty"`
	Parameters []ExamParameter `gorm:"foreignKey:ExamTypeID" json:"parameters,omitempty"`
}

func (ExamType) TableName() string {
	return "exam_types"
}

// SelectOptions para parámetros de tipo select
type SelectOptions []string

func (s *SelectOptions) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

func (s SelectOptions) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

// ExamParameter representa un parámetro de un examen
type ExamParameter struct {
	BaseModel
	ExamTypeID         uint          `gorm:"not null" json:"exam_type_id"`
	ParameterName      string        `gorm:"size:200;not null" json:"parameter_name" binding:"required"`
	ParameterCode      string        `gorm:"size:50" json:"parameter_code"`
	UnitOfMeasure      string        `gorm:"size:50" json:"unit_of_measure"`
	ReferenceMin       *float64      `gorm:"type:decimal(12,4)" json:"reference_min"`
	ReferenceMax       *float64      `gorm:"type:decimal(12,4)" json:"reference_max"`
	ReferenceValueText string        `gorm:"type:text" json:"reference_value_text"`
	DataType           string        `gorm:"size:20;not null" json:"data_type" binding:"required,oneof=numeric text boolean select"`
	SelectOptions      SelectOptions `gorm:"type:jsonb" json:"select_options"`
	DisplayOrder       int           `gorm:"default:0" json:"display_order"`
	IsCritical         bool          `gorm:"default:false" json:"is_critical"`
	IsRequired         bool          `gorm:"default:true" json:"is_required"`
	ValidationRules    Permissions   `gorm:"type:jsonb" json:"validation_rules"` // Reutilizamos el tipo
	Notes              string        `gorm:"type:text" json:"notes"`

	// Relaciones
	ExamType ExamType `gorm:"foreignKey:ExamTypeID" json:"exam_type,omitempty"`
}

func (ExamParameter) TableName() string {
	return "exam_parameters"
}

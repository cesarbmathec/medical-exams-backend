package dtos

// Para crear una orden con múltiples exámenes a la vez
type CreateOrderRequest struct {
	PatientID       uint               `json:"patient_id" binding:"required"`
	Priority        string             `json:"priority" binding:"required,oneof=normal urgente stat"`
	ReferringDoctor string             `json:"referring_doctor"`
	Diagnosis       string             `json:"diagnosis"`
	Exams           []OrderExamRequest `json:"exams" binding:"required,gt=0"`
}

type OrderExamRequest struct {
	ExamTypeID uint    `json:"exam_type_id" binding:"required"`
	Price      float64 `json:"price" binding:"required,gt=0"`
}

// Para registrar resultados de un examen
type UpdateResultRequest struct {
	ParameterID  uint     `json:"exam_parameter_id" binding:"required"`
	ValueNumeric *float64 `json:"value_numeric"`
	ValueText    string   `json:"value_text"`
}

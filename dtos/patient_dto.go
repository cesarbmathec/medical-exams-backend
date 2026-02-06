package dtos

import "time"

// Definimos lo que esperamos recibir exactamente del cliente
type CreatePatientRequest struct {
	DocumentType   string    `json:"document_type" binding:"required,oneof=cedula pasaporte rif otro"`
	DocumentNumber string    `json:"document_number" binding:"required"`
	FirstName      string    `json:"first_name" binding:"required"`
	LastName       string    `json:"last_name" binding:"required"`
	DateOfBirth    time.Time `json:"date_of_birth" binding:"required"`
	Gender         string    `json:"gender" binding:"omitempty,oneof=M F O"`
	Phone          string    `json:"phone"`
	Email          string    `json:"email" binding:"omitempty,email"`
	BloodType      string    `json:"blood_type" binding:"omitempty,oneof=A+ A- B+ B- AB+ AB- O+ O-"`
}

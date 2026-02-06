package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User representa un usuario del sistema
type User struct {
	BaseModel
	Username            string     `gorm:"size:50;uniqueIndex;not null" json:"username" binding:"required,min=3,max=50"`
	Email               string     `gorm:"size:100;uniqueIndex;not null" json:"email" binding:"required,email"`
	PasswordHash        string     `gorm:"size:255;not null" json:"-"`
	Password            string     `gorm:"-" json:"password,omitempty" binding:"required,min=6"`
	FullName            string     `gorm:"size:150;not null" json:"full_name" binding:"required"`
	Phone               string     `gorm:"size:20" json:"phone"`
	RoleID              uint       `gorm:"not null" json:"role_id" binding:"required"`
	IsActive            bool       `gorm:"default:true" json:"is_active"`
	LastLogin           *time.Time `json:"last_login,omitempty"`
	FailedLoginAttempts int        `gorm:"default:0" json:"-"`
	LockedUntil         *time.Time `json:"-"`

	// Relaciones
	Role Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// TableName especifica el nombre de la tabla
func (User) TableName() string {
	return "users"
}

// BeforeCreate hook para hashear la contraseña antes de crear
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.PasswordHash = string(hashedPassword)
		u.Password = "" // Limpiar el password en texto plano
	}
	return nil
}

// BeforeUpdate hook para hashear la contraseña si se actualiza
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.PasswordHash = string(hashedPassword)
		u.Password = "" // Limpiar el password en texto plano
	}
	return nil
}

// CheckPassword verifica si la contraseña es correcta
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// UserResponse es la estructura para respuestas sin datos sensibles
type UserResponse struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	FullName  string     `json:"full_name"`
	Phone     string     `json:"phone"`
	RoleID    uint       `json:"role_id"`
	RoleName  string     `json:"role_name"`
	IsActive  bool       `json:"is_active"`
	LastLogin *time.Time `json:"last_login"`
	CreatedAt time.Time  `json:"created_at"`
}

// ToResponse convierte User a UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FullName:  u.FullName,
		Phone:     u.Phone,
		RoleID:    u.RoleID,
		RoleName:  u.Role.Name,
		IsActive:  u.IsActive,
		LastLogin: u.LastLogin,
		CreatedAt: u.CreatedAt,
	}
}

// domain/user.go
package domain

import "gorm.io/gorm"

type Role uint32

const (
	RoleUser  Role = 2001
	RoleAdmin Role = 5150
)

type User struct {
	gorm.Model

	Email    string `json:"email" gorm:"size:255;unique;not null"`
	Password string `json:"password" gorm:"size:255;not null"`
	Name     string `json:"name" gorm:"size:255;not null"`
	Role     Role   `json:"role" gorm:"not null;default:2001"`

	Jobs []Job `json:"jobs,omitempty" gorm:"foreignKey:UserID"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required,min=3,max=25"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  Role   `json:"role"`
}

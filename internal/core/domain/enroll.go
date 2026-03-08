package domain

import "gorm.io/gorm"

type EnrollmentStatus string

const (
	StatusPending  EnrollmentStatus = "pending"
	StatusAccepted EnrollmentStatus = "accepted"
	StatusRejected EnrollmentStatus = "rejected"
)

type Enrollment struct {
	gorm.Model

	UserID uint             `json:"user_id" gorm:"not null"`
	JobID  uint             `json:"job_id"  gorm:"not null"`
	Status EnrollmentStatus `json:"status"  gorm:"size:50;not null;default:'pending'"`

	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Job  Job  `json:"job,omitempty"  gorm:"foreignKey:JobID"`
}

type EnrollJobRequest struct {
	JobID uint `json:"job_id" binding:"required"`
}

type UpdateEnrollmentStatusRequest struct {
	Status EnrollmentStatus `json:"status" binding:"required,oneof=pending accepted rejected"`
}

type EnrollmentResponse struct {
	ID        uint             `json:"id"`
	Status    EnrollmentStatus `json:"status"`
	Job       JobResponse      `json:"job"`
	AppliedBy UserResponse     `json:"applied_by"`
}

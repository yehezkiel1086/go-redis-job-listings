package domain

import "gorm.io/gorm"

type JobType string
type ExperienceLevel string

const (
	JobTypeFullTime   JobType = "full_time"
	JobTypePartTime   JobType = "part_time"
	JobTypeContract   JobType = "contract"
	JobTypeInternship JobType = "internship"
	JobTypeRemote     JobType = "remote"

	ExperienceLevelJunior ExperienceLevel = "junior"
	ExperienceLevelMid    ExperienceLevel = "mid"
	ExperienceLevelSenior ExperienceLevel = "senior"
	ExperienceLevelLead   ExperienceLevel = "lead"
)

type Job struct {
	gorm.Model

	UserID          uint            `json:"user_id"          gorm:"not null"`
	Title           string          `json:"title"            gorm:"size:255;not null"`
	Company         string          `json:"company"          gorm:"size:255;not null"`
	Location        string          `json:"location"         gorm:"size:255;not null"`
	Type            JobType         `json:"type"             gorm:"size:50;not null"`
	ExperienceLevel ExperienceLevel `json:"experience_level" gorm:"size:50;not null"`
	Description     string          `json:"description"      gorm:"type:text;not null"`
	Requirements    string          `json:"requirements"     gorm:"type:text"`
	SalaryMin       float64         `json:"salary_min"       gorm:"default:0"`
	SalaryMax       float64         `json:"salary_max"       gorm:"default:0"`
	IsActive        bool            `json:"is_active"        gorm:"not null;default:true"`

	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type CreateJobRequest struct {
	Title           string          `json:"title"            binding:"required,min=3,max=100"`
	Company         string          `json:"company"          binding:"required,min=2,max=100"`
	Location        string          `json:"location"         binding:"required"`
	Type            JobType         `json:"type"             binding:"required,oneof=full_time part_time contract internship remote"`
	ExperienceLevel ExperienceLevel `json:"experience_level" binding:"required,oneof=junior mid senior lead"`
	Description     string          `json:"description"      binding:"required,min=10"`
	Requirements    string          `json:"requirements"`
	SalaryMin       float64         `json:"salary_min"       binding:"omitempty,min=0"`
	SalaryMax       float64         `json:"salary_max"       binding:"omitempty,min=0"`
}

type UpdateJobRequest struct {
	Title           string          `json:"title"            binding:"omitempty,min=3,max=100"`
	Company         string          `json:"company"          binding:"omitempty,min=2,max=100"`
	Location        string          `json:"location"         binding:"omitempty"`
	Type            JobType         `json:"type"             binding:"omitempty,oneof=full_time part_time contract internship remote"`
	ExperienceLevel ExperienceLevel `json:"experience_level" binding:"omitempty,oneof=junior mid senior lead"`
	Description     string          `json:"description"      binding:"omitempty,min=10"`
	Requirements    string          `json:"requirements"     binding:"omitempty"`
	SalaryMin       float64         `json:"salary_min"       binding:"omitempty,min=0"`
	SalaryMax       float64         `json:"salary_max"       binding:"omitempty,min=0"`
	IsActive        *bool           `json:"is_active"        binding:"omitempty"`
}

type JobFilter struct {
	Type            JobType         `form:"type"             binding:"omitempty,oneof=full_time part_time contract internship remote"`
	ExperienceLevel ExperienceLevel `form:"experience_level" binding:"omitempty,oneof=junior mid senior lead"`
	Location        string          `form:"location"         binding:"omitempty"`
	Search          string          `form:"search"           binding:"omitempty"`
	IsActive        *bool           `form:"is_active"        binding:"omitempty"`
}

type JobResponse struct {
	ID              uint            `json:"id"`
	Title           string          `json:"title"`
	Company         string          `json:"company"`
	Location        string          `json:"location"`
	Type            JobType         `json:"type"`
	ExperienceLevel ExperienceLevel `json:"experience_level"`
	Description     string          `json:"description"`
	Requirements    string          `json:"requirements"`
	SalaryMin       float64         `json:"salary_min"`
	SalaryMax       float64         `json:"salary_max"`
	IsActive        bool            `json:"is_active"`
	PostedBy        UserResponse    `json:"posted_by"`
}

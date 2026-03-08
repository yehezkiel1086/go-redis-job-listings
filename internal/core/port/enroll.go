package port

import (
	"context"

	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/domain"
)

type EnrollmentRepository interface {
	CreateEnrollment(ctx context.Context, enrollment *domain.Enrollment) (*domain.Enrollment, error)
	GetEnrollmentByID(ctx context.Context, id uint) (*domain.Enrollment, error)
	GetEnrollmentByUserAndJob(ctx context.Context, userID, jobID uint) (*domain.Enrollment, error)
	GetEnrollmentsByUserID(ctx context.Context, userID uint) ([]domain.Enrollment, error)
	GetEnrollmentsByJobID(ctx context.Context, jobID uint) ([]domain.Enrollment, error)
	UpdateEnrollment(ctx context.Context, enrollment *domain.Enrollment) (*domain.Enrollment, error)
	DeleteEnrollmentByID(ctx context.Context, id uint) error
}

type EnrollmentService interface {
	EnrollJob(ctx context.Context, userID uint, jobID uint) (*domain.EnrollmentResponse, error)
	GetEnrollmentByID(ctx context.Context, id uint) (*domain.EnrollmentResponse, error)
	GetMyEnrollments(ctx context.Context, userID uint) ([]domain.EnrollmentResponse, error)
	GetJobEnrollments(ctx context.Context, jobID uint, requesterID uint, role domain.Role) ([]domain.EnrollmentResponse, error)
	UpdateEnrollmentStatus(ctx context.Context, id uint, requesterID uint, role domain.Role, req *domain.UpdateEnrollmentStatusRequest) (*domain.EnrollmentResponse, error)
	DeleteEnrollmentByID(ctx context.Context, id uint, requesterID uint, role domain.Role) error
}

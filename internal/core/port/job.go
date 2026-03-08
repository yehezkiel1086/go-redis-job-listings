package port

import (
	"context"

	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/domain"
)

type JobRepository interface {
	CreateJob(ctx context.Context, job *domain.Job) (*domain.Job, error)
	GetJobByID(ctx context.Context, id uint) (*domain.Job, error)
	GetAllJobs(ctx context.Context, filter domain.JobFilter) ([]domain.Job, error)
	GetJobsByUserID(ctx context.Context, userID uint) ([]domain.Job, error)
	UpdateJob(ctx context.Context, job *domain.Job) (*domain.Job, error)
	DeleteJobByID(ctx context.Context, id uint) error
}

type JobService interface {
	CreateJob(ctx context.Context, userID uint, req *domain.CreateJobRequest) (*domain.JobResponse, error)
	GetJobByID(ctx context.Context, id uint) (*domain.JobResponse, error)
	GetAllJobs(ctx context.Context, filter domain.JobFilter) ([]domain.JobResponse, error)
	GetJobsByUserID(ctx context.Context, userID uint) ([]domain.JobResponse, error)
	UpdateJob(ctx context.Context, id uint, userID uint, role domain.Role, req *domain.UpdateJobRequest) (*domain.JobResponse, error)
	DeleteJobByID(ctx context.Context, id uint, userID uint, role domain.Role) error
}

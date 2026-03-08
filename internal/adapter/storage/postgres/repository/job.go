package repository

import (
	"context"

	"github.com/yehezkiel1086/go-redis-job-listings/internal/adapter/storage/postgres"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/domain"
)

type JobRepository struct {
	db *postgres.DB
}

func NewJobRepository(db *postgres.DB) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) CreateJob(ctx context.Context, job *domain.Job) (*domain.Job, error) {
	db := r.db.GetDB()

	if err := db.WithContext(ctx).Create(job).Error; err != nil {
		return nil, err
	}

	return r.GetJobByID(ctx, job.ID)
}

func (r *JobRepository) GetJobByID(ctx context.Context, id uint) (*domain.Job, error) {
	db := r.db.GetDB()

	var job domain.Job
	if err := db.WithContext(ctx).Preload("User").First(&job, id).Error; err != nil {
		return nil, err
	}

	return &job, nil
}

func (r *JobRepository) GetAllJobs(ctx context.Context, filter domain.JobFilter) ([]domain.Job, error) {
	db := r.db.GetDB()

	query := db.WithContext(ctx).Preload("User")

	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.ExperienceLevel != "" {
		query = query.Where("experience_level = ?", filter.ExperienceLevel)
	}
	if filter.Location != "" {
		query = query.Where("location ILIKE ?", "%"+filter.Location+"%")
	}
	if filter.Search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ? OR company ILIKE ?",
			"%"+filter.Search+"%",
			"%"+filter.Search+"%",
			"%"+filter.Search+"%",
		)
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	var jobs []domain.Job
	if err := query.Find(&jobs).Error; err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r *JobRepository) GetJobsByUserID(ctx context.Context, userID uint) ([]domain.Job, error) {
	db := r.db.GetDB()

	var jobs []domain.Job
	if err := db.WithContext(ctx).Preload("User").Where("user_id = ?", userID).Find(&jobs).Error; err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r *JobRepository) UpdateJob(ctx context.Context, job *domain.Job) (*domain.Job, error) {
	db := r.db.GetDB()

	if err := db.WithContext(ctx).Save(job).Error; err != nil {
		return nil, err
	}

	return r.GetJobByID(ctx, job.ID)
}

func (r *JobRepository) DeleteJobByID(ctx context.Context, id uint) error {
	db := r.db.GetDB()

	return db.WithContext(ctx).Delete(&domain.Job{}, id).Error
}

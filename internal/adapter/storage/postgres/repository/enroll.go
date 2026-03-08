package repository

import (
	"context"

	"github.com/yehezkiel1086/go-redis-job-listings/internal/adapter/storage/postgres"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/domain"
)

type EnrollmentRepository struct {
	db *postgres.DB
}

func NewEnrollmentRepository(db *postgres.DB) *EnrollmentRepository {
	return &EnrollmentRepository{db: db}
}

func (r *EnrollmentRepository) CreateEnrollment(ctx context.Context, enrollment *domain.Enrollment) (*domain.Enrollment, error) {
	if err := r.db.GetDB().WithContext(ctx).Create(enrollment).Error; err != nil {
		return nil, err
	}

	return r.GetEnrollmentByID(ctx, enrollment.ID)
}

func (r *EnrollmentRepository) GetEnrollmentByID(ctx context.Context, id uint) (*domain.Enrollment, error) {
	var enrollment domain.Enrollment
	if err := r.db.GetDB().WithContext(ctx).
		Preload("User").
		Preload("Job").
		Preload("Job.User").
		First(&enrollment, id).Error; err != nil {
		return nil, err
	}

	return &enrollment, nil
}

func (r *EnrollmentRepository) GetEnrollmentByUserAndJob(ctx context.Context, userID, jobID uint) (*domain.Enrollment, error) {
	var enrollment domain.Enrollment
	if err := r.db.GetDB().WithContext(ctx).
		Where("user_id = ? AND job_id = ?", userID, jobID).
		First(&enrollment).Error; err != nil {
		return nil, err
	}

	return &enrollment, nil
}

func (r *EnrollmentRepository) GetEnrollmentsByUserID(ctx context.Context, userID uint) ([]domain.Enrollment, error) {
	var enrollments []domain.Enrollment
	if err := r.db.GetDB().WithContext(ctx).
		Preload("User").
		Preload("Job").
		Preload("Job.User").
		Where("user_id = ?", userID).
		Find(&enrollments).Error; err != nil {
		return nil, err
	}

	return enrollments, nil
}

func (r *EnrollmentRepository) GetEnrollmentsByJobID(ctx context.Context, jobID uint) ([]domain.Enrollment, error) {
	var enrollments []domain.Enrollment
	if err := r.db.GetDB().WithContext(ctx).
		Preload("User").
		Preload("Job").
		Preload("Job.User").
		Where("job_id = ?", jobID).
		Find(&enrollments).Error; err != nil {
		return nil, err
	}

	return enrollments, nil
}

func (r *EnrollmentRepository) UpdateEnrollment(ctx context.Context, enrollment *domain.Enrollment) (*domain.Enrollment, error) {
	if err := r.db.GetDB().WithContext(ctx).Save(enrollment).Error; err != nil {
		return nil, err
	}

	return r.GetEnrollmentByID(ctx, enrollment.ID)
}

func (r *EnrollmentRepository) DeleteEnrollmentByID(ctx context.Context, id uint) error {
	return r.db.GetDB().WithContext(ctx).Delete(&domain.Enrollment{}, id).Error
}

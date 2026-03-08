package service

import (
	"context"
	"time"

	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/domain"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/port"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/util"
)

const (
	enrollmentPrefix     = "enrollment"
	myEnrollmentsPrefix  = "enrollment:user"
	jobEnrollmentsPrefix = "enrollment:job"
	enrollmentCacheTTL   = 5 * time.Minute
)

type EnrollmentService struct {
	repo    port.EnrollmentRepository
	jobRepo port.JobRepository
	cache   port.CacheRepository
}

func NewEnrollmentService(repo port.EnrollmentRepository, jobRepo port.JobRepository, cache port.CacheRepository) *EnrollmentService {
	return &EnrollmentService{repo: repo, jobRepo: jobRepo, cache: cache}
}

func (s *EnrollmentService) EnrollJob(ctx context.Context, userID uint, jobID uint) (*domain.EnrollmentResponse, error) {
	// verify job exists and is active.
	job, err := s.jobRepo.GetJobByID(ctx, jobID)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	if !job.IsActive {
		return nil, domain.ErrInactiveJob
	}

	// guard against duplicate enrollment.
	existing, err := s.repo.GetEnrollmentByUserAndJob(ctx, userID, jobID)
	if err == nil && existing != nil {
		return nil, domain.ErrDuplicateEnroll
	}

	enrollment := &domain.Enrollment{
		UserID: userID,
		JobID:  jobID,
		Status: domain.StatusPending,
	}

	created, err := s.repo.CreateEnrollment(ctx, enrollment)
	if err != nil {
		return nil, err
	}

	// invalidate user's enrollment list and job's enrollment list.
	_ = s.cache.Delete(ctx, util.GenerateCacheKey(myEnrollmentsPrefix, userID))
	_ = s.cache.Delete(ctx, util.GenerateCacheKey(jobEnrollmentsPrefix, jobID))

	return toEnrollmentResponse(created), nil
}

func (s *EnrollmentService) GetEnrollmentByID(ctx context.Context, id uint) (*domain.EnrollmentResponse, error) {
	key := util.GenerateCacheKey(enrollmentPrefix, id)

	if cached, err := s.cache.Get(ctx, key); err == nil {
		var resp domain.EnrollmentResponse
		if err := util.Deserialize([]byte(cached), &resp); err == nil {
			return &resp, nil
		}
	}

	enrollment, err := s.repo.GetEnrollmentByID(ctx, id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	resp := toEnrollmentResponse(enrollment)

	if value, err := util.Serialize(resp); err == nil {
		_ = s.cache.Set(ctx, key, value, enrollmentCacheTTL)
	}

	return resp, nil
}

func (s *EnrollmentService) GetMyEnrollments(ctx context.Context, userID uint) ([]domain.EnrollmentResponse, error) {
	key := util.GenerateCacheKey(myEnrollmentsPrefix, userID)

	if cached, err := s.cache.Get(ctx, key); err == nil {
		var responses []domain.EnrollmentResponse
		if err := util.Deserialize([]byte(cached), &responses); err == nil {
			return responses, nil
		}
	}

	enrollments, err := s.repo.GetEnrollmentsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]domain.EnrollmentResponse, 0, len(enrollments))
	for _, e := range enrollments {
		responses = append(responses, *toEnrollmentResponse(&e))
	}

	if value, err := util.Serialize(responses); err == nil {
		_ = s.cache.Set(ctx, key, value, enrollmentCacheTTL)
	}

	return responses, nil
}

func (s *EnrollmentService) GetJobEnrollments(ctx context.Context, jobID uint, requesterID uint, role domain.Role) ([]domain.EnrollmentResponse, error) {
	// only admin or the job owner can see applicants.
	if role != domain.RoleAdmin {
		job, err := s.jobRepo.GetJobByID(ctx, jobID)
		if err != nil {
			return nil, domain.ErrNotFound
		}
		if job.UserID != requesterID {
			return nil, domain.ErrForbidden
		}
	}

	key := util.GenerateCacheKey(jobEnrollmentsPrefix, jobID)

	if cached, err := s.cache.Get(ctx, key); err == nil {
		var responses []domain.EnrollmentResponse
		if err := util.Deserialize([]byte(cached), &responses); err == nil {
			return responses, nil
		}
	}

	enrollments, err := s.repo.GetEnrollmentsByJobID(ctx, jobID)
	if err != nil {
		return nil, err
	}

	responses := make([]domain.EnrollmentResponse, 0, len(enrollments))
	for _, e := range enrollments {
		responses = append(responses, *toEnrollmentResponse(&e))
	}

	if value, err := util.Serialize(responses); err == nil {
		_ = s.cache.Set(ctx, key, value, enrollmentCacheTTL)
	}

	return responses, nil
}

func (s *EnrollmentService) UpdateEnrollmentStatus(ctx context.Context, id uint, requesterID uint, role domain.Role, req *domain.UpdateEnrollmentStatusRequest) (*domain.EnrollmentResponse, error) {
	enrollment, err := s.repo.GetEnrollmentByID(ctx, id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	// only admin or the job owner can accept/reject.
	if role != domain.RoleAdmin {
		job, err := s.jobRepo.GetJobByID(ctx, enrollment.JobID)
		if err != nil {
			return nil, domain.ErrNotFound
		}
		if job.UserID != requesterID {
			return nil, domain.ErrForbidden
		}
	}

	enrollment.Status = req.Status

	updated, err := s.repo.UpdateEnrollment(ctx, enrollment)
	if err != nil {
		return nil, err
	}

	// invalidate individual, user list, and job list caches.
	_ = s.cache.Delete(ctx,
		util.GenerateCacheKey(enrollmentPrefix, id),
		util.GenerateCacheKey(myEnrollmentsPrefix, enrollment.UserID),
		util.GenerateCacheKey(jobEnrollmentsPrefix, enrollment.JobID),
	)

	return toEnrollmentResponse(updated), nil
}

func (s *EnrollmentService) DeleteEnrollmentByID(ctx context.Context, id uint, requesterID uint, role domain.Role) error {
	enrollment, err := s.repo.GetEnrollmentByID(ctx, id)
	if err != nil {
		return domain.ErrNotFound
	}

	// only the applicant themselves or an admin can withdraw.
	if role != domain.RoleAdmin && enrollment.UserID != requesterID {
		return domain.ErrForbidden
	}

	if err := s.repo.DeleteEnrollmentByID(ctx, id); err != nil {
		return err
	}

	// invalidate individual, user list, and job list caches.
	_ = s.cache.Delete(ctx,
		util.GenerateCacheKey(enrollmentPrefix, id),
		util.GenerateCacheKey(myEnrollmentsPrefix, enrollment.UserID),
		util.GenerateCacheKey(jobEnrollmentsPrefix, enrollment.JobID),
	)

	return nil
}

func toEnrollmentResponse(e *domain.Enrollment) *domain.EnrollmentResponse {
	return &domain.EnrollmentResponse{
		ID:     e.ID,
		Status: e.Status,
		Job:    *toJobResponse(&e.Job),
		AppliedBy: domain.UserResponse{
			ID:    e.User.ID,
			Email: e.User.Email,
			Name:  e.User.Name,
			Role:  e.User.Role,
		},
	}
}

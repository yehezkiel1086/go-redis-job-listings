package service

import (
	"context"
	"time"

	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/domain"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/port"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/util"
)

const (
	jobPrefix      = "job"
	allJobsPrefix  = "job:all"
	userJobsPrefix = "job:user"
	jobCacheTTL    = 5 * time.Minute
)

type JobService struct {
	repo  port.JobRepository
	cache port.CacheRepository
}

func NewJobService(repo port.JobRepository, cache port.CacheRepository) *JobService {
	return &JobService{repo: repo, cache: cache}
}

func (s *JobService) CreateJob(ctx context.Context, userID uint, req *domain.CreateJobRequest) (*domain.JobResponse, error) {
	job := &domain.Job{
		UserID:          userID,
		Title:           req.Title,
		Company:         req.Company,
		Location:        req.Location,
		Type:            req.Type,
		ExperienceLevel: req.ExperienceLevel,
		Description:     req.Description,
		Requirements:    req.Requirements,
		SalaryMin:       req.SalaryMin,
		SalaryMax:       req.SalaryMax,
		IsActive:        true,
	}

	created, err := s.repo.CreateJob(ctx, job)
	if err != nil {
		return nil, err
	}

	// invalidate - new job affects all-jobs and this user's jobs
	_ = s.cache.DeleteByPrefix(ctx, allJobsPrefix)
	_ = s.cache.Delete(ctx, util.GenerateCacheKey(userJobsPrefix, userID))

	return toJobResponse(created), nil
}

func (s *JobService) GetJobByID(ctx context.Context, id uint) (*domain.JobResponse, error) {
	key := util.GenerateCacheKey(jobPrefix, id)

	if cached, err := s.cache.Get(ctx, key); err == nil {
		var resp domain.JobResponse
		if err := util.Deserialize([]byte(cached), &resp); err == nil {
			return &resp, nil
		}
	}

	job, err := s.repo.GetJobByID(ctx, id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	resp := toJobResponse(job)

	if value, err := util.Serialize(resp); err == nil {
		_ = s.cache.Set(ctx, key, value, jobCacheTTL)
	}

	return resp, nil
}

func (s *JobService) GetAllJobs(ctx context.Context, filter domain.JobFilter) ([]domain.JobResponse, error) {
	// cache key encodes filter params so different filters are cached separately.
	// e.g: job:all:type:full-time:experience:mid-level:location:remote:search:go:is_active:true
	key := util.GenerateCacheKey(allJobsPrefix, util.GenerateCacheKeyParams(
		string(filter.Type),
		string(filter.ExperienceLevel),
		filter.Location,
		filter.Search,
		filter.IsActive,
	))

	if cached, err := s.cache.Get(ctx, key); err == nil {
		var responses []domain.JobResponse
		if err := util.Deserialize([]byte(cached), &responses); err == nil {
			return responses, nil
		}
	}

	jobs, err := s.repo.GetAllJobs(ctx, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]domain.JobResponse, 0, len(jobs))
	for _, j := range jobs {
		responses = append(responses, *toJobResponse(&j))
	}

	if value, err := util.Serialize(responses); err == nil {
		_ = s.cache.Set(ctx, key, value, jobCacheTTL)
	}

	return responses, nil
}

func (s *JobService) GetJobsByUserID(ctx context.Context, userID uint) ([]domain.JobResponse, error) {
	key := util.GenerateCacheKey(userJobsPrefix, userID)

	if cached, err := s.cache.Get(ctx, key); err == nil {
		var responses []domain.JobResponse
		if err := util.Deserialize([]byte(cached), &responses); err == nil {
			return responses, nil
		}
	}

	jobs, err := s.repo.GetJobsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]domain.JobResponse, 0, len(jobs))
	for _, j := range jobs {
		responses = append(responses, *toJobResponse(&j))
	}

	if value, err := util.Serialize(responses); err == nil {
		_ = s.cache.Set(ctx, key, value, jobCacheTTL)
	}

	return responses, nil
}

func (s *JobService) UpdateJob(ctx context.Context, id uint, userID uint, role domain.Role, req *domain.UpdateJobRequest) (*domain.JobResponse, error) {
	job, err := s.repo.GetJobByID(ctx, id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	// only the owner or an admin can update.
	if role != domain.RoleAdmin && job.UserID != userID {
		return nil, domain.ErrForbidden
	}

	if req.Title != "" {
		job.Title = req.Title
	}
	if req.Company != "" {
		job.Company = req.Company
	}
	if req.Location != "" {
		job.Location = req.Location
	}
	if req.Type != "" {
		job.Type = req.Type
	}
	if req.ExperienceLevel != "" {
		job.ExperienceLevel = req.ExperienceLevel
	}
	if req.Description != "" {
		job.Description = req.Description
	}
	if req.Requirements != "" {
		job.Requirements = req.Requirements
	}
	if req.SalaryMin != 0 {
		job.SalaryMin = req.SalaryMin
	}
	if req.SalaryMax != 0 {
		job.SalaryMax = req.SalaryMax
	}
	if req.IsActive != nil {
		job.IsActive = *req.IsActive
	}

	updated, err := s.repo.UpdateJob(ctx, job)
	if err != nil {
		return nil, err
	}

	// invalidate individual job, all-jobs list, and owner's job list.
	_ = s.cache.Delete(ctx, util.GenerateCacheKey(jobPrefix, id))
	_ = s.cache.DeleteByPrefix(ctx, allJobsPrefix)
	_ = s.cache.Delete(ctx, util.GenerateCacheKey(userJobsPrefix, job.UserID))

	return toJobResponse(updated), nil
}

func (s *JobService) DeleteJobByID(ctx context.Context, id uint, userID uint, role domain.Role) error {
	job, err := s.repo.GetJobByID(ctx, id)
	if err != nil {
		return domain.ErrNotFound
	}

	// only the owner or an admin can delete.
	if role != domain.RoleAdmin && job.UserID != userID {
		return domain.ErrForbidden
	}

	if err := s.repo.DeleteJobByID(ctx, id); err != nil {
		return err
	}

	// invalidate individual job, all-jobs list, and owner's job list.
	_ = s.cache.Delete(ctx, util.GenerateCacheKey(jobPrefix, id))
	_ = s.cache.DeleteByPrefix(ctx, allJobsPrefix)
	_ = s.cache.Delete(ctx, util.GenerateCacheKey(userJobsPrefix, job.UserID))

	return nil
}

func toJobResponse(job *domain.Job) *domain.JobResponse {
	return &domain.JobResponse{
		ID:              job.ID,
		Title:           job.Title,
		Company:         job.Company,
		Location:        job.Location,
		Type:            job.Type,
		ExperienceLevel: job.ExperienceLevel,
		Description:     job.Description,
		Requirements:    job.Requirements,
		SalaryMin:       job.SalaryMin,
		SalaryMax:       job.SalaryMax,
		IsActive:        job.IsActive,
		PostedBy: domain.UserResponse{
			ID:    job.User.ID,
			Email: job.User.Email,
			Name:  job.User.Name,
			Role:  job.User.Role,
		},
	}
}

package service

import (
	"context"
	"time"

	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/domain"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/port"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/util"
	"golang.org/x/crypto/bcrypt"
)

const (
	userPrefix   = "user"
	allUsersKey  = "user:all"
	userCacheTTL = 5 * time.Minute
)

type UserService struct {
	repo  port.UserRepository
	cache port.CacheRepository
}

func NewUserService(repo port.UserRepository, cache port.CacheRepository) *UserService {
	return &UserService{repo: repo, cache: cache}
}

func (s *UserService) RegisterUser(ctx context.Context, user *domain.User) (*domain.UserResponse, error) {
	existing, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err == nil && existing != nil {
		return nil, domain.ErrDuplicateEmail
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	if user.Role == 0 {
		user.Role = domain.RoleUser
	}

	created, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// invalidate all users cache since the list has changed
	_ = s.cache.Delete(ctx, allUsersKey)

	return util.ToUserResponse(created), nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (*domain.UserResponse, error) {
	key := util.GenerateCacheKey(userPrefix, id)

	// cache hit: return early.
	if cached, err := s.cache.Get(ctx, key); err == nil {
		var resp domain.UserResponse
		if err := util.Deserialize([]byte(cached), &resp); err == nil {
			return &resp, nil
		}
	}

	// cache miss: fetch from DB.
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	resp := util.ToUserResponse(user)

	// populate cache, ignore error to not block the response.
	if value, err := util.Serialize(resp); err == nil {
		_ = s.cache.Set(ctx, key, value, userCacheTTL)
	}

	return resp, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]domain.UserResponse, error) {
	// cache hit: return early.
	if cached, err := s.cache.Get(ctx, allUsersKey); err == nil {
		var responses []domain.UserResponse
		if err := util.Deserialize([]byte(cached), &responses); err == nil {
			return responses, nil
		}
	}

	// cache miss: fetch from DB.
	users, err := s.repo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]domain.UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, *util.ToUserResponse(&u))
	}

	// populate cache.
	if value, err := util.Serialize(responses); err == nil {
		_ = s.cache.Set(ctx, allUsersKey, value, userCacheTTL)
	}

	return responses, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id uint, updates *domain.User) (*domain.UserResponse, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	if updates.Name != "" {
		user.Name = updates.Name
	}

	if updates.Email != "" && updates.Email != user.Email {
		existing, _ := s.repo.GetUserByEmail(ctx, updates.Email)
		if existing != nil {
			return nil, domain.ErrDuplicateEmail
		}
		user.Email = updates.Email
	}

	if updates.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(updates.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashed)
	}

	if updates.Role != 0 {
		user.Role = updates.Role
	}

	updated, err := s.repo.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// invalidate both the individual entry and the list.
	_ = s.cache.Delete(ctx, util.GenerateCacheKey(userPrefix, id), allUsersKey)

	return util.ToUserResponse(updated), nil
}

func (s *UserService) DeleteUserByID(ctx context.Context, id uint) error {
	_, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return domain.ErrNotFound
	}

	if err := s.repo.DeleteUserByID(ctx, id); err != nil {
		return err
	}

	// invalidate both the individual entry and the list.
	_ = s.cache.Delete(ctx, util.GenerateCacheKey(userPrefix, id), allUsersKey)

	return nil
}

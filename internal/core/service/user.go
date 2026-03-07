package service

import (
	"context"

	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/domain"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/port"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/util"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo port.UserRepository
}

func NewUserService(repo port.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, user *domain.User) (*domain.UserResponse, error) {
	existing, _ := s.repo.GetUserByEmail(ctx, user.Email)
	if existing != nil {
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

	return util.ToUserResponse(created), nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (*domain.UserResponse, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	return util.ToUserResponse(user), nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]domain.UserResponse, error) {
	users, err := s.repo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]domain.UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, *util.ToUserResponse(&u))
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

	return util.ToUserResponse(updated), nil
}

func (s *UserService) DeleteUserByID(ctx context.Context, id uint) error {
	_, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return domain.ErrNotFound
	}

	return s.repo.DeleteUserByID(ctx, id)
}

package port

import (
	"context"

	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, id uint) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	DeleteUserByID(ctx context.Context, id uint) error
	GetAllUsers(ctx context.Context) ([]domain.User, error)
}

type UserService interface {
	RegisterUser(ctx context.Context, user *domain.User) (*domain.UserResponse, error)
	GetUserByID(ctx context.Context, id uint) (*domain.UserResponse, error)
	GetAllUsers(ctx context.Context) ([]domain.UserResponse, error)
	UpdateUser(ctx context.Context, id uint, updates *domain.User) (*domain.UserResponse, error)
	DeleteUserByID(ctx context.Context, id uint) error
}

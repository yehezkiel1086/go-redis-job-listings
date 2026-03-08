package port

import (
	"context"

	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/domain"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (*domain.LoginResponse, error)
	Refresh(ctx context.Context, rawRefreshToken string) (*domain.LoginResponse, error)
	Logout(ctx context.Context, rawRefreshToken string) error
}

package service

import (
	"context"

	"github.com/yehezkiel1086/go-redis-job-listings/internal/adapter/config"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/domain"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/port"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/util"
)

const refreshTokenPrefix = "refresh_token"

type AuthService struct {
	conf     *config.JWT
	userRepo port.UserRepository
	cache    port.CacheRepository
}

func NewAuthService(conf *config.JWT, userRepo port.UserRepository, cache port.CacheRepository) *AuthService {
	return &AuthService{
		conf:     conf,
		userRepo: userRepo,
		cache:    cache,
	}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*domain.LoginResponse, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if err := util.ComparePassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	accessToken, err := util.GenerateJWTToken(domain.AccessToken, s.conf, user)
	if err != nil {
		return nil, domain.ErrInternalServerError
	}

	refreshToken, err := util.GenerateJWTToken(domain.RefreshToken, s.conf, user)
	if err != nil {
		return nil, domain.ErrInternalServerError
	}

	value, err := util.Serialize(user.ID)
	if err != nil {
		return nil, domain.ErrInternalServerError
	}

	if err := s.cache.Set(ctx, util.GenerateCacheKey(refreshTokenPrefix, refreshToken), value, util.RefreshTokenTTL(s.conf)); err != nil {
		return nil, domain.ErrInternalServerError
	}

	return &domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, rawRefreshToken string) (*domain.LoginResponse, error) {
	claims, err := util.ParseJWTToken(domain.RefreshToken, s.conf, rawRefreshToken)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	key := util.GenerateCacheKey(refreshTokenPrefix, rawRefreshToken)
	if _, err := s.cache.Get(ctx, key); err != nil {
		return nil, domain.ErrInvalidToken
	}

	user, err := s.userRepo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	// Rotate: revoke old token, issue new pair.
	_ = s.cache.Delete(ctx, key)

	accessToken, err := util.GenerateJWTToken(domain.AccessToken, s.conf, user)
	if err != nil {
		return nil, domain.ErrInternalServerError
	}

	newRefreshToken, err := util.GenerateJWTToken(domain.RefreshToken, s.conf, user)
	if err != nil {
		return nil, domain.ErrInternalServerError
	}

	value, err := util.Serialize(user.ID)
	if err != nil {
		return nil, domain.ErrInternalServerError
	}

	if err := s.cache.Set(ctx, util.GenerateCacheKey(refreshTokenPrefix, newRefreshToken), value, util.RefreshTokenTTL(s.conf)); err != nil {
		return nil, domain.ErrInternalServerError
	}

	return &domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, rawRefreshToken string) error {
	key := util.GenerateCacheKey(refreshTokenPrefix, rawRefreshToken)
	if _, err := s.cache.Get(ctx, key); err != nil {
		return domain.ErrInvalidToken
	}
	return s.cache.Delete(ctx, key)
}

func (s *AuthService) ValidateAccessToken(tokenString string) (*domain.JWTClaims, error) {
	claims, err := util.ParseJWTToken(domain.AccessToken, s.conf, tokenString)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}
	return claims, nil
}

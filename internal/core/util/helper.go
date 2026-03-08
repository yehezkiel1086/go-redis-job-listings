package util

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/adapter/config"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/domain"
)

// user helpers
func ToUserResponse(user *domain.User) *domain.UserResponse {
	return &domain.UserResponse{
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func ParseID(c *gin.Context) (uint, error) {
	raw, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(raw), nil
}

// auth helpers
func RefreshTokenTTL(conf *config.JWT) time.Duration {
	days, err := strconv.Atoi(conf.RefreshTokenDuration)
	if err != nil {
		return 7 * 24 * time.Hour
	}
	return time.Duration(days) * 24 * time.Hour
}

package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/domain"
)

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

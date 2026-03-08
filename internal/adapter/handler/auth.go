package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/adapter/config"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/domain"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/port"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/util"
)

type AuthHandler struct {
	conf *config.JWT
	svc  port.AuthService
}

func NewAuthHandler(conf *config.JWT, svc port.AuthService) *AuthHandler {
	return &AuthHandler{
		conf,
		svc,
	}
}

// Login godoc
// @Summary      Login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      domain.LoginRequest   true  "Credentials"
// @Success      200   {object}  domain.LoginResponse
// @Failure      400   {object}  util.ErrorResponse
// @Failure      401   {object}  util.ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse{Error: err.Error()})
		return
	}

	resp, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, util.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, util.ErrorResponse{Error: err.Error()})
		return
	}

	accessTokenDuration, err := strconv.Atoi(h.conf.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorResponse{Error: err.Error()})
		return
	}
	refreshTokenDuration, err := strconv.Atoi(h.conf.RefreshTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorResponse{Error: err.Error()})
		return
	}
	c.SetCookie(string(domain.AccessToken), resp.AccessToken, accessTokenDuration*60, "/api/v1", "", false, true)
	c.SetCookie(string(domain.RefreshToken), resp.RefreshToken, refreshTokenDuration*60*60*24, "/api/v1/refresh", "", false, true)
	c.SetCookie(string(domain.RefreshToken), resp.RefreshToken, refreshTokenDuration*60*60*24, "/api/v1/logout", "", false, true)

	c.JSON(http.StatusOK, resp)
}

// Refresh godoc
// @Summary      Refresh access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      domain.RefreshRequest  true  "Refresh token"
// @Success      200   {object}  domain.LoginResponse
// @Failure      400   {object}  util.ErrorResponse
// @Failure      401   {object}  util.ErrorResponse
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie(string(domain.RefreshToken))
	if err != nil {
		c.JSON(http.StatusUnauthorized, util.ErrorResponse{Error: domain.ErrUnauthorized.Error()})
		return
	}

	resp, err := h.svc.Refresh(c.Request.Context(), refreshToken)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidToken) {
			c.JSON(http.StatusUnauthorized, util.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, util.ErrorResponse{Error: err.Error()})
		return
	}

	accessTokenDuration, err := strconv.Atoi(h.conf.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorResponse{Error: err.Error()})
		return
	}
	c.SetCookie(string(domain.AccessToken), resp.AccessToken, accessTokenDuration*60, "/api/v1", "", false, true)
	// Also rotate the refresh token cookie.
	refreshTokenDuration, err := strconv.Atoi(h.conf.RefreshTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorResponse{Error: err.Error()})
		return
	}
	c.SetCookie(string(domain.RefreshToken), resp.RefreshToken, refreshTokenDuration*60*60*24, "/api/v1/refresh", "", false, true)
	c.SetCookie(string(domain.RefreshToken), resp.RefreshToken, refreshTokenDuration*60*60*24, "/api/v1/logout", "", false, true)

	c.JSON(http.StatusOK, resp)
}

// Logout godoc
// @Summary      Logout (revoke refresh token)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  domain.RefreshRequest  true  "Refresh token to revoke"
// @Success      204   "No Content"
// @Failure      400   {object}  util.ErrorResponse
// @Failure      401   {object}  util.ErrorResponse
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie(string(domain.RefreshToken))
	if err != nil {
		c.JSON(http.StatusUnauthorized, util.ErrorResponse{Error: domain.ErrUnauthorized.Error()})
		return
	}

	if err := h.svc.Logout(c.Request.Context(), refreshToken); err != nil {
		if errors.Is(err, domain.ErrInvalidToken) {
			c.JSON(http.StatusUnauthorized, util.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, util.ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

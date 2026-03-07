package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/domain"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/port"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/util"
)

type UserHandler struct {
	svc port.UserService
}

func NewUserHandler(svc port.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// RegisterUser godoc
// @Summary      Register a new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body  body      domain.User          true  "User payload"
// @Success      201   {object}  domain.UserResponse
// @Failure      400   {object}  errorResponse
// @Failure      409   {object}  errorResponse
// @Failure      500   {object}  errorResponse
// @Router       /users/register [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var user domain.RegisterRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse{Error: err.Error()})
		return
	}

	resp, err := h.svc.RegisterUser(c.Request.Context(), &domain.User{
		Email:    user.Email,
		Password: user.Password,
		Name:     user.Name,
	})
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateEmail) {
			c.JSON(http.StatusConflict, util.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, util.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetUserByID godoc
// @Summary      Get a user by ID
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  domain.UserResponse
// @Failure      400  {object}  errorResponse
// @Failure      404  {object}  errorResponse
// @Router       /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := util.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse{Error: "invalid user id"})
		return
	}

	resp, err := h.svc.GetUserByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, util.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, util.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetAllUsers godoc
// @Summary      List all users
// @Tags         users
// @Produce      json
// @Success      200  {array}   domain.UserResponse
// @Failure      500  {object}  errorResponse
// @Router       /users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	resp, err := h.svc.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateUser godoc
// @Summary      Update a user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      int          true  "User ID"
// @Param        body  body      domain.User  true  "Fields to update"
// @Success      200   {object}  domain.UserResponse
// @Failure      400   {object}  errorResponse
// @Failure      404   {object}  errorResponse
// @Failure      409   {object}  errorResponse
// @Failure      500   {object}  errorResponse
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := util.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse{Error: "invalid user id"})
		return
	}

	var updates domain.User
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse{Error: err.Error()})
		return
	}

	resp, err := h.svc.UpdateUser(c.Request.Context(), id, &updates)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			c.JSON(http.StatusNotFound, util.ErrorResponse{Error: err.Error()})
		case errors.Is(err, domain.ErrDuplicateEmail):
			c.JSON(http.StatusConflict, util.ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, util.ErrorResponse{Error: err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteUserByID godoc
// @Summary      Delete a user by ID
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      204  "No Content"
// @Failure      400  {object}  errorResponse
// @Failure      404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUserByID(c *gin.Context) {
	id, err := util.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse{Error: "invalid user id"})
		return
	}

	if err := h.svc.DeleteUserByID(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, util.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, util.ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

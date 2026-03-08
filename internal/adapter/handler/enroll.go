package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/domain"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/port"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/util"
)

type EnrollmentHandler struct {
	svc port.EnrollmentService
}

func NewEnrollmentHandler(svc port.EnrollmentService) *EnrollmentHandler {
	return &EnrollmentHandler{svc: svc}
}

// EnrollJob godoc
// @Summary      Apply to a job listing
// @Tags         enrollments
// @Produce      json
// @Param        id   path      int  true  "Job ID"
// @Success      201  {object}  domain.EnrollmentResponse
// @Failure      400  {object}  util.ErrorResponse
// @Failure      401  {object}  util.ErrorResponse
// @Failure      404  {object}  util.ErrorResponse
// @Failure      409  {object}  util.ErrorResponse
// @Failure      500  {object}  util.ErrorResponse
// @Router       /jobs/{id}/enroll [post]
func (h *EnrollmentHandler) EnrollJob(c *gin.Context) {
	claims, ok := GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, util.ErrorResponse{Error: domain.ErrUnauthorized.Error()})
		return
	}

	jobID, err := util.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse{Error: "invalid job id"})
		return
	}

	resp, err := h.svc.EnrollJob(c.Request.Context(), claims.UserID, jobID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			c.JSON(http.StatusNotFound, util.ErrorResponse{Error: err.Error()})
		case errors.Is(err, domain.ErrInactiveJob):
			c.JSON(http.StatusBadRequest, util.ErrorResponse{Error: err.Error()})
		case errors.Is(err, domain.ErrDuplicateEnroll):
			c.JSON(http.StatusConflict, util.ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, util.ErrorResponse{Error: err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetMyEnrollments godoc
// @Summary      Get the logged-in user's enrollments
// @Tags         enrollments
// @Produce      json
// @Success      200  {array}   domain.EnrollmentResponse
// @Failure      401  {object}  util.ErrorResponse
// @Failure      500  {object}  util.ErrorResponse
// @Router       /enrollments/me [get]
func (h *EnrollmentHandler) GetMyEnrollments(c *gin.Context) {
	claims, ok := GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, util.ErrorResponse{Error: domain.ErrUnauthorized.Error()})
		return
	}

	resp, err := h.svc.GetMyEnrollments(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetJobEnrollments godoc
// @Summary      Get all enrollments for a job (admin or job owner)
// @Tags         enrollments
// @Produce      json
// @Param        id   path      int  true  "Job ID"
// @Success      200  {array}   domain.EnrollmentResponse
// @Failure      400  {object}  util.ErrorResponse
// @Failure      401  {object}  util.ErrorResponse
// @Failure      403  {object}  util.ErrorResponse
// @Failure      404  {object}  util.ErrorResponse
// @Failure      500  {object}  util.ErrorResponse
// @Router       /jobs/{id}/enrollments [get]
func (h *EnrollmentHandler) GetJobEnrollments(c *gin.Context) {
	claims, ok := GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, util.ErrorResponse{Error: domain.ErrUnauthorized.Error()})
		return
	}

	jobID, err := util.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse{Error: "invalid job id"})
		return
	}

	resp, err := h.svc.GetJobEnrollments(c.Request.Context(), jobID, claims.UserID, claims.Role)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			c.JSON(http.StatusNotFound, util.ErrorResponse{Error: err.Error()})
		case errors.Is(err, domain.ErrForbidden):
			c.JSON(http.StatusForbidden, util.ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, util.ErrorResponse{Error: err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateEnrollmentStatus godoc
// @Summary      Accept or reject an enrollment (admin or job owner)
// @Tags         enrollments
// @Accept       json
// @Produce      json
// @Param        id    path      int                                     true  "Enrollment ID"
// @Param        body  body      domain.UpdateEnrollmentStatusRequest    true  "Status update"
// @Success      200   {object}  domain.EnrollmentResponse
// @Failure      400   {object}  util.ErrorResponse
// @Failure      401   {object}  util.ErrorResponse
// @Failure      403   {object}  util.ErrorResponse
// @Failure      404   {object}  util.ErrorResponse
// @Failure      500   {object}  util.ErrorResponse
// @Router       /enrollments/{id} [put]
func (h *EnrollmentHandler) UpdateEnrollmentStatus(c *gin.Context) {
	claims, ok := GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, util.ErrorResponse{Error: domain.ErrUnauthorized.Error()})
		return
	}

	id, err := util.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse{Error: "invalid enrollment id"})
		return
	}

	var req domain.UpdateEnrollmentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse{Error: err.Error()})
		return
	}

	resp, err := h.svc.UpdateEnrollmentStatus(c.Request.Context(), id, claims.UserID, claims.Role, &req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			c.JSON(http.StatusNotFound, util.ErrorResponse{Error: err.Error()})
		case errors.Is(err, domain.ErrForbidden):
			c.JSON(http.StatusForbidden, util.ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, util.ErrorResponse{Error: err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteEnrollmentByID godoc
// @Summary      Withdraw an enrollment (applicant or admin)
// @Tags         enrollments
// @Produce      json
// @Param        id   path      int  true  "Enrollment ID"
// @Success      204  "No Content"
// @Failure      400  {object}  util.ErrorResponse
// @Failure      401  {object}  util.ErrorResponse
// @Failure      403  {object}  util.ErrorResponse
// @Failure      404  {object}  util.ErrorResponse
// @Failure      500  {object}  util.ErrorResponse
// @Router       /enrollments/{id} [delete]
func (h *EnrollmentHandler) DeleteEnrollmentByID(c *gin.Context) {
	claims, ok := GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, util.ErrorResponse{Error: domain.ErrUnauthorized.Error()})
		return
	}

	id, err := util.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse{Error: "invalid enrollment id"})
		return
	}

	if err := h.svc.DeleteEnrollmentByID(c.Request.Context(), id, claims.UserID, claims.Role); err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			c.JSON(http.StatusNotFound, util.ErrorResponse{Error: err.Error()})
		case errors.Is(err, domain.ErrForbidden):
			c.JSON(http.StatusForbidden, util.ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, util.ErrorResponse{Error: err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

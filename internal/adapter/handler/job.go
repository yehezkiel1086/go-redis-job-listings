package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/domain"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/port"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/util"
)

type JobHandler struct {
	svc port.JobService
}

func NewJobHandler(svc port.JobService) *JobHandler {
	return &JobHandler{svc: svc}
}

// CreateJob godoc
// @Summary      Create a new job listing
// @Tags         jobs
// @Accept       json
// @Produce      json
// @Param        body  body      domain.CreateJobRequest  true  "Job payload"
// @Success      201   {object}  domain.JobResponse
// @Failure      400   {object}  util.ErrorResponse
// @Failure      401   {object}  util.ErrorResponse
// @Failure      500   {object}  util.ErrorResponse
// @Router       /jobs [post]
func (h *JobHandler) CreateJob(c *gin.Context) {
	claims, ok := GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, util.ErrorResponse{Error: domain.ErrUnauthorized.Error()})
		return
	}

	var req domain.CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse{Error: err.Error()})
		return
	}

	resp, err := h.svc.CreateJob(c.Request.Context(), claims.UserID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetJobByID godoc
// @Summary      Get a job by ID
// @Tags         jobs
// @Produce      json
// @Param        id   path      int  true  "Job ID"
// @Success      200  {object}  domain.JobResponse
// @Failure      400  {object}  util.ErrorResponse
// @Failure      404  {object}  util.ErrorResponse
// @Router       /jobs/{id} [get]
func (h *JobHandler) GetJobByID(c *gin.Context) {
	id, err := util.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse{Error: "invalid job id"})
		return
	}

	resp, err := h.svc.GetJobByID(c.Request.Context(), id)
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

// GetAllJobs godoc
// @Summary      List all job listings
// @Tags         jobs
// @Produce      json
// @Param        type              query     string  false  "Job type"
// @Param        experience_level  query     string  false  "Experience level"
// @Param        location          query     string  false  "Location"
// @Param        search            query     string  false  "Search keyword"
// @Param        is_active         query     bool    false  "Active status"
// @Success      200  {array}   domain.JobResponse
// @Failure      400  {object}  util.ErrorResponse
// @Failure      500  {object}  util.ErrorResponse
// @Router       /jobs [get]
func (h *JobHandler) GetAllJobs(c *gin.Context) {
	var filter domain.JobFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse{Error: err.Error()})
		return
	}

	resp, err := h.svc.GetAllJobs(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetMyJobs godoc
// @Summary      Get jobs posted by the logged-in user
// @Tags         jobs
// @Produce      json
// @Success      200  {array}   domain.JobResponse
// @Failure      401  {object}  util.ErrorResponse
// @Failure      500  {object}  util.ErrorResponse
// @Router       /jobs/me [get]
func (h *JobHandler) GetMyJobs(c *gin.Context) {
	claims, ok := GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, util.ErrorResponse{Error: domain.ErrUnauthorized.Error()})
		return
	}

	resp, err := h.svc.GetJobsByUserID(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetJobsByUserID godoc
// @Summary      Get jobs posted by a specific user (admin only)
// @Tags         jobs
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {array}   domain.JobResponse
// @Failure      400  {object}  util.ErrorResponse
// @Failure      500  {object}  util.ErrorResponse
// @Router       /jobs/user/{id} [get]
func (h *JobHandler) GetJobsByUserID(c *gin.Context) {
	userID, err := util.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse{Error: "invalid user id"})
		return
	}

	resp, err := h.svc.GetJobsByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateJob godoc
// @Summary      Update a job listing
// @Tags         jobs
// @Accept       json
// @Produce      json
// @Param        id    path      int                      true  "Job ID"
// @Param        body  body      domain.UpdateJobRequest  true  "Fields to update"
// @Success      200   {object}  domain.JobResponse
// @Failure      400   {object}  util.ErrorResponse
// @Failure      401   {object}  util.ErrorResponse
// @Failure      403   {object}  util.ErrorResponse
// @Failure      404   {object}  util.ErrorResponse
// @Failure      500   {object}  util.ErrorResponse
// @Router       /jobs/{id} [put]
func (h *JobHandler) UpdateJob(c *gin.Context) {
	claims, ok := GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, util.ErrorResponse{Error: domain.ErrUnauthorized.Error()})
		return
	}

	id, err := util.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse{Error: "invalid job id"})
		return
	}

	var req domain.UpdateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse{Error: err.Error()})
		return
	}

	resp, err := h.svc.UpdateJob(c.Request.Context(), id, claims.UserID, claims.Role, &req)
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

// DeleteJobByID godoc
// @Summary      Delete a job listing
// @Tags         jobs
// @Produce      json
// @Param        id   path      int  true  "Job ID"
// @Success      204  "No Content"
// @Failure      400  {object}  util.ErrorResponse
// @Failure      401  {object}  util.ErrorResponse
// @Failure      403  {object}  util.ErrorResponse
// @Failure      404  {object}  util.ErrorResponse
// @Failure      500  {object}  util.ErrorResponse
// @Router       /jobs/{id} [delete]
func (h *JobHandler) DeleteJobByID(c *gin.Context) {
	claims, ok := GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, util.ErrorResponse{Error: domain.ErrUnauthorized.Error()})
		return
	}

	id, err := util.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse{Error: "invalid job id"})
		return
	}

	if err := h.svc.DeleteJobByID(c.Request.Context(), id, claims.UserID, claims.Role); err != nil {
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

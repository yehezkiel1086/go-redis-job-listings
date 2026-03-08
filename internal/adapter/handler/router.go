package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/adapter/config"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/domain"
)

type Router struct {
	r *gin.Engine
}

func New(
	conf *config.JWT,
	userHandler *UserHandler,
	authHandler *AuthHandler,
	jobHandler *JobHandler,
	enrollmentHandler *EnrollmentHandler,
) (*Router, error) {
	r := gin.New()

	// RBAC
	pb := r.Group("/api/v1")
	us := pb.Group("/", Authenticate(conf))
	ad := us.Group("/", Authorize(domain.RoleAdmin))

	// auth routes - public
	pb.POST("/register", userHandler.RegisterUser)
	pb.POST("/login", authHandler.Login)
	pb.POST("/refresh", authHandler.Refresh)

	// auth routes - authenticated
	us.POST("/logout", authHandler.Logout)

	// user routes - authenticated
	us.GET("/users/:id", AuthorizeOwnerOrAdmin(), userHandler.GetUserByID)
	us.PUT("/users/:id", AuthorizeOwnerOrAdmin(), userHandler.UpdateUser)
	us.DELETE("/users/:id", AuthorizeOwnerOrAdmin(), userHandler.DeleteUserByID)

	// user routes - admin
	ad.GET("/users", userHandler.GetAllUsers)

	// job routes — public
	pb.GET("/jobs", jobHandler.GetAllJobs)
	pb.GET("/jobs/:id", jobHandler.GetJobByID)

	// job routes — authenticated
	us.GET("/jobs/me", jobHandler.GetMyJobs)
	us.GET("/jobs/user/:id", AuthorizeOwnerOrAdmin(), jobHandler.GetJobsByUserID)

	// job routes — admin
	ad.POST("/jobs", jobHandler.CreateJob)
	ad.PUT("/jobs/:id", jobHandler.UpdateJob)
	ad.DELETE("/jobs/:id", jobHandler.DeleteJobByID)

	// enrollment routes — authenticated
	us.POST("/jobs/:id/enroll", enrollmentHandler.EnrollJob)
	us.GET("/enrollments/me", enrollmentHandler.GetMyEnrollments)
	us.DELETE("/enrollments/:id", AuthorizeOwnerOrAdmin(), enrollmentHandler.DeleteEnrollmentByID)

	// enrollment routes - admin
	ad.PUT("/enrollments/:id", enrollmentHandler.UpdateEnrollmentStatus)
	ad.GET("/jobs/:id/enrollments", enrollmentHandler.GetJobEnrollments)

	return &Router{r: r}, nil
}

func (r *Router) Run(conf *config.HTTP) error {
	return r.r.Run(conf.Host + ":" + conf.Port)
}

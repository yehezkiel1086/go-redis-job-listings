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
) (*Router, error) {
	r := gin.New()

	// RBAC
	pb := r.Group("/api/v1")
	us := pb.Group("/", Authenticate(conf))
	ad := us.Group("/", Authorize(domain.RoleAdmin))

	// public user routes
	pb.POST("/register", userHandler.RegisterUser)
	pb.POST("/login", authHandler.Login)
	pb.POST("/refresh", authHandler.Refresh)

	// user user routes
	us.GET("/users/:id", AuthorizeOwnerOrAdmin(), userHandler.GetUserByID)
	us.PUT("/users/:id", AuthorizeOwnerOrAdmin(), userHandler.UpdateUser)
	us.DELETE("/users/:id", AuthorizeOwnerOrAdmin(), userHandler.DeleteUserByID)
	us.POST("/logout", authHandler.Logout)

	// admin user routes
	ad.GET("/users", userHandler.GetAllUsers)

	return &Router{r: r}, nil
}

func (r *Router) Run(conf *config.HTTP) error {
	uri := conf.Host + ":" + conf.Port
	return r.r.Run(uri)
}

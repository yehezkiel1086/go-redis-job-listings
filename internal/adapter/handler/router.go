package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/adapter/config"
)

type Router struct {
	r *gin.Engine
}

func New(
	userHandler *UserHandler,
) (*Router, error) {
	r := gin.New()

	// public routes
	pb := r.Group("/api/v1")
	{
		// user routes
		pb.GET("/users", userHandler.GetAllUsers)
		pb.POST("/register", userHandler.RegisterUser)
		pb.GET("/users/:id", userHandler.GetUserByID)
		pb.PUT("/users/:id", userHandler.UpdateUser)
		pb.DELETE("/users/:id", userHandler.DeleteUserByID)
	}

	return &Router{r: r}, nil
}

func (r *Router) Run(conf *config.HTTP) error {
	uri := conf.Host + ":" + conf.Port
	return r.r.Run(uri)
}

package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/yehezkiel1086/go-redis-job-listings/internal/adapter/config"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/adapter/handler"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/adapter/storage/postgres"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/adapter/storage/postgres/repository"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/adapter/storage/redis"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/domain"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/service"
)

func handleError(err error, msg string) {
	if err != nil {
		slog.Error(msg, "error", err)
		os.Exit(1)
	}
}

func main() {
	// load .env configs
	conf, err := config.New()
	handleError(err, "failed to load .env configs")
	slog.Info(".env configs loaded successfully", "app", conf.App.Name, "env", conf.App.Env)

	ctx := context.Background()

	// init postgres db
	db, err := postgres.New(ctx, conf.DB)
	handleError(err, "failed to init postgres db")
	slog.Info("postgres db initialized successfully", "db", conf.DB.Name)

	// migrate
	err = db.Migrate(&domain.User{})
	handleError(err, "failed to migrate postgres db")
	slog.Info("postgres db migrated successfully")

	// init redis
	redis, err := redis.New(ctx, conf.Redis)
	handleError(err, "failed to init redis")
	slog.Info("redis initialized successfully")

	defer redis.Close()

	// dependency injection
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	// init router
	router, err := handler.New(userHandler)
	handleError(err, "failed to init router")
	slog.Info("router initialized successfully")

	// serve backend api
	err = router.Run(conf.HTTP)
	handleError(err, "failed to serve backend api")
	slog.Info("backend api served successfully")
}

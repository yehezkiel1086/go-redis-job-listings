package config

import (
	"os"

	"github.com/joho/godotenv"
)

type (
	Container struct {
		App   *App
		HTTP  *HTTP
		DB    *DB
		Redis *Redis
		JWT   *JWT
	}

	App struct {
		Name string
		Env  string
	}

	HTTP struct {
		Host string
		Port string
	}

	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}

	Redis struct {
		Host     string
		Port     string
		Password string
		DB       string
	}

	JWT struct {
		AccessTokenSecret    string
		RefreshTokenSecret   string
		AccessTokenDuration  string
		RefreshTokenDuration string
	}
)

func New() (*Container, error) {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	}

	App := &App{
		Name: os.Getenv("APP_NAME"),
		Env:  os.Getenv("APP_ENV"),
	}

	HTTP := &HTTP{
		Host: os.Getenv("HTTP_HOST"),
		Port: os.Getenv("HTTP_PORT"),
	}

	DB := &DB{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}

	Redis := &Redis{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       os.Getenv("REDIS_DB"),
	}

	JWT := &JWT{
		AccessTokenSecret:    os.Getenv("ACCESS_TOKEN_SECRET"),
		RefreshTokenSecret:   os.Getenv("REFRESH_TOKEN_SECRET"),
		AccessTokenDuration:  os.Getenv("ACCESS_TOKEN_DURATION"),
		RefreshTokenDuration: os.Getenv("REFRESH_TOKEN_DURATION"),
	}

	return &Container{
		App:   App,
		HTTP:  HTTP,
		DB:    DB,
		Redis: Redis,
		JWT:   JWT,
	}, nil
}

package util

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/adapter/config"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/domain"
)

func GenerateJWTToken(tokenType domain.TokenType, conf *config.JWT, user *domain.User) (string, error) {
	var signingKey []byte
	var duration time.Duration

	switch tokenType {
	case domain.AccessToken:
		signingKey = []byte(conf.AccessTokenSecret)
		durationInt, err := strconv.Atoi(conf.AccessTokenDuration)
		if err != nil {
			return "", err
		}
		duration = time.Duration(durationInt) * time.Minute
	case domain.RefreshToken:
		signingKey = []byte(conf.RefreshTokenSecret)
		durationInt, err := strconv.Atoi(conf.RefreshTokenDuration)
		if err != nil {
			return "", err
		}
		duration = time.Duration(durationInt) * time.Hour * 24
	default:
		return "", domain.ErrInvalidToken
	}

	claims := &domain.JWTClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}

func ParseJWTToken(tokenType domain.TokenType, conf *config.JWT, tokenString string) (*domain.JWTClaims, error) {
	var signingKey []byte
	switch tokenType {
	case domain.AccessToken:
		signingKey = []byte(conf.AccessTokenSecret)
	case domain.RefreshToken:
		signingKey = []byte(conf.RefreshTokenSecret)
	default:
		return nil, domain.ErrInvalidToken
	}

	token, err := jwt.ParseWithClaims(tokenString, &domain.JWTClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrInvalidToken
		}
		return signingKey, nil
	})
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(*domain.JWTClaims)
	if !ok || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	return claims, nil
}

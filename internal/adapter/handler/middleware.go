package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/adapter/config"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/domain"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/util"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
	ClaimsKey           = "claims"
)

func Authenticate(conf *config.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		header := c.GetHeader(AuthorizationHeader)
		if strings.HasPrefix(header, BearerPrefix) {
			tokenString = strings.TrimPrefix(header, BearerPrefix)
		} else {
			cookie, err := c.Cookie(string(domain.AccessToken))
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": domain.ErrUnauthorized.Error()})
				return
			}
			tokenString = cookie
		}

		claims, err := util.ParseJWTToken(domain.AccessToken, conf, tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": domain.ErrInvalidToken.Error()})
			return
		}

		c.Set(ClaimsKey, claims)
		c.Next()
	}
}

func Authorize(roles ...domain.Role) gin.HandlerFunc {
	allowed := make(map[domain.Role]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(c *gin.Context) {
		claims, ok := c.MustGet(ClaimsKey).(*domain.JWTClaims)
		if !ok || claims == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": domain.ErrUnauthorized.Error()})
			return
		}

		if _, permitted := allowed[claims.Role]; !permitted {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": domain.ErrForbidden.Error()})
			return
		}

		c.Next()
	}
}

func AuthorizeOwnerOrAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := c.MustGet(ClaimsKey).(*domain.JWTClaims)
		if !ok || claims == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": domain.ErrUnauthorized.Error()})
			return
		}

		if claims.Role == domain.RoleAdmin {
			c.Next()
			return
		}

		resourceID, err := util.ParseID(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		if claims.UserID != resourceID {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": domain.ErrForbidden.Error()})
			return
		}

		c.Next()
	}
}

func GetClaims(c *gin.Context) (*domain.JWTClaims, bool) {
	raw, exists := c.Get(ClaimsKey)
	if !exists {
		return nil, false
	}
	claims, ok := raw.(*domain.JWTClaims)
	return claims, ok
}

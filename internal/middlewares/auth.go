package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/j94veron/auth-service-insu/pkg/redis"
	"github.com/j94veron/auth-service-insu/pkg/token"
)

type AuthMiddleware struct {
	tokenService *token.TokenService
	redisClient  *redis.Client
}

func NewAuthMiddleware(tokenService *token.TokenService, redisClient *redis.Client) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService: tokenService,
		redisClient:  redisClient,
	}
}

func (am *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := am.tokenService.VerifyToken(tokenString, false)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			c.Abort()
			return
		}

		// Verificar si el token está en Redis/blacklist
		ctx := context.Background()
		userID, err := am.redisClient.GetUserID(ctx, claims.TokenUuid)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token revoked or expired"})
			c.Abort()
			return
		}

		// Verificar que el userID en Redis coincida con el del token
		if userID != claims.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token ownership"})
			c.Abort()
			return
		}

		// Agregar los datos del token al contexto para usar en los handlers
		c.Set("userID", claims.UserID)
		c.Set("userName", claims.Name)
		c.Set("userLastName", claims.LastName)
		c.Set("commercialZone", claims.CommercialZone)
		c.Set("warehouse", claims.Warehouse)
		c.Set("roleID", claims.RoleID)
		c.Set("tokenUuid", claims.TokenUuid)

		c.Next()
	}
}

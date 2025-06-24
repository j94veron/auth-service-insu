package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/j94veron/auth-service-insu/internal/role"
	"github.com/j94veron/auth-service-insu/internal/user"
)

type PermissionMiddleware struct {
	userRepo     user.Repository
	roleRepo     role.Repository
	allowedRoles []string
}

func NewPermissionMiddleware(userRepo user.Repository, roleRepo role.Repository) *PermissionMiddleware {
	return &PermissionMiddleware{
		userRepo:     userRepo,
		roleRepo:     roleRepo,
		allowedRoles: []string{"USER_ROLE_SCAN", "ADMIN"},
	}
}

func (pm *PermissionMiddleware) HasPermission(endpoint string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		roleID, exists := c.Get("roleID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Role information missing"})
			c.Abort()
			return
		}

		// Check if the role matches the allowed roles
		roleName, err := pm.roleRepo.GetRoleName(roleID.(uint))
		if err != nil || !pm.roleIsAllowed(roleName) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this resource."})
			c.Abort()
			return
		}

		// Check if the user has any specific restrictions
		userRestricted, err := pm.userRepo.IsUserRestricted(userID.(uint))
		if err != nil || userRestricted {
			c.JSON(http.StatusForbidden, gin.H{"error": "Your account has restrictions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// roleIsAllowed checks if the given role is in the list of allowed roles
func (pm *PermissionMiddleware) roleIsAllowed(roleName string) bool {
	for _, allowedRole := range pm.allowedRoles {
		if allowedRole == roleName {
			return true
		}
	}
	return false
}

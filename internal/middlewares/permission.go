package middlewares

import (
	"net/http"

	"auth-aca/internal/role"
	"auth-aca/internal/user"
	"github.com/gin-gonic/gin"
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

		// Verificar si el rol coincide con los roles permitidos
		roleName, err := pm.roleRepo.GetRoleName(roleID.(uint)) // Asegúrate de tener este método para obtener el nombre del rol
		if err != nil || !pm.roleIsAllowed(roleName) {
			c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permiso para acceder a este recurso"})
			c.Abort()
			return
		}

		// Verificar si el rol tiene permiso para este endpoint
		//hasPermission, err := pm.roleRepo.CheckPermission(roleID.(uint), endpoint, c.Request.Method)
		//if err != nil || !hasPermission {
		//	c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permiso para acceder a este recurso"})
		//	c.Abort()
		//	return
		//}

		// Ejemplo de verificación adicional basada en el usuario
		// Por ejemplo, verifica si el usuario tiene alguna restricción específica
		userRestricted, err := pm.userRepo.IsUserRestricted(userID.(uint))
		if err != nil || userRestricted {
			c.JSON(http.StatusForbidden, gin.H{"error": "Tu cuenta tiene restricciones"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// roleIsAllowed verifica si el rol dado está en la lista de roles permitidos
func (pm *PermissionMiddleware) roleIsAllowed(roleName string) bool {
	for _, allowedRole := range pm.allowedRoles {
		if allowedRole == roleName {
			return true
		}
	}
	return false
}

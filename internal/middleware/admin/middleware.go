package middleware

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/bonsus/go-saas/internal/config"
	myredis "github.com/bonsus/go-saas/internal/redis"
	token "github.com/bonsus/go-saas/internal/utils/jwt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Permission(db *gorm.DB, requiredPermission string) fiber.Handler {
	// initRedis()
	return func(c *fiber.Ctx) error {
		// Ambil token dari header Authorization
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or invalid token",
			})
		}

		// Format token "Bearer <token>", jadi harus dipotong
		var tokenStr string
		fmt.Sscanf(tokenString, "Bearer %s", &tokenStr)

		// Parse token JWT
		cfg := config.GetConfig()
		claims, err := token.ParseToken(tokenStr, cfg.JWT.Key)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}
		cacheId := "admin:" + claims.Id
		var user *AdminRole
		err = myredis.GetData(cacheId, &user)
		if err != nil {
			repo := NewRepository(db)
			user, err = repo.getPermission(claims.Id)
			if err != nil {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Access Denied!",
				})
			}
			myredis.SetData(cacheId, user, 30*time.Minute)
		}

		c.Locals("admin_id", user.Id)
		c.Locals("admin_name", user.Name)
		c.Locals("admin_email", user.Email)
		c.Locals("admin_username", user.Username)
		c.Locals("admin_status", user.Status)
		c.Locals("role_id", user.RoleId)
		c.Locals("role_name", user.RoleName)
		c.Locals("role_permission", user.Permission)

		if user.Status != "active" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access Denied!",
			})
		}

		// cek permission
		if requiredPermission != "" {
			allowed, err := HasPermission(requiredPermission, user.Permission)
			if err != nil {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Access Denied!",
				})
			}
			if !allowed {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Access Denied!",
				})
			}
		}
		return c.Next()
	}
}

func HasPermission(requiredPermission string, permissionJSON json.RawMessage) (bool, error) {
	// Parse JSON RawMessage ke slice []string
	var permissions []string
	err := json.Unmarshal(permissionJSON, &permissions)
	if err != nil {
		return false, err
	}

	// Loop untuk cek permission
	for _, p := range permissions {
		if strings.EqualFold(p, requiredPermission) { // Case-insensitive check
			return true, nil
		}
	}
	return false, nil
}

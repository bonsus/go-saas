package user

import (
	"log/slog"

	"github.com/bonsus/go-saas/apps/user/auth"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitModule(router fiber.Router, db *gorm.DB) {
	slog.Debug("starting module user")

	auth.Run(router, db)
}

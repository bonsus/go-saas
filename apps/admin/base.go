package admin

import (
	"log/slog"

	"github.com/bonsus/go-saas/apps/admin/admins"
	"github.com/bonsus/go-saas/apps/admin/apps"
	"github.com/bonsus/go-saas/apps/admin/auth"
	"github.com/bonsus/go-saas/apps/admin/option"
	"github.com/bonsus/go-saas/apps/admin/product"
	"github.com/bonsus/go-saas/apps/admin/promotion"
	"github.com/bonsus/go-saas/apps/admin/users"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitModule(router fiber.Router, db *gorm.DB) {
	slog.Debug("starting module admin")

	option.Run(router, db)
	auth.Run(router, db)
	admins.Run(router, db)
	users.Run(router, db)
	apps.Run(router, db)
	product.Run(router, db)
	promotion.Run(router, db)
}

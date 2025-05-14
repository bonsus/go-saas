package option

import (
	"github.com/bonsus/go-saas/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitModule(router fiber.Router, db *gorm.DB) {
	repo := NewRepository(db)
	svc := *NewService(repo)
	handler := newHandler(svc)

	authRouter := router.Group("admin/options")
	{
		authRouter.Put("", middleware.Permission(db, ""), handler.updateHandler)
		authRouter.Get(":name", middleware.Permission(db, ""), handler.getHandler)
	}
}

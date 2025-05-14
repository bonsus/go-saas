package auth

import (
	middleware "github.com/bonsus/go-saas/internal/middleware/admin"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Run(router fiber.Router, db *gorm.DB) {
	repo := NewRepository(db)
	svc := *NewService(repo)
	handler := newHandler(svc)

	authRouter := router.Group("admin/auth")
	{
		authRouter.Post("login", handler.LoginHandler)
		authRouter.Post("register", handler.RegisterHandler)
		authRouter.Post("me", middleware.Permission(db, ""), handler.meHandler)
		authRouter.Put("me", middleware.Permission(db, ""), handler.updateHandler)
		authRouter.Put("update-password", middleware.Permission(db, ""), handler.updatePasswordHandler)
	}
}

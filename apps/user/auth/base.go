package auth

import (
	"github.com/bonsus/go-saas/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Run(router fiber.Router, db *gorm.DB) {
	repo := NewRepository(db)
	svc := *NewService(repo)
	handler := newHandler(svc)

	authRouter := router.Group("user/auth")
	{
		authRouter.Post("login", handler.LoginHandler)
		authRouter.Post("register", handler.RegisterHandler)
		authRouter.Get("me", middleware.Permission(db, ""), handler.meHandler)
		authRouter.Put("me", middleware.Permission(db, ""), handler.updateHandler)
		authRouter.Put("update-password", middleware.Permission(db, ""), handler.updatePasswordHandler)
	}
}

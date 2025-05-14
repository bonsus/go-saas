package promotion

import (
	middleware "github.com/bonsus/go-saas/internal/middleware/admin"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Run(router fiber.Router, db *gorm.DB) {
	repo := NewRepository(db)
	svc := *NewService(repo)
	handler := newHandler(svc)

	authRouter := router.Group("admin/promotions")
	{
		authRouter.Post("/create", middleware.Permission(db, ""), handler.CreateHandler)
		authRouter.Get("/", middleware.Permission(db, ""), handler.IndexHandler)
		authRouter.Get("/:id", middleware.Permission(db, ""), handler.ReadHandler)
		authRouter.Put("/:id", middleware.Permission(db, ""), handler.UpdateHandler)
		authRouter.Delete("/:id", middleware.Permission(db, ""), handler.DeleteHandler)
		authRouter.Put("/:id/status", middleware.Permission(db, ""), handler.UpdateStatusHandler)
	}
}

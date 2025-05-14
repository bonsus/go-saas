package media

import (
	"github.com/bonsus/go-saas/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitModule(router fiber.Router, db *gorm.DB) {
	repo := NewRepository(db)
	svc := *NewService(repo)
	handler := newHandler(svc)

	authRouter := router.Group("medias")
	{
		authRouter.Get("/", middleware.Permission(db, ""), handler.IndexHandler)
		authRouter.Post("/uploads", middleware.Permission(db, ""), handler.UploadHandler)
		authRouter.Get("/:id", middleware.Permission(db, ""), handler.ReadHandler)
		authRouter.Put("/:id", middleware.Permission(db, ""), handler.UpdateHandler)
		authRouter.Delete("/:id", middleware.Permission(db, ""), handler.DeleteHandler)
	}
}

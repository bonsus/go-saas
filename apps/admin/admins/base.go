package admins

import (
	middleware "github.com/bonsus/go-saas/internal/middleware/admin"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Run(router fiber.Router, db *gorm.DB) {
	repo := NewRepository(db)
	svc := *NewService(repo)
	handler := newHandler(svc)

	authRouter := router.Group("admin/admins")
	{
		authRouter.Get("/roles", middleware.Permission(db, ""), handler.RoleIndexHandler)
		authRouter.Post("/roles/create", middleware.Permission(db, ""), handler.RoleCreateHandler)
		authRouter.Put("/roles/:id", middleware.Permission(db, ""), handler.RoleUpdateHandler)
		authRouter.Delete("/roles/:id", middleware.Permission(db, ""), handler.RoleDeleteHandler)

		authRouter.Post("/create", middleware.Permission(db, ""), handler.CreateHandler)
		authRouter.Get("/", middleware.Permission(db, ""), handler.IndexHandler)
		authRouter.Get("/:id", middleware.Permission(db, ""), handler.ReadHandler)
		authRouter.Put("/:id", middleware.Permission(db, ""), handler.UpdateHandler)
		authRouter.Delete("/:id", middleware.Permission(db, ""), handler.DeleteHandler)
		authRouter.Put("/:id/status", middleware.Permission(db, ""), handler.UpdateStatusHandler)
		authRouter.Put("/:id/password", middleware.Permission(db, ""), handler.UpdatePasswordHandler)
	}
}

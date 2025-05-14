package apps

import (
	middleware "github.com/bonsus/go-saas/internal/middleware/admin"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Run(router fiber.Router, db *gorm.DB) {
	repo := NewRepository(db)
	svc := *NewService(repo)
	handler := newHandler(svc)

	authRouter := router.Group("admin/apps")
	{
		authRouter.Get("/plugins", middleware.Permission(db, ""), handler.PluginIndexHandler)
		authRouter.Post("/plugins/create", middleware.Permission(db, ""), handler.PluginCreateHandler)
		authRouter.Put("/plugins/:id", middleware.Permission(db, ""), handler.PluginUpdateHandler)
		authRouter.Delete("/plugins/:id", middleware.Permission(db, ""), handler.PluginDeleteHandler)
		authRouter.Get("/plugins/:id", middleware.Permission(db, ""), handler.PluginReadHandler)

		authRouter.Get("/moduls", middleware.Permission(db, ""), handler.ModulIndexHandler)
		authRouter.Post("/moduls/create", middleware.Permission(db, ""), handler.ModulCreateHandler)
		authRouter.Put("/moduls/:id", middleware.Permission(db, ""), handler.ModulUpdateHandler)
		authRouter.Put("/moduls/:id/status", middleware.Permission(db, ""), handler.ModulUpdateStatusHandler)
		authRouter.Delete("/moduls/:id", middleware.Permission(db, ""), handler.ModulDeleteHandler)
		authRouter.Get("/moduls/:id", middleware.Permission(db, ""), handler.ModulReadHandler)

		authRouter.Post("/moduls/features/create", middleware.Permission(db, ""), handler.FeatureCreateHandler)
		authRouter.Get("/moduls/features/:id", middleware.Permission(db, ""), handler.FeatureReadHandler)
		authRouter.Put("/moduls/features/:id", middleware.Permission(db, ""), handler.FeatureUpdateHandler)
		authRouter.Put("/moduls/features/:id/status", middleware.Permission(db, ""), handler.FeatureUpdateStatusHandler)
		authRouter.Delete("/moduls/features/:id", middleware.Permission(db, ""), handler.FeatureDeleteHandler)
		authRouter.Post("/moduls/features/bulk-delete", middleware.Permission(db, ""), handler.FeatureBulkDeleteHandler)

		authRouter.Post("/create", middleware.Permission(db, ""), handler.CreateHandler)
		authRouter.Get("/", middleware.Permission(db, ""), handler.IndexHandler)
		authRouter.Get("/:id", middleware.Permission(db, ""), handler.ReadHandler)
		authRouter.Get("/:id/data", middleware.Permission(db, ""), handler.ReadDataHandler)
		authRouter.Put("/:id", middleware.Permission(db, ""), handler.UpdateHandler)
		authRouter.Delete("/:id", middleware.Permission(db, ""), handler.DeleteHandler)
		authRouter.Put("/:id/status", middleware.Permission(db, ""), handler.UpdateStatusHandler)
		authRouter.Put("/:id/db", middleware.Permission(db, ""), handler.UpdateDbHandler)
		authRouter.Post("/:id/db-test", middleware.Permission(db, ""), handler.DbTestHandler)
	}
}

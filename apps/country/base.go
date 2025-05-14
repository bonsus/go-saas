package country

import (
	"github.com/bonsus/go-saas/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitModule(router fiber.Router, db *gorm.DB) {
	repo := NewRepository(db)
	svc := *NewService(repo)
	handler := newHandler(svc)

	authRouter := router.Group("countries")
	{
		authRouter.Get("/provinces", middleware.Permission(db, ""), handler.ProvinceHandler)
		authRouter.Get("/cities", middleware.Permission(db, ""), handler.CityHandler)
		authRouter.Get("/districts", middleware.Permission(db, ""), handler.DistrictHandler)
		authRouter.Get("/zips", middleware.Permission(db, ""), handler.ZipHandler)
		authRouter.Get("/search", middleware.Permission(db, ""), handler.SearchHandler)
		authRouter.Get("/", middleware.Permission(db, ""), handler.CountryHandler)
	}
}

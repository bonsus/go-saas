package main

import (
	"os"
	"path/filepath"

	"github.com/bonsus/go-saas/apps/admin"
	"github.com/bonsus/go-saas/apps/user"
	_ "github.com/bonsus/go-saas/cmd/api/docs"
	"github.com/bonsus/go-saas/internal/config"
	"github.com/bonsus/go-saas/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// @title API Dokumentasi
// @version 1.0
// @description Ini adalah dokumentasi API untuk SaaS
// @host localhost:3001
// @BasePath /api/v1

func main() {
	filename := "./config.yaml"
	devMode := os.Getenv("DEV_MODE") == "true"
	if !devMode {
		execPath, _ := os.Executable()
		filename = filepath.Join(filepath.Dir(execPath), "config.yaml")
	}
	if err := config.LoadConfig(filename); err != nil {
		panic(err)
	}
	db, err := database.ConnectDB(config.Cfg.DB)
	if err != nil {
		panic(err)
	}

	// redis.InitRedis()

	app := fiber.New(fiber.Config{
		// Prefork: true,
		AppName: config.Cfg.App.Name,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173,http://localhost:5174,http://192.168.66.6:4321,http://192.168.100.98:4321",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Admin
	admin.InitModule(app, db)
	user.InitModule(app, db)

	// users.InitModule(app, db)
	// option.InitModule(app, db)
	// media.InitModule(app, db)
	// country.InitModule(app, db)

	app.Listen(":" + config.Cfg.App.Port)
}

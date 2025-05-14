package main

import (
	"github.com/bonsus/go-saas/internal/config"
	"github.com/bonsus/go-saas/internal/database"
)

func main() {
	filename := "./config.yaml"
	if err := config.LoadConfig(filename); err != nil {
		panic(err)
	}
	db, err := database.ConnectDB(config.Cfg.DB)
	if err != nil {
		panic(err)
	}
	database.RunMigrate(db)
}

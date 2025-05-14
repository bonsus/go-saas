package database

import (
	"fmt"

	"github.com/bonsus/go-saas/internal/database/migration"
	"gorm.io/gorm"
)

func RunMigrate(DB *gorm.DB) {
	err := DB.AutoMigrate(
		&migration.Option{},

		&migration.Admin{},
		&migration.AdminRole{},

		&migration.User{},
		&migration.UserPermission{},
		// &migration.UserApp{},
		// &migration.UserAppPlugin{},

		&migration.App{},
		&migration.AppModul{},
		&migration.AppModulFeature{},
		&migration.AppPlugin{},

		&migration.Product{},
		&migration.ProductModul{},
		&migration.ProductModulFeature{},
		&migration.ProductPrice{},
		&migration.ProductLimit{},
		&migration.Promotion{},
		&migration.PromotionItem{},
	)
	if err != nil {
		panic("error migrate user:" + err.Error())
	}
	fmt.Println("Migrate success")
}

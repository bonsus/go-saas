package media

import (
	"testing"

	"github.com/bonsus/go-saas/internal/config"
	"github.com/bonsus/go-saas/internal/database"
)

var svc service

func init() {
	filename := "../../config.yaml"
	err := config.LoadConfig(filename)
	if err != nil {
		panic(err)
	}
	db, err := database.ConnectDB(config.Cfg.DB)
	if err != nil {
		panic(err)
	}

	repo := NewRepository(db)
	svc = *NewService(repo)
}

func TestMedia(t *testing.T) {
	t.Run("UploadSuccess", func(t *testing.T) {
		// filePath := "../../storage/temp/80e168aa-d588-4685-95da-684d5cdf66ff_1726721191752-xxx4586kedgsizb7f4yc.png"

		// _, err := svc.Upload(filePath)
		// assert.Nil(t, err)
		// fmt.Println(err)
	})
}

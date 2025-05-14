package country

import (
	"context"
	"testing"

	"github.com/bonsus/go-saas/internal/config"
	"github.com/bonsus/go-saas/internal/database"
	"github.com/stretchr/testify/assert"
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

func TestAccount(t *testing.T) {
	t.Run("ProvinceSuccess", func(t *testing.T) {
		param := Params{
			Province: "Jawa Barat",
		}
		_, err := svc.Provinces(context.Background(), param)
		assert.Nil(t, err)
	})
	t.Run("CitySuccess", func(t *testing.T) {
		param := Params{
			Province: "Jawa Barat",
			City:     "Kota Sukabumi",
		}
		_, err := svc.Cities(context.Background(), param)
		assert.Nil(t, err)
	})
	t.Run("DistrictSuccess", func(t *testing.T) {
		param := Params{
			Province: "Jawa Barat",
			City:     "Kota Sukabumi",
			District: "Cikole",
		}
		_, err := svc.Districts(context.Background(), param)
		assert.Nil(t, err)
	})
	t.Run("SearchSuccess", func(t *testing.T) {
		param := Params{
			Search: "ci",
		}
		_, err := svc.Search(context.Background(), param)
		assert.Nil(t, err)
	})
}

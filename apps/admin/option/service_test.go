package option

import (
	"encoding/json"
	"fmt"
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

func TestSaveOption(t *testing.T) {
	t.Run("SaveOptionSuccess", func(t *testing.T) {
		req := Request{
			Name:  "general_test",
			Value: json.RawMessage(`{"zzzz":"halo message","zzz":"dua message"}`),
		}
		_, _, err := svc.Save(t.Context(), req)
		assert.Nil(t, err)
		fmt.Println(err)
	})
	t.Run("getOptionSuccess", func(t *testing.T) {
		data, err := svc.Get("general_test")
		assert.Nil(t, err)
		fmt.Println(err)
		fmt.Println(data)
	})
}

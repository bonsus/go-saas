package auth

import (
	"fmt"
	"testing"

	"github.com/bonsus/go-saas/internal/config"
	"github.com/bonsus/go-saas/internal/database"
	"github.com/google/uuid"
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

var email = fmt.Sprintf("%v@mail.id", uuid.NewString())
var username = uuid.NewString()

func TestAuth(t *testing.T) {
	t.Run("RegisterSuccess", func(t *testing.T) {
		req := RegisterRequest{
			Name:     "Joko",
			Email:    email,
			Username: username,
			Password: "P123456$$",
			RoleId:   "root",
			Status:   "active",
		}
		_, _, err := svc.Register(req)
		assert.Nil(t, err)
		fmt.Println(err)
	})
	t.Run("LoginSuccess", func(t *testing.T) {
		req := Request{
			Username: username,
			Password: "P123456$$",
		}
		_, _, err := svc.Login(t.Context(), req)
		assert.Nil(t, err)
		fmt.Println(err)
	})
}

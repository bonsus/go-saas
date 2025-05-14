package database

import (
	"testing"

	"github.com/bonsus/go-saas/internal/config"
	"github.com/stretchr/testify/require"
)

func init() {
	filename := "../../config.yaml"
	err := config.LoadConfig(filename)
	if err != nil {
		panic(err)
	}
}
func TestConnectDB(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, err := ConnectDB(config.Cfg.DB)
		require.Nil(t, err)
		require.NotNil(t, db)
	})
}

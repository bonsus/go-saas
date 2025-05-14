package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		filename := "./config.yaml"
		err := LoadConfig(filename)

		require.Nil(t, err)
		// log.Printf("%+v\n", Cfg)
	})
}
func BenchmarkLoadConfig(b *testing.B) {
	b.Run("success", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			filename := "./config.yaml"
			LoadConfig(filename)
		}
	})
}

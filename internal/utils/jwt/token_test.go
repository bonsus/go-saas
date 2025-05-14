package token

import (
	"fmt"
	"testing"
	"time"

	"github.com/bonsus/go-saas/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cfg := config.GetConfig()
		claims := Claims{
			Id:         "iniadalahid",
			ExpireTime: time.Duration(time.Now().Unix() + int64(cfg.JWT.ExpireTime)),
		}
		token, err := GenerateJWT(claims, cfg.JWT.Key)

		assert.Nil(t, err)
		fmt.Println(err)
		fmt.Println(token)
	})
}

func TestParseToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cfg := config.GetConfig()
		claims := Claims{
			Id:         "iniadalahidx",
			ExpireTime: time.Duration(time.Now().Unix() + int64(cfg.JWT.ExpireTime)),
		}
		token, err := GenerateJWT(claims, cfg.JWT.Key)
		fmt.Println(err)
		fmt.Println(token)

		result, err := ParseToken(token, cfg.JWT.Key)
		// result, err := ParseJWT(token, cfg.JWT.Key)
		fmt.Println(result.Id)
		fmt.Println(result.ExpireTime)
		fmt.Println(err)
		assert.Nil(t, err)
	})
}

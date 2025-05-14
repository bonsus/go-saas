package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Id         string        `json:"id"`
	Type       string        `json:"type"`
	ExpireTime time.Duration `json:"expire_time"`
	jwt.RegisteredClaims
}

func GenerateJWT(data Claims, secretKey string) (string, error) {
	var token string
	var err error

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	token, err = tok.SignedString([]byte(secretKey))
	return token, err
}

func ParseToken(token string, secretKey string) (claims *Claims, err error) {

	tok, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return claims, err
	}
	claimsTok, ok := tok.Claims.(*Claims)
	if !ok || !tok.Valid {
		return claims, errors.New("token not valid")
	}
	claims = claimsTok

	now := time.Now().Unix()
	tokenTime := claims.ExpireTime

	if int64(tokenTime) < now {
		return claims, errors.New("token expired")
	}

	return claims, nil
}

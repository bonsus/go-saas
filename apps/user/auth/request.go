package auth

import (
	"errors"

	"github.com/bonsus/go-saas/internal/utils/encryption"
	token "github.com/bonsus/go-saas/internal/utils/jwt"
)

type Request struct {
	Company                 string `json:"company"`
	Name                    string `json:"name"`
	Phone                   string `json:"phone"`
	Username                string `json:"username"`
	Email                   string `json:"email"`
	Password                string `json:"password"`
	NewPassword             string `json:"new_password"`
	NewPasswordConfirmation string `json:"new_password_confirmation"`
	ParentId                string `json:"parent_id"`
}

type RegisterRequest struct {
	Company              string `json:"company"`
	Name                 string `json:"name"`
	Phone                string `json:"phone"`
	Username             string `json:"username"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

func (u User) ValidatePassword(password string) error {
	if err := encryption.ValidatePassword(u.Password, password); err != nil {
		return errors.New("email or password is not match")
	}

	return nil
}

func (u User) GenerateToken(data token.Claims, secretKey string) (string, error) {
	token, err := token.GenerateJWT(data, secretKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

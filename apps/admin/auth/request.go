package auth

import (
	"errors"

	"github.com/bonsus/go-saas/internal/utils/encryption"
	token "github.com/bonsus/go-saas/internal/utils/jwt"
)

type Request struct {
	Email    string `json:"email" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Name                    string `json:"name" validate:"required"`
	Email                   string `json:"email" validate:"required,email"`
	Username                string `json:"username" validate:"required"`
	Password                string `json:"password" validate:"required"`
	NewPassword             string `json:"new_password" validate:"required"`
	NewPasswordConfirmation string `json:"new_password_confirmation" validate:"required"`
	RoleId                  string `json:"role_id" validate:"required"`
	Status                  string `json:"status"`
}

func (u Admin) ValidatePassword(password string) error {
	if err := encryption.ValidatePassword(u.Password, password); err != nil {
		return errors.New("email or password is not match")
	}

	return nil
}

func (u Admin) GenerateToken(data token.Claims, secretKey string) (string, error) {
	token, err := token.GenerateJWT(data, secretKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

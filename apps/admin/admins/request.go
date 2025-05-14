package admins

import "encoding/json"

type Request struct {
	Name                 string `json:"name" validate:"required"`
	Username             string `json:"username" validate:"required"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required"`
	RoleId               string `json:"role_id" validate:"required"`
	Status               string `json:"status"`
}
type ParamIndex struct {
	Page    int    `query:"page"`
	Perpage int    `query:"perpage"`
	Search  string `query:"search"`
}

type RoleRequest struct {
	Name       string          `json:"name"`
	Permission json.RawMessage `json:"permission"`
}

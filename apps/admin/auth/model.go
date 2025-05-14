package auth

import (
	"encoding/json"
	"time"
)

type Response struct {
	Token string `json:"token"`
	Admin Admin  `json:"admin"`
}
type Admin struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	RoleId    string    `json:"role_id"`
	Role      AdminRole `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AdminRole struct {
	Id         string          `json:"id"`
	Name       string          `json:"name"`
	Permission json.RawMessage `json:"permission"`
}

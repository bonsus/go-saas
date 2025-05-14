package users

import (
	"encoding/json"
	"time"
)

type User struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	RoleId    string    `json:"role_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Role      UserRole  `json:"role"`
}
type UserRole struct {
	Id         string          `json:"id"`
	Name       string          `json:"name"`
	Permission json.RawMessage `json:"permission"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

type UserIndex struct {
	Data      []User `json:"data"`
	Page      int    `json:"page"`
	PerPage   int    `json:"per_page"`
	TotalPage int    `json:"total_page"`
	Total     int64  `json:"total"`
}

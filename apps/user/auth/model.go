package auth

import (
	"encoding/json"
	"time"
)

type Response struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
type User struct {
	Id          string       `json:"id"`
	Company     string       `json:"company"`
	Name        string       `json:"name"`
	Phone       string       `json:"phone"`
	Username    string       `json:"username"`
	Email       string       `json:"email"`
	Password    string       `json:"password"`
	Type        string       `json:"type"`
	Status      string       `json:"status"`
	EmailStatus string       `json:"email_status"`
	PhoneStatus string       `json:"phone_status"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Permission  []Permission `json:"permissions" gorm:"-"`
}

type Permission struct {
	Id         string          `json:"id"`
	Name       string          `json:"name"`
	Permission json.RawMessage `json:"permission"`
}

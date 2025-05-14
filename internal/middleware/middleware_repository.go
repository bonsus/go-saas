package middleware

import (
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repository {
	return repository{
		db: db,
	}
}

type User struct {
	Id         string         `json:"id"`
	Company    string         `json:"company"`
	Name       string         `json:"name"`
	Phone      string         `json:"Phone"`
	Username   string         `json:"username"`
	Email      string         `json:"email"`
	Type       string         `json:"type"`
	Status     string         `json:"status"`
	Permission UserPermission `json:"permission"`
}
type UserPermission struct {
	Id         string          `json:"id"`
	UserId     string          `json:"user_id"`
	AppId      string          `json:"app_id"`
	Permission json.RawMessage `json:"permission"`
}

func (r *repository) getPermission(userId string) (user *User, err error) {
	result := r.db.Model(User{}).Find(&user)

	if result.Error != nil {
		return user, result.Error
	}
	if result.RowsAffected == 0 {
		return user, errors.New("user not found")
	}
	return user, nil
}

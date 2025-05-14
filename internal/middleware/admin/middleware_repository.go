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

type AdminRole struct {
	Id         string          `json:"id"`
	Name       string          `json:"name"`
	Email      string          `json:"email"`
	Username   string          `json:"username"`
	Status     string          `json:"status"`
	RoleId     string          `json:"role_id"`
	RoleName   string          `json:"role_name"`
	Permission json.RawMessage `json:"permission" gorm:"type:json"`
}

func (r *repository) getPermission(userId string) (role *AdminRole, err error) {
	result := r.db.Table("admins").
		Joins("JOIN admin_roles ON admins.role_id = admin_roles.id").
		Where("admins.id = ?", userId).
		Where("admins.status = ?", "active").
		Select("admins.id, admins.name as name, admins.email as email,admins.username as username, admins.status, admin_roles.id as role_id, admin_roles.name as role_name, admin_roles.permission as permission").
		Scan(&role)

	if result.Error != nil {
		return role, result.Error
	}
	if result.RowsAffected == 0 {
		return role, errors.New("user not found")
	}
	return role, nil
}

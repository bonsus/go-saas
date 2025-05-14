package auth

import (
	"context"
	"errors"
	"time"

	"github.com/bonsus/go-saas/internal/utils/encryption"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Login(admin Admin) (*Admin, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repository {
	return repository{
		db: db,
	}
}

func (r *repository) FindUserByEmail(ctx context.Context, req Request) (admin *Admin, err error) {
	result := r.db.Where("admins.email = ?", req.Email).
		Select("id, name, email, username, password, role_id, status,created_at,updated_at").
		Preload("Role", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "permission")
		}).
		First(&admin)

	if result.Error != nil {
		return admin, result.Error
	}
	if result.RowsAffected == 0 {
		return admin, errors.New("user not found")
	}
	return admin, nil
}

func (r *repository) FindUserByUsername(ctx context.Context, req Request) (admin *Admin, err error) {
	result := r.db.Where("admins.username = ?", req.Username).
		Select("id, name, email, username, password, role_id, status,created_at,updated_at").
		Preload("Role", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "permission")
		}).
		Find(&admin)

	if result.Error != nil {
		return admin, result.Error
	}
	if result.RowsAffected == 0 {
		return admin, errors.New("user not found")
	}
	return admin, nil
}

func (r *repository) FindByEmail(email string, Id ...string) bool {
	var exists bool
	if Id != nil {
		query := `SELECT EXISTS (SELECT 1 FROM admins WHERE email = ? AND id != ?)`
		if err := r.db.Raw(query, email, Id).Scan(&exists).Error; err != nil {
			return false
		}
	} else {
		query := `SELECT EXISTS (SELECT 1 FROM admins WHERE email = ?)`
		if err := r.db.Raw(query, email).Scan(&exists).Error; err != nil {
			return false
		}
	}
	return exists
}

func (r *repository) FindByUsername(username string, Id ...string) bool {
	var exists bool
	if Id != nil {
		query := `SELECT EXISTS (SELECT 1 FROM admins WHERE username = ? AND id != ?)`
		if err := r.db.Raw(query, username, Id).Scan(&exists).Error; err != nil {
			return false
		}
	} else {
		query := `SELECT EXISTS (SELECT 1 FROM admins WHERE username = ?)`
		if err := r.db.Raw(query, username).Scan(&exists).Error; err != nil {
			return false
		}
	}
	return exists
}

func (r *repository) Register(req RegisterRequest) (*Admin, error) {
	Id := uuid.NewString()
	query := `
		INSERT INTO admins (id, role_id, name, email, username, password, status, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	hashPassword, _ := encryption.HashPassword(req.Password)
	result := r.db.Exec(query, Id, req.RoleId, req.Name, req.Email, req.Username, hashPassword, req.Status, time.Now(), time.Now())
	if result.Error != nil {
		return nil, result.Error
	}
	user, _ := r.FindUserById(Id)
	return user, nil
}

func (r *repository) Update(req RegisterRequest, Id string) (*Admin, error) {

	result := r.db.Model(&Admin{}).Where("id = ?", Id).Updates(map[string]interface{}{
		"name":       req.Name,
		"email":      req.Email,
		"username":   req.Username,
		"updated_at": time.Now(),
	})

	if result.Error != nil {
		return nil, result.Error
	}
	user, _ := r.FindUserById(Id)
	return user, nil
}

func (r *repository) UpdatePassword(req RegisterRequest, Id string) (*Admin, error) {

	hashPassword, _ := encryption.HashPassword(req.NewPassword)
	result := r.db.Model(&Admin{}).Where("id = ?", Id).Updates(map[string]interface{}{
		"password":   hashPassword,
		"updated_at": time.Now(),
	})

	if result.Error != nil {
		return nil, result.Error
	}
	user, _ := r.FindUserById(Id)
	return user, nil
}

func (r *repository) FindRoleByID(id string) bool {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM admin_roles WHERE id = ?)`
	if err := r.db.Raw(query, id).Scan(&exists).Error; err != nil {
		return false
	}
	return exists
}

func (r *repository) FindUserById(adminId string) (admin *Admin, err error) {
	result := r.db.Model(Admin{}).Where("admins.id = ?", adminId).
		Select("id, role_id, name, email, username, status, created_at,updated_at").
		Preload("Role", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "permission")
		}).
		First(&admin)

	if result.Error != nil {
		return admin, result.Error
	}
	if result.RowsAffected == 0 {
		return admin, errors.New("user not found")
	}

	return admin, nil
}

func (r *repository) FindUserPasswordById(userId string) (admin *Admin, err error) {
	result := r.db.Where("admins.id = ?", userId).Select("id, password").Find(&admin)

	if result.Error != nil {
		return admin, result.Error
	}
	if result.RowsAffected == 0 {
		return admin, errors.New("user not found")
	}

	return admin, nil
}

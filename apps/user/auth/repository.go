package auth

import (
	"errors"
	"time"

	"github.com/bonsus/go-saas/internal/utils/encryption"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Login(user User) (*User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repository {
	return repository{
		db: db,
	}
}

func (r *repository) Register(req RegisterRequest) (*User, error) {
	Id := uuid.NewString()
	timeNow := time.Now()
	hashPassword, _ := encryption.HashPassword(req.Password)
	create := User{
		Id:          Id,
		Company:     req.Company,
		Name:        req.Name,
		Phone:       req.Phone,
		Username:    req.Username,
		Email:       req.Email,
		Password:    hashPassword,
		Type:        "owner",
		Status:      "unverify",
		EmailStatus: "unverify",
		PhoneStatus: "unverify",
		CreatedAt:   timeNow,
		UpdatedAt:   timeNow,
	}
	result := r.db.Create(&create)
	if result.Error != nil {
		return nil, errors.New("failed to register")
	}
	res, _ := r.Read(Id)
	return res, nil
}

func (r *repository) Login(req Request) (*User, error) {
	var user User
	check := r.db.Model(User{}).
		Where("email = ? or phone = ? or username = ?", req.Email, req.Phone, req.Username).
		Find(&user)
	if check.RowsAffected == 0 {
		if check.Error != nil {
			return nil, check.Error
		}
		return nil, errors.New("user not found")
	}
	return &user, nil
}
func (r *repository) LoginById(id string) (*User, error) {
	var user User
	check := r.db.Model(User{}).Where("id = ?", id).Find(&user)
	if check.RowsAffected == 0 {
		if check.Error != nil {
			return nil, check.Error
		}
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (r *repository) Update(req Request, Id string) (*User, error) {
	data := map[string]interface{}{
		"company":    req.Company,
		"name":       req.Name,
		"updated_at": time.Now(),
	}
	result := r.db.Model(&User{}).Where("id = ?", Id).Updates(&data)

	if result.Error != nil {
		return nil, result.Error
	}
	user, _ := r.Read(Id)
	return user, nil
}

func (r *repository) UpdatePassword(req Request, Id string) (*User, error) {

	hashPassword, _ := encryption.HashPassword(req.NewPassword)
	result := r.db.Model(&User{}).Where("id = ?", Id).Updates(map[string]interface{}{
		"password":   hashPassword,
		"updated_at": time.Now(),
	})

	if result.Error != nil {
		return nil, result.Error
	}
	user, _ := r.Read(Id)
	return user, nil
}
func (r *repository) Read(id string) (*User, error) {
	var user User
	check := r.db.Model(User{}).Where("id = ?", id).Find(&user)
	if check.RowsAffected == 0 {
		if check.Error != nil {
			return nil, check.Error
		}
		return nil, errors.New("user not found")
	}
	user.Password = ""
	return &user, nil
}
func (r *repository) ReadDetail(id string) (*User, error) {
	var user User
	check := r.db.Model(User{}).Where("id = ?", id).Find(&user)
	if check.RowsAffected == 0 {
		if check.Error != nil {
			return nil, check.Error
		}
		return nil, errors.New("user not found")
	}
	user.Password = ""
	return &user, nil
}

func (r *repository) CheckEmail(email string) bool {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE email = ?)`
	if err := r.db.Raw(query, email).Scan(&exists).Error; err != nil {
		return false
	}
	return exists
}
func (r *repository) CheckUsername(username string) bool {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE username = ?)`
	if err := r.db.Raw(query, username).Scan(&exists).Error; err != nil {
		return false
	}
	return exists
}
func (r *repository) CheckPhone(phone string) bool {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE phone = ?)`
	if err := r.db.Raw(query, phone).Scan(&exists).Error; err != nil {
		return false
	}
	return exists
}

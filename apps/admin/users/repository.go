package users

import (
	"errors"
	"time"

	"github.com/bonsus/go-saas/internal/utils/encryption"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	// Create(account Request) (*Request, error)
	// Index(param ParamIndex) (*AccountIndex, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repository {
	return repository{
		db: db,
	}
}

func (r *repository) Create(req Request) (*User, error) {
	Id := uuid.NewString()
	timeNow := time.Now()
	hashPassword, _ := encryption.HashPassword(req.Password)
	eventContest := User{
		Id:        Id,
		Name:      req.Name,
		Email:     req.Email,
		Username:  req.Username,
		Password:  hashPassword,
		Status:    req.Status,
		RoleId:    req.RoleId,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	result := r.db.Create(&eventContest)

	if result.Error != nil {
		return nil, result.Error
	}
	user, _ := r.Read(Id)
	return user, nil
}

func (r *repository) Update(req Request, Id string) (user *User, err error) {
	updateData := map[string]interface{}{
		"name":       req.Name,
		"email":      req.Email,
		"username":   req.Username,
		"status":     req.Status,
		"role_id":    req.RoleId,
		"updated_at": time.Now(),
	}
	result := r.db.Model(&User{}).Where("id = ?", Id).Updates(updateData)
	if result.Error != nil {
		return nil, errors.New("update failed")
	}
	user, _ = r.Read(Id)
	return user, nil
}

func (r *repository) UpdateStatus(req Request, Id string) (user *User, err error) {
	result := r.db.Model(&User{}).Where("id = ?", Id).Updates(map[string]interface{}{
		"status":     req.Status,
		"updated_at": time.Now(),
	})

	if result.Error != nil {
		return nil, result.Error
	}
	user, _ = r.Read(Id)
	return user, nil
}

func (r *repository) UpdatePassword(req Request, Id string) (user *User, err error) {
	hashPassword, _ := encryption.HashPassword(req.Password)
	result := r.db.Model(&User{}).Where("id = ?", Id).Updates(map[string]interface{}{
		"password":   hashPassword,
		"updated_at": time.Now(),
	})

	if result.Error != nil {
		return nil, result.Error
	}
	user = &User{
		Id: Id,
	}
	return user, nil
}

func (r *repository) Index(param ParamIndex) (*UserIndex, error) {
	var data []User
	var total int64

	if param.Page < 1 {
		param.Page = 1
	}
	if param.Perpage < 1 {
		param.Perpage = 20
	}

	offset := (param.Page - 1) * param.Perpage
	dbQuery := r.db.Model(&User{}).Preload("Role").
		Select("id,name,username,email,created_at,updated_at").
		Order("created_at DESC")
	if param.Search != "" {
		search := "%" + param.Search + "%"
		dbQuery = dbQuery.Where(
			"name ILIKE ? OR email ILIKE ? OR username ILIKE ?",
			search, search, search,
		)
	}
	dbQuery.Count(&total)
	dbQuery.Limit(param.Perpage).Offset(offset).Find(&data)

	totalPage := int((total + int64(param.Perpage) - 1) / int64(param.Perpage))

	return &UserIndex{
		Data:      data,
		Total:     total,
		TotalPage: totalPage,
		Page:      param.Page,
		PerPage:   param.Perpage,
	}, nil
}

func (r *repository) Read(id string) (*User, error) {
	var user User
	check := r.db.Model(User{}).Where("id = ?", id).Preload("Role").Find(&user)
	if check.RowsAffected == 0 {
		if check.Error != nil {
			return nil, check.Error
		}
		return nil, errors.New("id not found")
	}
	user.Password = ""
	return &user, nil
}

func (r *repository) Delete(id string) error {
	var user User
	if err := r.db.Where("id = ?", id).Find(&user).Error; err != nil {
		return errors.New("id not found")
	}

	delete := r.db.Where("id = ?", id).Delete(&User{})
	if delete.Error != nil {
		return errors.New("data cannot be deleted")
	}

	return nil
}

func (r *repository) IsEmailExists(email string, Id string) bool {
	if email == "" {
		return false
	}

	var count int64
	var err error

	if Id != "" {
		err = r.db.Raw("SELECT COUNT(id) FROM riders WHERE email = ? AND id != ?", email, Id).Scan(&count).Error
	} else {
		err = r.db.Raw("SELECT COUNT(id) FROM riders WHERE email = ?", email).Scan(&count).Error
	}

	if err != nil {
		return false
	}

	return count > 0
}

func (r *repository) IsPhoneExists(phone string, Id string) bool {
	if phone == "" {
		return false
	}

	var count int64
	var err error

	if Id != "" {
		err = r.db.Raw("SELECT COUNT(id) FROM riders WHERE phone = ? AND id != ?", phone, Id).Scan(&count).Error
	} else {
		err = r.db.Raw("SELECT COUNT(id) FROM riders WHERE phone = ?", phone).Scan(&count).Error
	}

	if err != nil {
		return false
	}

	return count > 0
}

func (r *repository) IsValidTeamId(team_id string) bool {
	var count int64
	err := r.db.Raw("SELECT COUNT(id) FROM rider_teams WHERE id = ?", team_id).Scan(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

func (r *repository) FindRoleByID(id string) bool {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM user_roles WHERE id = ?)`
	if err := r.db.Raw(query, id).Scan(&exists).Error; err != nil {
		return false
	}
	return exists
}

func (r *repository) FindRoleByName(name string, id string) bool {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM user_roles WHERE name = ? AND id != ?)`
	if err := r.db.Raw(query, name, id).Scan(&exists).Error; err != nil {
		return false
	}
	return exists
}

func (r *repository) FindByEmail(email string, Id ...string) bool {
	var exists bool
	if Id != nil {
		query := `SELECT EXISTS (SELECT 1 FROM users WHERE email = ? AND id != ?)`
		if err := r.db.Raw(query, email, Id).Scan(&exists).Error; err != nil {
			return false
		}
	} else {
		query := `SELECT EXISTS (SELECT 1 FROM users WHERE email = ?)`
		if err := r.db.Raw(query, email).Scan(&exists).Error; err != nil {
			return false
		}
	}
	return exists
}

func (r *repository) FindByUsername(username string, Id ...string) bool {
	var exists bool
	if Id != nil {
		query := `SELECT EXISTS (SELECT 1 FROM users WHERE username = ? AND id != ?)`
		if err := r.db.Raw(query, username, Id).Scan(&exists).Error; err != nil {
			return false
		}
	} else {
		query := `SELECT EXISTS (SELECT 1 FROM users WHERE username = ?)`
		if err := r.db.Raw(query, username).Scan(&exists).Error; err != nil {
			return false
		}
	}
	return exists
}

func (r *repository) RoleIndex() ([]UserRole, error) {
	var data []UserRole
	r.db.Model(&UserRole{}).Order("name ASC").Find(&data)
	return data, nil
}
func (r *repository) RoleCreate(req RoleRequest) (*UserRole, error) {
	Id := uuid.NewString()
	timeNow := time.Now()
	role := UserRole{
		Id:         Id,
		Name:       req.Name,
		Permission: req.Permission,
		CreatedAt:  timeNow,
		UpdatedAt:  timeNow,
	}

	result := r.db.Create(&role)

	if result.Error != nil {
		return nil, result.Error
	}
	role = UserRole{
		Id:         Id,
		Name:       req.Name,
		Permission: req.Permission,
		CreatedAt:  timeNow,
		UpdatedAt:  timeNow,
	}
	return &role, nil
}

func (r *repository) RoleUpdate(req RoleRequest, Id string) (*UserRole, error) {
	timeNow := time.Now()
	updateData := map[string]interface{}{
		"name":       req.Name,
		"permission": req.Permission,
		"updated_at": time.Now(),
	}
	result := r.db.Model(&UserRole{}).Where("id = ?", Id).Updates(updateData)
	if result.Error != nil {
		return nil, errors.New("update failed")
	}
	role := UserRole{
		Id:         Id,
		Name:       req.Name,
		Permission: req.Permission,
		CreatedAt:  timeNow,
		UpdatedAt:  timeNow,
	}
	return &role, nil
}

func (r *repository) RoleDelete(id string) error {
	if err := r.db.Where("id = ?", id).Find(&UserRole{}).Error; err != nil {
		return errors.New("id not found")
	}

	delete := r.db.Where("id = ?", id).Delete(&UserRole{})
	if delete.Error != nil {
		return errors.New("data cannot be deleted")
	}

	return nil
}

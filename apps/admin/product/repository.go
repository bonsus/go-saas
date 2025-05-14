package product

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

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

func (r *repository) Create(req Request) (*Product, error) {
	Id := uuid.NewString()
	timeNow := time.Now()
	slug, err := r.ensureUniqueSlug(req.Name, req.Slug, Id)
	if err != nil {
		return nil, err
	}
	eventContest := Product{
		Id:          Id,
		AppId:       req.AppId,
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		Image:       req.Image,
		Status:      req.Status,
		CreatedAt:   timeNow,
		UpdatedAt:   timeNow,
	}

	result := r.db.Create(&eventContest)

	if result.Error != nil {
		return nil, result.Error
	}
	res, _ := r.Read(Id)
	return res, nil
}

func (r *repository) Update(req Request, Id string) (res *Product, err error) {
	slug, _ := r.ensureUniqueSlug(req.Name, req.Slug, Id)
	updateData := map[string]interface{}{
		"app_id":      req.AppId,
		"name":        req.Name,
		"slug":        slug,
		"description": req.Description,
		"image":       req.Image,
		"status":      req.Status,
		"updated_at":  time.Now(),
	}
	result := r.db.Model(&Product{}).Where("id = ?", Id).Updates(updateData)
	if result.Error != nil {
		return nil, errors.New("update failed")
	}
	res, _ = r.Read(Id)
	return res, nil
}

func (r *repository) UpdateStatus(req Request, Id string) (res *Product, err error) {
	result := r.db.Model(&Product{}).Where("id = ?", Id).Updates(map[string]interface{}{
		"status":     req.Status,
		"updated_at": time.Now(),
	})

	if result.Error != nil {
		return nil, result.Error
	}
	res, _ = r.Read(Id)
	return res, nil
}

func (r *repository) Index(param ParamIndex) (*ProductIndex, error) {
	var data []Product
	var total int64

	if param.Page < 1 {
		param.Page = 1
	}
	if param.Perpage < 1 {
		param.Perpage = 20
	}

	offset := (param.Page - 1) * param.Perpage
	dbQuery := r.db.Model(&Product{}).
		Preload("App").
		Order("created_at DESC")
	if param.Search != "" {
		search := "%" + param.Search + "%"
		dbQuery = dbQuery.Where(
			"name ILIKE ? OR code ILIKE ? OR description ILIKE ?",
			search, search, search,
		)
	}
	dbQuery.Count(&total)
	dbQuery.Limit(param.Perpage).Offset(offset).Find(&data)

	totalPage := int((total + int64(param.Perpage) - 1) / int64(param.Perpage))

	return &ProductIndex{
		Data:      data,
		Total:     total,
		TotalPage: totalPage,
		Page:      param.Page,
		PerPage:   param.Perpage,
	}, nil
}

func (r *repository) Read(id string) (*Product, error) {
	var product Product
	check := r.db.Model(Product{}).
		Preload("App").
		Preload("Prices").
		Preload("Limits").
		Preload("Moduls.Modul").
		Preload("Moduls.Features.Feature").
		Where("id = ?", id).Find(&product)
	if check.RowsAffected == 0 {
		if check.Error != nil {
			return nil, check.Error
		}
		return nil, errors.New("id not found")
	}
	return &product, nil
}

func (r *repository) Delete(id string) error {
	delete := r.db.Where("id = ?", id).Delete(&Product{})
	if delete.Error != nil {
		return errors.New("data cannot be deleted")
	}

	return nil
}

func (r *repository) ensureUniqueSlug(name string, slug string, id ...string) (string, error) {
	var exists bool
	slug = generateSlug(name, slug)
	originalSlug := slug
	count := 0

	for {
		if id != nil {
			query := `SELECT EXISTS (SELECT 1 FROM products WHERE slug = ? AND id != ?)`
			if err := r.db.Raw(query, slug, id).Scan(&exists).Error; err != nil {
				return "", err
			}
		} else {
			query := `SELECT EXISTS (SELECT 1 FROM products WHERE slug = ?)`
			if err := r.db.Raw(query, slug).Scan(&exists).Error; err != nil {
				return "", err
			}
		}

		if !exists {
			break
		}

		count++
		slug = fmt.Sprintf("%s-%d", originalSlug, count)
	}

	return slug, nil
}

func generateSlug(name string, slug string) string {
	if slug == "" {
		slug = strings.ToLower(name)
	} else {
		slug = strings.ToLower(slug)
	}
	slug = strings.ToLower(slug)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	slug = re.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}

func (r *repository) CheckAppId(appId string) bool {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM apps WHERE id = ?)`
	if err := r.db.Raw(query, appId).Scan(&exists).Error; err != nil {
		return false
	}
	return exists
}
func (r *repository) CheckAppModulId(appModulId string) bool {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM app_moduls WHERE id = ?)`
	if err := r.db.Raw(query, appModulId).Scan(&exists).Error; err != nil {
		return false
	}
	return exists
}
func (r *repository) CheckAppModulFeatureId(id string, appModulId string) bool {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM app_modul_features WHERE id = ? and app_modul_id = ?)`
	if err := r.db.Raw(query, id, appModulId).Scan(&exists).Error; err != nil {
		return false
	}
	return exists
}

func (r *repository) CheckProductId(Id string) bool {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM products WHERE id = ?)`
	if err := r.db.Raw(query, Id).Scan(&exists).Error; err != nil {
		return false
	}
	return exists
}
func (r *repository) CheckProductModulId(Id string) bool {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM product_moduls WHERE id = ?)`
	if err := r.db.Raw(query, Id).Scan(&exists).Error; err != nil {
		return false
	}
	return exists
}

func (r *repository) PriceUpdate(req Request, id string) (res *Product, err error) {
	err = r.db.Transaction(func(tx *gorm.DB) error {
		tx.Where("product_id = ?", id).Delete(&ProductPrice{})
		for _, item := range req.Prices {
			timeNow := time.Now()
			itemId := uuid.NewString()
			record := ProductPrice{
				Id:             itemId,
				ProductId:      id,
				Type:           item.Type,
				Price:          item.Price,
				RecurringPrice: item.RecurringPrice,
				Qty:            item.Qty,
				CreatedAt:      timeNow,
				UpdatedAt:      timeNow,
			}
			if err := tx.Create(&record).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("failed to create price")
	}
	res, _ = r.Read(id)
	return res, nil
}

func (r *repository) LimitUpdate(req Request, id string) (res *Product, err error) {
	err = r.db.Transaction(func(tx *gorm.DB) error {
		tx.Where("product_id = ?", id).Delete(&ProductLimit{})
		for _, item := range req.Limits {
			timeNow := time.Now()
			itemId := uuid.NewString()
			record := ProductLimit{
				Id:        itemId,
				ProductId: id,
				Name:      item.Name,
				Key:       item.Key,
				Value:     item.Value,
				CreatedAt: timeNow,
				UpdatedAt: timeNow,
			}
			if err := tx.Create(&record).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("failed to create price")
	}
	res, _ = r.Read(id)
	return res, nil
}

func (r *repository) ModulUpdate(req Request, id string) (res *Product, err error) {
	err = r.db.Transaction(func(tx *gorm.DB) error {
		tx.Where("product_id = ?", id).Delete(&ProductModulFeature{})
		tx.Where("product_id = ?", id).Delete(&ProductModul{})
		for _, item := range req.Moduls {
			itemId := uuid.NewString()
			record := ProductModul{
				Id:         itemId,
				ProductId:  id,
				AppModulId: item.AppModulId,
			}
			if err := tx.Create(&record).Error; err != nil {
				return err
			}
			for _, item2 := range item.Features {
				item2Id := uuid.NewString()
				record := ProductModulFeature{
					Id:                item2Id,
					ProductId:         id,
					ProductModulId:    itemId,
					AppModulFeatureId: item2.AppModulFeatureId,
				}
				if err := tx.Create(&record).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("failed to create modul")
	}
	res, _ = r.Read(id)
	return res, nil
}

package promotion

import (
	"errors"
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

func (r *repository) Create(req Request) (*Promotion, error) {
	Id := uuid.NewString()
	timeNow := time.Now()
	err := r.db.Transaction(func(tx *gorm.DB) error {
		data := Promotion{
			Id:          Id,
			Code:        req.Code,
			Name:        req.Name,
			Description: req.Description,
			TargetUser:  req.TargetUser,
			StartDate:   req.StartDate,
			EndDate:     req.EndDate,
			Status:      req.Status,
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
		}
		if err := tx.Create(&data).Error; err != nil {
			return err
		}

		for _, item := range req.Items {
			itemId := uuid.NewString()
			itemRecord := PromotionItem{
				Id:          itemId,
				PromotionId: Id,
				ProductType: item.ProductType,
				ProductId:   item.ProductId,
				PriceType:   item.PriceType,
				Type:        item.Type,
				Value:       item.Value,
				CreatedAt:   timeNow,
				UpdatedAt:   timeNow,
			}
			if err := tx.Create(&itemRecord).Error; err != nil {
				return err
			}
		}
		return nil // commit jika semua berhasil
	})

	if err != nil {
		return nil, errors.New("failed to create")
	}
	res, _ := r.Read(Id)
	return res, nil
}

func (r *repository) Update(req Request, id string) (res *Promotion, err error) {
	err = r.db.Transaction(func(tx *gorm.DB) error {
		timeNow := time.Now()
		updateData := map[string]interface{}{
			"code":        req.Code,
			"name":        req.Name,
			"description": req.Description,
			"target_user": req.TargetUser,
			"start_date":  req.StartDate,
			"end_date":    req.EndDate,
			"updated_at":  time.Now(),
		}
		result := tx.Model(&Promotion{}).Where("id = ?", id).Updates(updateData)
		if result.Error != nil {
			return errors.New("update failed")
		}
		tx.Where("promotion_id = ?", id).Delete(&PromotionItem{})
		for _, item := range req.Items {
			itemId := uuid.NewString()
			itemRecord := PromotionItem{
				Id:          itemId,
				PromotionId: id,
				ProductType: item.ProductType,
				ProductId:   item.ProductId,
				PriceType:   item.PriceType,
				Type:        item.Type,
				Value:       item.Value,
				UpdatedAt:   timeNow,
			}
			if err := tx.Create(&itemRecord).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("failed to update")
	}

	res, _ = r.Read(id)
	return res, nil
}

func (r *repository) UpdateStatus(req Request, id string) (res *Promotion, err error) {
	result := r.db.Model(&Promotion{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     req.Status,
		"updated_at": time.Now(),
	})

	if result.Error != nil {
		return nil, result.Error
	}
	res, _ = r.Read(id)
	return res, nil
}

func (r *repository) Index(param ParamIndex) (*PromotionIndex, error) {
	var data []Promotion
	var total int64

	if param.Page < 1 {
		param.Page = 1
	}
	if param.Perpage < 1 {
		param.Perpage = 20
	}

	offset := (param.Page - 1) * param.Perpage
	dbQuery := r.db.Model(&Promotion{}).Preload("Items").Order("created_at DESC")
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

	return &PromotionIndex{
		Data:      data,
		Total:     total,
		TotalPage: totalPage,
		Page:      param.Page,
		PerPage:   param.Perpage,
	}, nil
}

func (r *repository) Read(id string) (*Promotion, error) {
	var data Promotion
	check := r.db.Model(Promotion{}).Where("id = ?", id).Find(&data)
	if check.RowsAffected == 0 {
		if check.Error != nil {
			return nil, check.Error
		}
		return nil, errors.New("id not found")
	}
	return &data, nil
}
func (r *repository) ReadDetail(id string) (*Promotion, error) {
	var data Promotion
	check := r.db.Model(Promotion{}).Where("id = ?", id).Preload("Items").Find(&data)
	if check.RowsAffected == 0 {
		if check.Error != nil {
			return nil, check.Error
		}
		return nil, errors.New("id not found")
	}
	for i, item := range data.Items {
		if item.ProductType == "product" {
			var myItem Product
			r.db.Model(Product{}).Where("id = ?", item.ProductId).
				Select("id", "name", "slug", "description", "status").
				Find(&myItem)
			data.Items[i].Product = myItem
		}
		if item.ProductType == "plan" {
			var myItem Product
			r.db.Model(Product{}).Where("id = ?", item.ProductId).
				Preload("Plans", "id = ?", item.ProductPlanId).
				Find(&myItem)
			data.Items[i].Product = myItem
		}
		if item.ProductType == "price" {
			var myItem Product
			r.db.Model(Product{}).Where("id = ?", item.ProductId).
				Preload("Plans", "id = ?", item.ProductPlanId).
				Preload("Plans.Prices", "type = ?", item.PriceType).
				Find(&myItem)
			data.Items[i].Product = myItem
		}
	}
	return &data, nil
}

func (r *repository) Delete(id string) error {

	err := r.db.Transaction(func(tx *gorm.DB) error {
		delete := r.db.Where("promotion_id = ?", id).Delete(&PromotionItem{})
		if delete.Error != nil {
			return errors.New("item cannot be deleted")
		}
		delete = r.db.Where("id = ?", id).Delete(&Promotion{})
		if delete.Error != nil {
			return errors.New("data cannot be deleted")
		}
		return nil // commit jika semua berhasil
	})

	if err != nil {
		return errors.New("failed to delete")
	}
	return nil
}

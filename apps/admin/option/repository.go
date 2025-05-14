package option

import (
	"encoding/json"

	"github.com/bonsus/go-saas/internal/utils/option"
	"gorm.io/gorm"
)

type OptionRepository interface {
	Save(Option Option) (*Option, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repository {
	return repository{
		db: db,
	}
}

func (r *repository) Save(option Option) (*Option, error) {
	check := r.db.Select("name").Where("name = ?", option.Name).Find(&option)
	if check.RowsAffected == 0 {
		create := r.db.Create(&option)
		if create.Error != nil {
			return nil, create.Error
		}
	}
	update := r.db.Model(&option).Where("name = ?", option.Name).Select("value").Updates(option)
	if update.Error != nil {
		return nil, update.Error
	}
	return &option, nil
}

func (r *repository) Get(name string) (json.RawMessage, error) {
	res, err := option.GetOption(r.db, name)
	if err != nil {
		return nil, err
	}
	return res, nil
}

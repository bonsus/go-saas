package option

import (
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type Option struct {
	ID    uint            `json:"id"`
	Name  string          `json:"name"`
	Value json.RawMessage `json:"value"`
}

// GetOption - Ambil option berdasarkan name
func GetOption(db *gorm.DB, name string) (json.RawMessage, error) {
	var option Option
	result := db.Select("id,name,value").Where("name = ?", name).Find(&option)
	if result.RowsAffected == 0 {
		return nil, errors.New("option not found")
	}
	return option.Value, nil
}

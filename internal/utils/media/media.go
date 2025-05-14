package media

import (
	"encoding/json"
	"time"

	"github.com/bonsus/go-saas/internal/config"
	"gorm.io/gorm"
)

type Medias struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Alt         string      `json:"alt"`
	Type        string      `json:"type"`
	Status      string      `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Files       []MediaFile `json:"files" gorm:"foreignKey:MediaId"`
}

type MediaFile struct {
	Id       string `json:"id"`
	MediaId  string `json:"media_id"`
	File     string `json:"file"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Filesize string `json:"filesize"`
	Url      string `json:"url" gorm:"-"`
}

func GetMedia(db *gorm.DB, Id string) json.RawMessage {
	var media Medias
	var cfg = config.GetConfig()
	err := db.Table("medias").Preload("Files").Where("id =?", Id).First(&media).Error
	if err != nil || media.Id == "" {
		return nil
	}
	for i, f := range media.Files {
		media.Files[i].Url = cfg.S3.Domain + f.File
	}
	data, _ := json.Marshal(media)
	return data
}

func GetMedias(db *gorm.DB, Ids json.RawMessage) json.RawMessage {
	var media []Medias
	var cfg = config.GetConfig()
	var ids []string
	if err := json.Unmarshal(Ids, &ids); err != nil {
		return nil
	}
	err := db.Table("medias").Preload("Files").Where("id IN (?)", ids).Find(&media).Error
	if err != nil {
		return nil
	}
	for _, m := range media {
		for i, f := range m.Files {
			m.Files[i].Url = cfg.S3.Domain + f.File
		}
	}
	data, _ := json.Marshal(media)
	return data
}

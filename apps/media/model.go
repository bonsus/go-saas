package media

import (
	"time"
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

type MediaIndex struct {
	Data      []Medias `json:"data"`
	Page      int      `json:"page"`
	PerPage   int      `json:"per_page"`
	TotalPage int      `json:"total_page"`
	Total     int64    `json:"total"`
}

package product

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type App struct {
	Id          string    `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Url         string    `json:"url"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt
}
type AppModul struct {
	Id          string    `json:"id"`
	AppId       string    `json:"app_id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt
}
type AppModulFeature struct {
	Id          string          `json:"id"`
	AppId       string          `json:"app_id"`
	AppModulId  string          `json:"app_modul_id"`
	Code        string          `json:"code"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Permission  json.RawMessage `json:"permission"`
	Status      string          `json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index"`
}

type Product struct {
	Id          string    `json:"id"`
	AppId       string    `json:"app_id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug" gorm:"unique"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt

	App    App            `json:"app" gorm:"foreignKey:AppId"`
	Moduls []ProductModul `json:"moduls" gorm:"foreignKey:ProductId"`
	Prices []ProductPrice `json:"prices" gorm:"foreignKey:ProductId"`
	Limits []ProductLimit `json:"limits" gorm:"foreignKey:ProductId"`
}
type ProductModul struct {
	Id         string                `json:"Id"`
	ProductId  string                `json:"product_id"`
	AppModulId string                `json:"app_modul_id"`
	Modul      AppModul              `json:"modul" gorm:"foreignKey:AppModulId"`
	Features   []ProductModulFeature `json:"features" gorm:"foreignKey:ProductModulId"`
}
type ProductModulFeature struct {
	Id                string          `json:"Id"`
	ProductId         string          `json:"product_id"`
	ProductModulId    string          `json:"product_modul_id"`
	AppModulFeatureId string          `json:"app_modul_feature_id"`
	Feature           AppModulFeature `json:"feature" gorm:"foreignKey:AppModulFeatureId"`
}

type ProductLimit struct {
	Id        string    `json:"id" gorm:"primaryKey"`
	ProductId string    `json:"product_id"`
	Name      string    `json:"name"`
	Key       string    `json:"key"`
	Value     int       `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type ProductPrice struct {
	Id             string    `json:"id"`
	ProductId      string    `json:"product_id"`
	Type           string    `json:"type"`
	Price          int64     `json:"price"`
	RecurringPrice int64     `json:"recurring_price"`
	Qty            int       `json:"qty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ProductIndex struct {
	Data      []Product `json:"data"`
	Page      int       `json:"page"`
	PerPage   int       `json:"per_page"`
	TotalPage int       `json:"total_page"`
	Total     int64     `json:"total"`
}

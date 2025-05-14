package migration

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	Id          string         `json:"id" gorm:"primaryKey"`
	AppId       string         `json:"app_id"`
	Name        string         `json:"name"`
	Slug        string         `json:"slug" gorm:"unique"`
	Description string         `json:"description"`
	Image       string         `json:"image"`
	Status      string         `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	App App `gorm:"foreignKey:AppId"`
}
type ProductModul struct {
	Id         string `json:"id" gorm:"primaryKey"`
	ProductId  string `json:"product_id"`
	AppModulId string `json:"app_modul_id"`

	Product  Product  `gorm:"foreignKey:ProductId"`
	AppModul AppModul `gorm:"foreignKey:AppModulId"`
}
type ProductModulFeature struct {
	Id                string `json:"id" gorm:"primaryKey"`
	ProductId         string `json:"product_id"`
	ProductModulId    string `json:"product_modul_id"`
	AppModulFeatureId string `json:"app_modul_feature_id"`

	Product         Product         `gorm:"foreignKey:ProductId"`
	ProductModul    ProductModul    `gorm:"foreignKey:ProductModulId"`
	AppModulFeature AppModulFeature `gorm:"foreignKey:AppModulFeatureId"`
}
type ProductLimit struct {
	Id        string    `json:"id" gorm:"primaryKey"`
	ProductId string    `json:"product_id"`
	Name      string    `json:"name"`
	Key       string    `json:"key"`
	Value     int       `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Product Product `gorm:"foreignKey:ProductId"`
}
type ProductPrice struct {
	Id             string    `json:"id" gorm:"primaryKey"`
	ProductId      string    `json:"product_id"`
	Type           string    `json:"type"`
	Price          int64     `json:"price"`
	RecurringPrice int64     `json:"recurring_price"`
	Qty            int       `json:"qty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Product Product `gorm:"foreignKey:ProductId"`
}

type Promotion struct {
	Id          string `json:"id" gorm:"primaryKey"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TargetUser  string `json:"target_user"`
	StartDate   time.Time
	EndDate     time.Time
	Status      string `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type PromotionItem struct {
	Id          string    `json:"id" gorm:"primaryKey"`
	PromotionId string    `json:"promotion_id"`
	ProductId   string    `json:"product_id"`
	PriceType   string    `json:"price_type"`
	Type        string    `json:"type"`
	Value       float64   `json:"value"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Promotion Promotion `gorm:"foreignKey:PromotionId"`
}

package promotion

import (
	"time"
)

type Promotion struct {
	Id          string    `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	TargetUser  string    `json:"target_user"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Items []PromotionItem `json:"items" gorm:"foreignKey:PromotionId"`
}
type PromotionItem struct {
	Id            string    `json:"id"`
	PromotionId   string    `json:"promotion_id"`
	ProductType   string    `json:"product_type"` // all, product, plan, dll
	ProductId     string    `json:"product_id"`
	ProductPlanId string    `json:"product_plan_id"`
	PriceType     string    `json:"price_type"`
	Type          string    `json:"type"`
	Value         float64   `json:"value"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Product       Product   `json:"product" gorm:"-"`
}
type Product struct {
	Id          string        `json:"id"`
	Name        string        `json:"name"`
	Slug        string        `json:"slug"`
	Description string        `json:"description"`
	Status      string        `json:"status"`
	Plans       []ProductPlan `json:"plans" gorm:"foreignKey:ProductId"`
}
type ProductPlan struct {
	Id          string             `json:"id"`
	ProductId   string             `json:"product_id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Status      string             `json:"status"`
	Prices      []ProductPlanPrice `json:"prices" gorm:"foreignKey:ProductPlanId"`
}
type ProductPlanPrice struct {
	ProductPlanId string `json:"product_plan_id"`
	Type          string `json:"type"`
	Price         int64  `json:"price"`
}
type PromotionIndex struct {
	Data      []Promotion `json:"data"`
	Page      int         `json:"page"`
	PerPage   int         `json:"per_page"`
	TotalPage int         `json:"total_page"`
	Total     int64       `json:"total"`
}

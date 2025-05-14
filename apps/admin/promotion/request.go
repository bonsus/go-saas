package promotion

import "time"

type Request struct {
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	TargetUser  string    `json:"target_user"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Status      string    `json:"status"`

	Items []RequestItem `json:"items"`
}
type RequestItem struct {
	Id          string  `json:"id"`
	PromotionId string  `json:"promotion_id"`
	ProductType string  `json:"product_type"`
	ProductId   string  `json:"product_id"`
	PriceType   string  `json:"price_type"`
	Type        string  `json:"type"`
	Value       float64 `json:"value"`
}

type ParamIndex struct {
	Page    int    `query:"page"`
	Perpage int    `query:"perpage"`
	Search  string `query:"search"`
}

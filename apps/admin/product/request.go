package product

type Request struct {
	AppId       string         `json:"app_id"`
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Description string         `json:"description"`
	Image       string         `json:"image"`
	Status      string         `json:"status"`
	Moduls      []ModulRequest `json:"moduls"`
	Limits      []LimitRequest `json:"limits"`
	Prices      []PriceRequest `json:"prices"`
}
type ModulRequest struct {
	ProductId  string           `json:"product_id"`
	AppModulId string           `json:"app_modul_id"`
	Features   []FeatureRequest `json:"features"`
}
type FeatureRequest struct {
	ProductModulId    string `json:"product_modul_id"`
	AppModulFeatureId string `json:"app_modul_feature_id"`
}
type LimitRequest struct {
	Id        string `json:"id"`
	ProductId string `json:"product_id"`
	Name      string `json:"name"`
	Key       string `json:"key"`
	Value     int    `json:"value"`
}
type PriceRequest struct {
	Id             string `json:"id"`
	ProductId      string `json:"product_id"`
	Type           string `json:"type"`
	Price          int64  `json:"price"`
	RecurringPrice int64  `json:"recurring_price"`
	Qty            int    `json:"qty"`
}

type ParamIndex struct {
	Page    int    `query:"page"`
	Perpage int    `query:"perpage"`
	Search  string `query:"search"`
}

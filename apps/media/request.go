package media

type updateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Alt         string `json:"alt"`
	Status      string `json:"status"`
}

type ParamIndex struct {
	Page    int    `query:"page"`
	Perpage int    `query:"perpage"`
	Search  string `query:"search"`
	Status  string `query:"status"`
}

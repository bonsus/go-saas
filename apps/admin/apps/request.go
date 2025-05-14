package apps

import "encoding/json"

type Request struct {
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Image       string   `json:"image"`
	Url         string   `json:"url"`
	BackendUrl  string   `json:"backend_url"`
	Server      string   `json:"server"`
	Db          string   `json:"db"`
	DbHost      string   `json:"db_host"`
	DbPort      string   `json:"db_port"`
	DbUser      string   `json:"db_user"`
	DbPass      string   `json:"db_pass"`
	DbName      string   `json:"db_name"`
	Status      string   `json:"status"`
	Ids         []string `json:"ids"`
}
type ParamIndex struct {
	Page    int    `query:"page"`
	Perpage int    `query:"perpage"`
	Search  string `query:"search"`
}

type ModulRequest struct {
	AppId       string           `json:"app_id"`
	Code        string           `json:"code"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Status      string           `json:"status"`
	Features    []FeatureRequest `json:"features"`
}
type FeatureRequest struct {
	AppId       string          `json:"app_id"`
	AppModulId  string          `json:"app_modul_id"`
	Code        string          `json:"code"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Permission  json.RawMessage `json:"permission"`
	Status      string          `json:"status"`
}
type PluginRequest struct {
	AppId       string `json:"app_id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Status      string `json:"status"`
}

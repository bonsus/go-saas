package apps

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
	BackendUrl  string    `json:"backend_url"`
	DbHost      string    `json:"db_host"`
	DbPort      string    `json:"db_port"`
	DbUser      string    `json:"db_user"`
	DbPass      string    `json:"db_pass"`
	DbName      string    `json:"db_name"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Moduls   []AppModul        `json:"moduls" gorm:"foreignKey:AppId"`
	Features []AppModulFeature `json:"features" gorm:"foreignKey:AppId"`
	Plugins  []AppPlugin       `json:"plugins"`
}

type AppData struct {
	DbHost string `json:"db_host"`
	DbPort string `json:"db_port"`
	DbUser string `json:"db_user"`
	DbPass string `json:"db_pass"`
	DbName string `json:"db_name"`
}

type AppModul struct {
	Id          string            `json:"id"`
	AppId       string            `json:"app_id"`
	Code        string            `json:"code"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Status      string            `json:"status"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	DeletedAt   gorm.DeletedAt    `gorm:"index"`
	Features    []AppModulFeature `json:"features" gorm:"foreignKey:AppModulId"`
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
type AppPlugin struct {
	Id        string    `json:"id"`
	AppId     string    `json:"app_id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AppIndex struct {
	Data      []App `json:"data"`
	Page      int   `json:"page"`
	PerPage   int   `json:"per_page"`
	TotalPage int   `json:"total_page"`
	Total     int64 `json:"total"`
}

type AppModulIndex struct {
	Data      []AppModul `json:"data"`
	Page      int        `json:"page"`
	PerPage   int        `json:"per_page"`
	TotalPage int        `json:"total_page"`
	Total     int64      `json:"total"`
}

type AppPluginIndex struct {
	Data      []AppPlugin `json:"data"`
	Page      int         `json:"page"`
	PerPage   int         `json:"per_page"`
	TotalPage int         `json:"total_page"`
	Total     int64       `json:"total"`
}

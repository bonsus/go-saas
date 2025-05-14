package migration

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type App struct {
	Id          string         `json:"id" gorm:"primaryKey"`
	Code        string         `json:"code"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Image       string         `json:"image"`
	Url         string         `json:"url"`
	BackendUrl  string         `json:"backend_url"`
	DbHost      string         `json:"db_host"`
	DbPort      string         `json:"db_port"`
	DbUser      string         `json:"db_user"`
	DbPass      string         `json:"db_pass"`
	DbName      string         `json:"db_name"`
	Status      string         `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
type AppModul struct {
	Id          string         `json:"id" gorm:"primaryKey"`
	AppId       string         `json:"app_id"`
	Code        string         `json:"code"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	App App `gorm:"foreignKey:AppId"`
}
type AppModulFeature struct {
	Id          string          `json:"id" gorm:"primaryKey"`
	AppId       string          `json:"app_id"`
	AppModulId  string          `json:"app_modul_id"`
	Code        string          `json:"code"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Permission  json.RawMessage `json:"permission" gorm:"type:jsonb"`
	Status      string          `json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index"`

	App      App      `gorm:"foreignKey:AppId"`
	AppModul AppModul `gorm:"foreignKey:AppModulId"`
}
type AppPlugin struct {
	Id          string         `json:"id" gorm:"primaryKey"`
	AppId       string         `json:"app_id"`
	Code        string         `json:"code"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Image       string         `json:"image"`
	Status      string         `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	App App `gorm:"foreignKey:AppId"`
}

type UserApp struct {
	Id         string `json:"id" gorm:"primaryKey;unique;not null;size:36"`
	Code       string `json:"code" gorm:"unique;not null"`
	UserId     string `json:"user_id"`
	AppId      string `json:"app_id"`
	Url        string `json:"url"`
	BackendUrl string `json:"backend_url"`
	Server     string `json:"server"`
	Db         string `json:"db"`
	DbHost     string `json:"db_host"`
	DbPort     string `json:"db_port"`
	DbUser     string `json:"db_user"`
	DbPass     string `json:"db_pass"`
	DbName     string `json:"db_name"`
	Status     string `json:"status"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`

	User User `gorm:"foreignKey:UserId"`
	App  App  `gorm:"foreignKey:AppId"`
}
type UserAppPlugin struct {
	Id          string `json:"id" gorm:"primaryKey;unique;not null;size:36"`
	UserAppId   string `json:"user_app_id"`
	AppId       string `json:"app_id"`
	AppPluginId string `json:"app_plugin_id"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	UserApp   UserApp   `gorm:"foreignKey:UserAppId"`
	App       App       `gorm:"foreignKey:AppId"`
	AppPlugin AppPlugin `gorm:"foreignKey:AppPluginId"`
}

type User struct {
	Id          string `json:"id" gorm:"primaryKey;unique;not null;size:36"`
	Company     string `json:"company"`
	Name        string `json:"name" gorm:"not null;size:100"`
	Phone       string `json:"phone"`
	Username    string `json:"username" gorm:"not null;unique;size:100;index"`
	Email       string `json:"email" gorm:"not null;unique;size:100;index"`
	Password    string `json:"password"`
	Type        string `json:"type"`
	ParentId    string `json:"parent_id" gorm:"nullable"`
	Status      string `json:"status" gorm:"default:'inactive';size:20"`
	EmailStatus string `json:"email_status"`
	PhoneStatus string `json:"phone_status"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
type UserPermission struct {
	Id         string          `json:"id" gorm:"primaryKey;unique;not null;size:36"`
	UserId     string          `json:"user_id"`
	AppId      string          `json:"app_id"`
	Permission json.RawMessage `json:"permission"`

	User User `gorm:"foreignKey:UserId"`
	App  App  `gorm:"foreignKey:AppId"`
}

// type UserRole struct {
// 	Id         string          `json:"id" gorm:"primaryKey;unique;not null;size:36"`
// 	UserAppId  string          `json:"user_app_id"`
// 	Name       string          `json:"name" gorm:"not null;size:50;index"`
// 	Permission json.RawMessage `json:"permission" gorm:"type:jsonb"`
// 	CreatedAt  time.Time
// 	UpdatedAt  time.Time
// 	DeletedAt  gorm.DeletedAt `gorm:"index"`
// }

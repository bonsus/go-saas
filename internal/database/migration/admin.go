package migration

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	Id        string `json:"id" gorm:"primaryKey;unique;not null;size:36"`
	Name      string `json:"name" gorm:"not null;size:100"`
	Username  string `json:"username" gorm:"not null;unique;size:100;index"`
	Email     string `json:"email" gorm:"not null;unique;size:100;index"`
	Password  string `json:"password"`
	RoleId    string `json:"role_id" gorm:"nullable"`
	Status    string `json:"status" gorm:"default:'inactive';size:20"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Role AdminRole `gorm:"foreignKey:RoleId"`
}

type AdminRole struct {
	Id         string          `json:"id" gorm:"primaryKey;unique;not null;size:36"`
	Name       string          `json:"name" gorm:"not null;size:50;index"`
	Permission json.RawMessage `json:"permission" gorm:"type:jsonb"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

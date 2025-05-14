package migration

import (
	"encoding/json"
	"time"
)

type Option struct {
	Id        uint            `json:"id" gorm:"primarykey"`
	Name      string          `json:"name" gorm:"not null;uniqueIndex"`
	Value     json.RawMessage `json:"value" gorm:"type:jsonb"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

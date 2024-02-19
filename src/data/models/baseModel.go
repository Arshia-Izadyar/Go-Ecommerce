package models

import "time"

type BaseModel struct {
	Id        int       `gorm: "primaryKey"`
	CreatedAt time.Time `gorm:"type:TIMESTAMP with time zone;not null"`
	UpdatedAt time.Time `gorm:"type:TIMESTAMP with time zone;not null`
}

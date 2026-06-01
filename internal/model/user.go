package model

import (
	"time"
)

type User struct {
	ID           uint    `gorm:"primaryKey"`
	TenantID     *string `gorm:"type:varchar(36);index"`
	Username     string  `gorm:"type:varchar(255);unique;not null"`
	Password     string  `gorm:"type:text;not null"`
	SessionToken string  `gorm:"type:varchar(255)"`
	CSRFToken    string  `gorm:"type:varchar(255)"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

package model

import "gorm.io/gorm"

type UserLogin struct {
	gorm.Model
	Tij         string `gorm:"uniqueIndex; not null; type:varchar(255)"`
	RefeshToken string `gorm:"uniqueIndex; not null; type:varchar(255)"`
	ExpiresAt   int    `gorm:"not null"`
	IsRevoked   bool   `gorm:"not null"`
	UserName    string `gorm:"not null"`
}

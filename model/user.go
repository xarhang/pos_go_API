package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex; not null; type:varchar(255)"`
	Password string `gorm:"not null; type:varchar(255)"`
	Fullname string `gorm:"not null;"`
	Avatar   string
	StatusID uint
	RuleID   uint
	Status   Status
	Rule     Rule
}

package model

import "gorm.io/gorm"

type Status struct {
	gorm.Model
	StatusName string `gorm:"type:varchar(255)"`
	User       []User `gorm:"foreignKey:StatusID"`
}

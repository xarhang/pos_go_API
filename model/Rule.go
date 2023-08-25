package model

import "gorm.io/gorm"

type Rule struct {
	gorm.Model
	RuleName string `gorm:"type:varchar(255);not null"`
	User     []User `gorm:"foriegnKey:RuleID"`
}

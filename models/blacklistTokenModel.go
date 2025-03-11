package models

import "gorm.io/gorm"

type BlacklistToken struct {
	gorm.Model
	Token string `gorm:"unique"`
}

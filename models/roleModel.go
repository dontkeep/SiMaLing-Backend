package models

import "gorm.io/gorm"

type Roles struct {
	gorm.Model
	Id        uint `gorm:"primaryKey"`
	Role_Name string
}

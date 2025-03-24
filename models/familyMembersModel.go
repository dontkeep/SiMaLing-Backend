package models

import "gorm.io/gorm"

type FamilyMembers struct {
	gorm.Model
	ID              uint `gorm:"primaryKey"`
	NIK             string
	Name            string
	Status          string
	HeadOfFamily    User `gorm:"foreignKey:HeadOfFamily_Id"`
	HeadOfFamily_Id uint
}

// status only has 2 values, "wife" and "child"

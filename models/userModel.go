package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey"`
	Phone_No string `gorm:"unique"`
	NIK      string `gorm:"unique"`
	Name     string
	Address  string
	Funds    string
	Role     Roles `gorm:"foreignKey:Role_Id"`
	Role_Id  uint
}

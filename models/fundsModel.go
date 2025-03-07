package models

import "gorm.io/gorm"

type Funds struct {
	gorm.Model
	ID        uint `gorm:"primaryKey"`
	User      User `gorm:"foreignKey:User_Id"`
	User_Id   uint
	Amount    float64
	Is_Income bool
}

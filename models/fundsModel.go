package models

import "gorm.io/gorm"

type Funds struct {
	gorm.Model
	ID          uint `gorm:"primaryKey"`
	User        User `gorm:"foreignKey:User_Id"`
	User_Id     uint
	Block       string
	Amount      float64
	Image       string
	Description string
	Is_Income   bool
	Status      string
}

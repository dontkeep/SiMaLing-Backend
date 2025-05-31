package models

import "gorm.io/gorm"

type SecurityRecord struct {
	gorm.Model
	Security_ID User `gorm:"foreignKey:Security_Id"`
	Security_Id uint
	Block       string
	Phone_No    string
	Longitude   string
	Latitude    string
}

package models

import "gorm.io/gorm"

type SecurityRecord struct {
	gorm.Model
	Security_ID User `gorm:"foreignKey:Security_Id"`
	Security_Id uint
	Longitude   string
	Latitude    string
}

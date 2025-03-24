package models

import "gorm.io/gorm"

type SecurityRecord struct {
	gorm.Model
	Security_ID User `gorm:"foreignKey:Security_Id"`
	Security_Id uint
	NIK         string
	Longitude   string
	Latitude    string
}

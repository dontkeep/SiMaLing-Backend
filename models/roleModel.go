package models

import "gorm.io/gorm"

type Roles struct {
	gorm.Model
	Role_Id   uint `gorm:"primaryKey"`
	Role_Name string
}

// Roles struct is a model that represents the roles table in the database. It has the following fields:
// Role_Id: the primary key of the roles table.
// Role_Name: the name of the role.

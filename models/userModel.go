package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey"`
	Phone_No string `gorm:"unique"`
	NIK      string `gorm:"unique"`
	Name     string
	Address  string
	Role     Roles `gorm:"foreignKey:Role_Id"`
	Role_Id  uint
}

// User struct is a model that represents the users table in the database. It has the following fields:
// ID: the primary key of the users table.
// Phone_No: the phone number of the user.
// NIK: the national identification number of the user.
// Name: the name of the user.
// Address: the address of the user.
// Funds: the funds of the user.
// Role: the role of the user. It is a foreign key that references the Role_Id field in the roles table.
// Role_Id: the foreign key that references the Role_Id field in the roles table.
// The User struct also embeds the gorm.Model struct, which provides the fields ID, CreatedAt, UpdatedAt, and DeletedAt to the User struct. These fields are used to manage the lifecycle of the user records in the database.
// The User struct is used to interact with the users table in the database. It defines the structure of the user records and the relationships between the user records and other tables in the database. The User struct is defined in the models/userModel.go file.

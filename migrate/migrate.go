package main

import (
	"github.com/dontkeep/simaling-backend/initializers"
	"github.com/dontkeep/simaling-backend/models"
)

func init() {
	initializers.LoadEnvVar()
	initializers.DatabaseConnection()
}

func main() {
	initializers.DB.AutoMigrate(&models.User{}, &models.Roles{}, &models.Funds{}, &models.BlacklistToken{}, &models.SecurityRecord{})
}

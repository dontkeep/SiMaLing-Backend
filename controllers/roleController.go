package controllers

import (
	"github.com/dontkeep/simaling-backend/initializers"
	"github.com/dontkeep/simaling-backend/models"
	"github.com/gin-gonic/gin"
)

func GetAllRoles(c *gin.Context) {
	var roles []models.Roles
	result := initializers.DB.Find(&roles)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to get roles",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": roles,
	})
}

func CreateDefaultRoles() error {
	var roles = []models.Roles{
		{
			Role_Name: "Admin",
		},
		{
			Role_Name: "User",
		},
		{
			Role_Name: "Security",
		},
	}

	for _, role := range roles {
		var existingRole models.Roles
		result := initializers.DB.Where("role_name = ?", role.Role_Name).First(&existingRole)

		if result.Error == nil {
			continue
		}

		if result := initializers.DB.Create(&role); result.Error != nil {
			return result.Error
		}
	}
	return nil
}

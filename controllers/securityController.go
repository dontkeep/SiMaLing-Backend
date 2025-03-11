package controllers

import (
	"github.com/dontkeep/simaling-backend/initializers"
	"github.com/dontkeep/simaling-backend/models"
	"github.com/gin-gonic/gin"
)

func GetAllSecurityRecord(c *gin.Context) {
	var securityRecord []models.SecurityRecord

	result := initializers.DB.Find(&securityRecord)

	if result != nil {
		c.JSON(400, gin.H{
			"message": "Failed to get security record",
		})
		return
	}
}

func CreateSecurityRecord(c *gin.Context) {

}

package controllers

import (
	"github.com/dontkeep/simaling-backend/initializers"
	"github.com/dontkeep/simaling-backend/models"
	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) {
	var users []models.User
	result := initializers.DB.Find(&users)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to get users",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": users,
	})
}

func GetUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	result := initializers.DB.First(&user, id)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to get user",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": user,
	})
}

func CreateUser(c *gin.Context) {
	var body struct {
		Phone_No string
		NIK      string
		Name     string
		Address  string
		Funds    string
		Role_Id  uint
	}

	c.BindJSON(&body)

	user := models.User{
		Phone_No: body.Phone_No,
		NIK:      body.NIK,
		Name:     body.Name,
		Address:  body.Address,
		Funds:    body.Funds,
		Role_Id:  body.Role_Id,
	}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to create user",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "User has been created",
	})
}

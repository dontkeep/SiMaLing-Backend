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
		Password string
		Name     string
		Address  string
		Funds    string
		Role_Id  uint
	}

	c.BindJSON(&body)

	user := models.User{
		Phone_No: body.Phone_No,
		NIK:      body.NIK,
		Password: body.Password,
		Name:     body.Name,
		Address:  body.Address,
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

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	var body struct {
		Phone_No string
		NIK      string
		Password string
		Name     string
		Address  string
		Role_Id  uint
	}

	c.BindJSON(&body)

	user.Phone_No = body.Phone_No
	user.NIK = body.NIK
	user.Password = body.Password
	user.Name = body.Name
	user.Address = body.Address
	user.Role_Id = body.Role_Id

	result := initializers.DB.Update(id, &user)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to update user",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "User has been updated",
	})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	result := initializers.DB.Delete(&user, id)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to delete user",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "User has been deleted",
	})
}

package controllers

import (
	"github.com/dontkeep/simaling-backend/initializers"
	"github.com/dontkeep/simaling-backend/models"
	"github.com/gin-gonic/gin"
)

func PostCreate(c *gin.Context) {
	//Get Data from the request
	var body struct {
		Body  string
		Title string
	}

	c.BindJSON(&body)

	//Save to the database
	post := models.Post{
		Title: body.Title,
		Body:  body.Body,
	}

	result := initializers.DB.Create(&post)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to create post",
		})
		return
	}

	//Return the response
	c.JSON(200, gin.H{
		"message": "Post has been created",
		"data": gin.H{
			"id":    post.ID,
			"title": post.Title,
			"body":  post.Body,
		},
	})
}

func GetAllPost(c *gin.Context) {
	var posts []models.Post
	result := initializers.DB.Find(&posts)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to get posts",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": posts,
	})
}

func GetPost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post

	result := initializers.DB.Find(&post, id)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to get post",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": post,
	})
}

func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post

	//Update post data
	var body struct {
		Body  string
		Title string
	}

	c.BindJSON(&body)

	post.Title = body.Title
	post.Body = body.Body

	result := initializers.DB.Update(id, &post)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to update post",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Post has been updated",
		"data": gin.H{
			"id":    post.ID,
			"title": post.Title,
			"body":  post.Body,
		},
	})
}

func DeletePost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post

	result := initializers.DB.Delete(&post, id)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to delete post",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Post has been deleted",
	})
}

package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"service/config"
	"service/models"
)

func GetUsers(c *gin.Context) {

	var users []models.User

	config.DB.Find(&users)

	c.JSON(http.StatusOK, users)
}

func CreateUser(c *gin.Context) {

	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	config.DB.Create(&user)

	c.JSON(200, user)
}

func GetUserByID(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(404, gin.H{
			"message": "User not found",
		})
		return
	}

	c.JSON(200, user)
}

func UpdateUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
		})
		return
	}

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	config.DB.Save(&user)

	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
		})
		return
	}

	config.DB.Delete(&user)

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}
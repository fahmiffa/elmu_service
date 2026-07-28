package controllers

import (
	"net/http"
	"service/config"
	"service/models"

	"github.com/gin-gonic/gin"
)

// GetUnits mengambil seluruh data Unit (cabang)
func GetUnits(c *gin.Context) {
	var units []models.Unit
	if err := config.DB.Find(&units).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data unit: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": units,
	})
}

// GetPrograms mengambil seluruh data Program (paket)
func GetPrograms(c *gin.Context) {
	var programs []models.Program
	if err := config.DB.Find(&programs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data program: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": programs,
	})
}

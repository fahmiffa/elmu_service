package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// GetAppStatus returns the current application status and mode
func GetAppStatus(c *gin.Context) {
	mode := os.Getenv("APP_MODE")
	if mode == "" {
		mode = "development"
	}

	version := os.Getenv("APP_VERSION")
	if version == "" {
		version = "0.1.0"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "online",
		"mode":    mode,
		"version": version,
	})
}

package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"service/config"
	"service/models"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest

	// Validasi body request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email dan password wajib diisi dengan format yang benar",
		})
		return
	}

	// Cari user berdasarkan email
	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Email atau password salah",
		})
		return
	}

	// Cek apakah akun aktif
	if user.Status == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Akun Anda tidak aktif. Hubungi administrator.",
		})
		return
	}

	if user.Role != 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Akun Anda tidak valid. Hubungi administrator.",
		})
		return
	}

	// Verifikasi password (bcrypt — sama dengan yang digunakan Laravel)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Email atau password salah",
		})
		return
	}

	// Generate JWT token
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "JWT_SECRET tidak dikonfigurasi",
		})
		return
	}

	claims := jwt.MapClaims{
		"sub":   user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(24 * time.Hour).Unix(), // token berlaku 24 jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal membuat token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      tokenString,
		"token_type": "Bearer",
		"expires_in": 86400, // 24 jam dalam detik
		"user": gin.H{
			"id":     user.ID,
			"name":   user.Name,
			"email":  user.Email,
			"role":   user.Role,
			"status": user.Status,
		},
	})
}

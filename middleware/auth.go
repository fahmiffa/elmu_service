package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth adalah middleware untuk memvalidasi Bearer Token JWT pada setiap request.
// Secret JWT diambil dari environment variable JWT_SECRET.
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header tidak ditemukan",
			})
			return
		}

		// Pastikan format: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Format Authorization header tidak valid. Gunakan: Bearer <token>",
			})
			return
		}

		tokenString := parts[1]
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "JWT_SECRET tidak dikonfigurasi di server",
			})
			return
		}

		// Parse dan validasi token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan algoritma adalah HMAC (HS256)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token tidak valid atau sudah kadaluarsa",
			})
			return
		}

		// Simpan claims ke context agar bisa diakses di controller
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("claims", claims)
			// Laravel JWT menyimpan user id di "sub"
			if sub, exists := claims["sub"]; exists {
				c.Set("user_id", sub)
			}
		}

		c.Next()
	}
}

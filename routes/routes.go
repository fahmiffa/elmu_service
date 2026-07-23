package routes

import (
	"service/controllers"
	"service/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           12 * time.Hour,
	}))

	// ─── Public routes (tanpa autentikasi) ───────────────────────────────────
	r.POST("/auth/login", controllers.Login)

	// ─── Protected routes (wajib JWT Bearer Token) ───────────────────────────
	auth := r.Group("/")
	auth.Use(middleware.JWTAuth())
	{
		auth.GET("/users", controllers.GetUsers)
		auth.GET("/users/:id", controllers.GetUserByID)
		auth.POST("/users", controllers.CreateUser)
		auth.PUT("/users/:id", controllers.UpdateUser)
		auth.DELETE("/users/:id", controllers.DeleteUser)
		auth.GET("/dashboard/salary", controllers.GetSalaryDashboard)
		auth.POST("/salary/generate", controllers.GenerateSalary)
		auth.GET("/dashboard/pembayaran/bulanan", controllers.GetMonthlyPayments)
	}

	return r
}
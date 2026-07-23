package main

import (
	"log"
	"os"

	"service/config"
	// "service/models"
	"service/routes"

	"github.com/joho/godotenv"
)

func main() {

	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env tidak ditemukan, menggunakan environment system")
	} else {
		log.Println("✅ .env berhasil dimuat")
	}

	// Connect Database
	config.ConnectDB()
	log.Println("✅ Database Connected")

	// Auto Migrate
	// config.DB.AutoMigrate(&models.User{})
	// log.Println("✅ Database Migrated")

	// Setup Router
	r := routes.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("===================================")
	log.Printf("🚀 Server berjalan di http://localhost:%s\n", port)
	log.Println("===================================")

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
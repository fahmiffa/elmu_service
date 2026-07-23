package models

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"` // disembunyikan dari JSON response
	Role     int    `json:"role"`
	Status   int    `json:"status"`
}